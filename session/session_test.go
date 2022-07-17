package session

import (
	"fmt"
	"gonetlib/generator/idl"
	. "gonetlib/node"
	"gonetlib/util"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func SendMessage(node *UserNode, message string) {
	id, protocol := idl.NEW_PACKET_REQ_ECHO()

	protocol.Message = message

	header := NodeHeader{}
	header.PacketID = id
	header.Length = uint16(util.Sizeof(reflect.ValueOf(protocol)))

	node.Send(header, protocol)
}

var (
	serverSession *Session
	clientSession *Session

	wg sync.WaitGroup
)

func TestConnect(t *testing.T) {
	RegisterTask()

	server, client := net.Pipe()

	serverSession = NewSession()
	serverSession.Setup(1, server, NewUserNode(serverSession))

	serverSession.Start()

	clientSession = NewSession()
	clientSession.Setup(2, client, NewUserNode(clientSession))

	clientSession.Start()

	wg.Add(1)
	go communication()

	wg.Wait()

	time.Sleep(10 * time.Second)
}

func RegisterTask() {
	idl.Add_PACKET_REQ_ECHO_TASK_REGISTER()
	idl.Add_PACKET_RES_ECHO_TASK_REGISTER()
}

func communication() {
	if clientSession == nil {
		fmt.Println("client session is nullptr")
		return
	}

	tick := time.Tick(time.Second)
	terminate := time.After(10 * time.Second)

	for {
		select {
		case <-terminate:
			fmt.Println("now client session close...")
			clientSession.Close()
			wg.Done()
			return

		case <-tick:
			for i := 1; i <= 10; i++ {
				SendMessage(clientSession.GetNode().(*UserNode), fmt.Sprintf("hi my name is... %d", i))
			}
		}
	}
}
