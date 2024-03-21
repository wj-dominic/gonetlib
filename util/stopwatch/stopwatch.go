package stopwatch

import "time"

type IStopWatch interface {
	Start() IStopWatch
	Stop() IStopWatch
	ElapsedTime() time.Duration
	Reset() IStopWatch
}

type StopWatch struct {
	start time.Time
	stop  time.Time
}

func (watch *StopWatch) Start() IStopWatch {
	watch.start = time.Now()
	return watch
}

func (watch *StopWatch) Stop() IStopWatch {
	watch.stop = time.Now()
	return watch
}

func (watch *StopWatch) Reset() IStopWatch {
	watch.start = time.Time{}
	watch.stop = time.Time{}
	return watch
}

func (watch *StopWatch) ElapsedTime() time.Duration {
	return watch.stop.Sub(watch.start)
}

func New() IStopWatch {
	return &StopWatch{
		start: time.Time{},
		stop:  time.Time{},
	}
}
