package commandcenter

import (
	. "message"
	. "scv"
	. "session"
)

type CommandCenter struct{
	SCVs	[]SCV
}

func NewCommandCenter(scvCount... uint32) *CommandCenter {
	maxSCVs := uint32(1)
	if len(scvCount) > 0 {
		maxSCVs = scvCount[0]
	}

	return &CommandCenter{
		SCVs : make([]SCV, maxSCVs),
	}
}

func (center *CommandCenter) OnRecv(session *Session, msg *Message) {



}


