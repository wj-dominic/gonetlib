package session

import (
	"fmt"
	. "gonetlib/message"
	. "gonetlib/netlogger"
	"net"
	"testing"
	"time"
)

var (
	serverSession *Session
	clientSession *Session
)

type MyNode struct{
	id int
}

func (node MyNode) OnConnect(){
	fmt.Println("call OnConnect!")
}

func (node MyNode) OnDisconnect(){
	fmt.Println("call OnDisConnect!")
}

func (node MyNode) OnRecv(packet *Message) bool{
	if packet == nil {
		GetLogger().Error("packet is nullptr")
		return false
	}

	fmt.Printf("call OnRecv! | recvLength[%d]", packet.GetLength())

	return true
}

func (node MyNode) OnSend(sendBytes int){
	fmt.Printf("call OnSend! | sendBytes[%d]", sendBytes)
}

func TestConnect(t *testing.T){
	server, client := net.Pipe()

	fmt.Println("starting server...!")

	serverSession = NewSession()
	serverNode := MyNode{}
	serverSession.Setup(1, server, serverNode)

	serverSession.Start()



	clientSession = NewSession()
	clientNode := MyNode{}
	clientSession.Setup(2, client, clientNode)

	clientSession.Start()


	communication()

}

func communication() {
	if clientSession == nil {
		fmt.Println("client session is nullptr")
		return
	}

	for {
		packet := NewMessage(true)
		packet.SetHeader(SYN, XOR)
		packet.Push("hello world")

		if clientSession.SendPost(packet) == false {
			fmt.Println("send post failed..")
			break
		}

		fmt.Println("success to send to server : ", packet.GetPayloadBuffer())

		time.Sleep(10 * time.Second)
	}
}