package session

import (
	"fmt"
	. "gonetlib/message"
	. "gonetlib/netlogger"
	"net"
	"sync"
	"testing"
	"time"
)

var (
	serverSession *Session
	clientSession *Session
	wg	sync.WaitGroup
)

type MyNode struct{
	session *Session
}

func (node MyNode) OnConnect(){
	fmt.Printf("call OnConnect! | sessionID[%d]\n", node.session.id)
}

func (node MyNode) OnDisconnect(){
	fmt.Printf("call OnDisConnect! | sessionID[%d]\n", node.session.id)
}

func (node MyNode) OnRecv(packet *Message) bool{
	if packet == nil {
		GetLogger().Error("packet is nullptr")
		return false
	}

	fmt.Printf("call OnRecv! | sessionID[%d] recvLength[%d]\n", node.session.id, packet.GetLength())

	return true
}

func (node MyNode) OnSend(sendBytes int){
	fmt.Printf("call OnSend! | sessionID[%d] sendBytes[%d]\n", node.session.id, sendBytes)
}

func TestConnect(t *testing.T){
	server, client := net.Pipe()

	fmt.Println("starting server...!")

	serverSession = NewSession()
	serverNode := MyNode{serverSession}
	serverSession.Setup(1, server, serverNode)

	serverSession.Start()



	clientSession = NewSession()
	clientNode := MyNode{clientSession}
	clientSession.Setup(2, client, clientNode)

	clientSession.Start()

	wg.Add(1)
	go communication()

	wg.Wait()

	for {
		select {
		case <- time.After(10 * time.Second):
			return
		}
	}
}

func communication() {
	if clientSession == nil {
		fmt.Println("client session is nullptr")
		return
	}

	tick := time.Tick(time.Second)
	terminate := time.After(3 * time.Second)

	for {
		select {
		case <-terminate:
			fmt.Println("now client session close...")
			clientSession.Close()
			wg.Done()
			return

		case <-tick:
			packet := NewMessage(true)
			packet.Push("hello world")
			packet.SetHeader(ESTABLISHED, NONE)

			for i := 0 ; i < 10 ; i++ {
				if clientSession.SendPost(packet) == false {
					fmt.Println("send post failed..")
					break
				}
			}

			fmt.Println("success to send to server : ", packet.GetPayloadBuffer())
		}
	}
}