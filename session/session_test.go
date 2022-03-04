package session

import (
	"fmt"
	"gonetlib/gym"
	. "gonetlib/node"
	"gonetlib/routine"
	"gonetlib/util"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func SendSample(node *UserNode, name string, value1 uint64, value2 uint32) {
	protocol := routine.ReqSampleProtocol{}
	protocol.Name = name
	protocol.Value1 = value1
	protocol.Value2 = value2

	header := NodeHeader{}
	header.PacketID = routine.SampleProtocolID
	header.Length = uint16(util.Sizeof(reflect.ValueOf(protocol)))

	node.Send(header, protocol)
}

var (
	serverSession *Session
	clientSession *Session

	serverNode *UserNode
	clientNode *UserNode

	wg sync.WaitGroup
)

func TestConnect(t *testing.T) {
	RegisterRoutine()

	server, client := net.Pipe()

	serverSession = NewSession()
	serverNode = NewUserNode()
	serverNode.SetSession(serverSession)

	serverSession.Setup(1, server, serverNode)

	serverSession.Start()

	clientSession = NewSession()
	clientNode = NewUserNode()
	clientNode.SetSession(clientSession)

	clientSession.Setup(2, client, clientNode)

	clientSession.Start()

	wg.Add(1)
	go communication()

	wg.Wait()

	for {
		select {
		case <-time.After(10 * time.Second):
			return
		}
	}
}

func RegisterRoutine() {
	gyms := gym.GetGyms()
	gyms.CreateGym(gym.GymMain, 1, 1)

	routineMaker := routine.GetRoutineMaker()
	routineMaker.AddRegister(routine.SampleProtocolID, routine.NewSampleRoutineRegister())
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
				SendSample(clientNode, "dogSyeon", uint64(i*100), uint32(i))
			}
		}
	}
}
