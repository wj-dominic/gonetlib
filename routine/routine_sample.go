package routine

import (
	"fmt"
	. "gonetlib/message"
	. "gonetlib/netlogger"
)


const (
	SampleProtocolID uint16 = 10001
)

type ReqSampleProtocol struct{
	Name string
	Value1 uint64
	Value2 uint32
}

type ResSampleProtocol struct{
	Result 	string
	Value 	uint32
}

// SampleRoutine : 패킷 ID와 대응되는 루틴
// TODO 여기 샘플들을 자동화 코드로 만들 수 있어야 함
// TODO 제너릭을 사용할 수 있는지 확인 필요
type SampleRoutine struct {
	request 	ReqSampleProtocol
	response 	ResSampleProtocol
}

// 프로토콜에 해당하는 루틴 로직
func (routine *SampleRoutine) Workout() bool {
	fmt.Println("Proc... sample routine...")
	fmt.Printf("name[%s] value1[%d] value2[%d]\n", routine.request.Name, routine.request.Value1, routine.request.Value2)

	return true
}

// SampleRoutineRegister : 루틴을 생성할 수 있는 루틴 생성 클래스 정의
type SampleRoutineRegister struct {

}

func NewSampleRoutineRegister() *SampleRoutineRegister{
	return &SampleRoutineRegister{}
}

// Make : 대응되는 루틴을 이 루틴 생성 함수에서 생성한다.
func (register *SampleRoutineRegister) Make(packet *Message) Routine{
	if packet == nil {
		GetLogger().Error("packet is nullptr")
		return nil
	}

	var routine SampleRoutine

	//마샬링
	packet.Pop(&routine.request)

	return &routine
}


