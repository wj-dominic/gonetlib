package routine

import (
	. "gonetlib/Message"
	. "gonetlib/netlogger"
)

// SampleRoutine : 패킷 ID와 대응되는 루틴
// TODO 여기 샘플들을 자동화 코드로 만들 수 있어야 함
// TODO 제너릭을 사용할 수 있는지 확인 필요
type SampleRoutine struct {
	name 	string
	value1 	uint64
	value2	uint32
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
	packet.Pop(&routine.name)
	packet.Pop(&routine.value1)
	packet.Pop(&routine.value2)

	return &routine
}


