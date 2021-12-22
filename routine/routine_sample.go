package routine

import (
	"fmt"
	. "gonetlib/Message"
	. "gonetlib/netlogger"
)


type ReqSampleProtocol struct{
	name string
	value1 uint64
	value2 uint32
}

type ResSampleProtocol struct{
	result 	string
	value 	uint32
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

	return true
}

// SampleRoutineRegister : 루틴을 생성할 수 있는 루틴 생성 클래스 정의
type SampleRoutineRegister struct {

}

func NewSampleRoutineRegister() *SampleRoutineRegister{
	return &SampleRoutineRegister{}
}

// Make : 대응되는 루틴을 이 루틴 생성 함수에서 생성한다.
func (register *SampleRoutineRegister) Make(packet *Message) *SampleRoutine{
	if packet == nil {
		GetLogger().Error("packet is nullptr")
		return nil
	}

	var routine SampleRoutine

	//마샬링
	packet.Pop(&routine)

	return &routine
}


