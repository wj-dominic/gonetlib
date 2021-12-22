package trainer

import (
	"context"
	. "gonetlib/netlogger"
	. "gonetlib/routine"
)

type Trainer struct {
	id	 			uint8
	routines		*chan Routine

	cancelCtx		context.Context
	cancelFunc		context.CancelFunc
}

func NewTrainer(id uint8, routines *chan Routine) *Trainer{
	if routines == nil {
		GetLogger().Error("routines is nullptr")
		return nil
	}

	trainer := Trainer{
		id : id,
		routines: routines,
	}

	trainer.cancelCtx, trainer.cancelFunc = context.WithCancel(context.Background())

	return &trainer
}

func (trainer *Trainer) Start() {
	go trainer.workout()
	GetLogger().Info("start to workout | trainerID[%d]", trainer.id)
}

func (trainer *Trainer) Stop(){
	trainer.cancelFunc()
	GetLogger().Info("stop workout | trainerID[%d]", trainer.id)
}

func (trainer *Trainer) workout() {
	for {
		select {
		case routine := <-*trainer.routines:
			routine.Workout()
			break

		case _ = <-trainer.cancelCtx.Done():
			GetLogger().Debug("workout done..! | trainerID[%d]", trainer.id)
			return
		}
	}
}