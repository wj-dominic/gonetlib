package mmo_server_test

import (
	"fmt"
	mmo_server "gonetlib/example/mmo/server"
	"gonetlib/logger"
	"gonetlib/message"
	"gonetlib/server"
	"gonetlib/util"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

type RequestLogin struct {
	Id   uint64
	Name string
}

type ResponseLogin struct {
	Result  int
	Message string
}

func RequestLoginHandler(ctx *mmo_server.PacketContextTwoWay[RequestLogin, ResponseLogin]) error {
	//패킷 데이터 참조하기
	fmt.Println("request : ", ctx.Request)

	//비동기 작업 요청하기
	ctx.Async(func(i ...interface{}) (interface{}, error) {
		//실제로 수행되는 비동기 작업
		time.Sleep(time.Second * 5)
		return i[0].(int) + i[1].(int), nil
	}).Await(func(result interface{}, err error) {
		//비동기 끝나고 호출되는 콜백
		fmt.Println("async result is ", result)

		ctx.Response.Result = result.(int)
		ctx.Response.Message = fmt.Sprintf("HELLO %s, Wellcome to the world!", ctx.Request.Name)

		//노드 획득 후 사용하기
		ctx.SendResponse()
	}).Start(1, 2)

	return nil
}

func TestServer(t *testing.T) {
	config := logger.CreateLoggerConfig().
		WriteToConsole().
		WriteToFile(
			logger.WriteToFile{
				Filepath:        "./MMOServer.log",
				RollingInterval: logger.RollingIntervalDay,
			}).
		MinimumLevel(logger.DebugLevel).
		TickDuration(1000)
	_logger := config.CreateLogger()

	builder := server.CreateServerBuilder()
	builder.Configuration(server.ServerInfo{
		Id:         1,
		Address:    server.Endpoint{IP: "0.0.0.0", Port: 50000},
		Protocols:  server.TCP | server.UDP,
		MaxSession: 10000,
	})
	builder.Logger(_logger)
	builder.Handler(mmo_server.CreateMMOServer())

	//패킷을 받았을 때 수행할 핸들러를 등록한다.
	//TODO:Packet type으로 ID를 획득할 수 있도록 해야 함
	mmo_server.AddPacketContext(1, mmo_server.CreatePacketContextTwoWay(RequestLoginHandler))

	server := builder.Build()
	server.Run()

	wg.Add(1)
	go ConnectToServer()

	wg.Wait()

	server.Stop()
	_logger.Info("Success to stop the server")
}

func ConnectToServer() {
	defer func() {
		wg.Done()
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:50000")
	if err != nil {
		return
	}

	defer func() {
		conn.Close()
	}()

	Send(conn, 1, RequestLogin{
		Id:   1,
		Name: "GONETLIB",
	})

	buffer := make([]byte, 1024)
	recvBytes, err := conn.Read(buffer)
	if err != nil {
		return
	}

	packet := message.NewMessage()
	packet.MoveRear(uint16(recvBytes))
	copy(packet.GetBuffer(), buffer[:recvBytes])

	response := ResponseLogin{}
	packet.Pop(&response)

	fmt.Println("response : ", response)
}

func Send[TStructure any](conn net.Conn, id uint16, request TStructure) {
	packet := message.NewMessage()
	packet.Push(id)
	packet.Push(uint16(util.Sizeof(reflect.ValueOf(request))))
	packet.Push(request)
	packet.MakeHeader()

	conn.Write(packet.GetBuffer())
}
