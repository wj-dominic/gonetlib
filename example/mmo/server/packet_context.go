package mmo_server

import (
	"gonetlib/task"
	"math/rand"
	"sync"
)

// /////////////////////////////////////////////////////////////////////////////////////////////////
// packet context task
// /////////////////////////////////////////////////////////////////////////////////////////////////
type IPacketContextTask interface {
	Await(func(interface{}, error)) IPacketContextTask
	Start(...interface{})
}

type PacketContextTask[Out any] struct {
	ctx *PacketContext
	job task.ITask[Out]
}

func (task *PacketContextTask[Out]) Await(job func(interface{}, error)) IPacketContextTask {
	if task.ctx != nil && task.job != nil {
		task.ctx.wg.Add(1)
		task.job.Await(func(o Out, err error) {
			defer func() {
				task.ctx.wg.Done()
			}()

			job(o, err)
		})
	}

	return task
}

func (task *PacketContextTask[Out]) Start(i ...interface{}) {
	if task.ctx != nil && task.job != nil {
		task.job.Start(i...)
	}
}

// /////////////////////////////////////////////////////////////////////////////////////////////////
// packet context
// /////////////////////////////////////////////////////////////////////////////////////////////////s
type IPacketContext interface {
	GetNode() INode
	GetPacket() interface{}
	Async(func(...interface{}) (interface{}, error), ...uint8) IPacketContextTask
	Wait()
}

type PacketContext struct {
	node   INode
	packet interface{}
	wg     sync.WaitGroup
}

func CreateContext(node INode, packet interface{}) IPacketContext {
	return &PacketContext{
		node:   node,
		packet: packet,
		wg:     sync.WaitGroup{},
	}
}

func (ctx *PacketContext) GetNode() INode {
	return ctx.node
}

func (ctx *PacketContext) GetPacket() interface{} {
	return ctx.packet
}

func (ctx *PacketContext) Async(job func(...interface{}) (interface{}, error), numOfThread ...uint8) IPacketContextTask {
	_numOfThread := uint8(rand.Intn(7) + 1)
	if len(numOfThread) > 0 {
		_numOfThread = numOfThread[0]
	}

	ctx.wg.Add(1)
	_job := task.New(func(i ...interface{}) (interface{}, error) {
		defer func() {
			ctx.wg.Done()
		}()

		return job(i...)
	}, _numOfThread)

	return &PacketContextTask[interface{}]{
		ctx: ctx,
		job: _job,
	}
}

func (ctx *PacketContext) Wait() {
	ctx.wg.Wait()
}
