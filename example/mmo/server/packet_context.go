package mmo_server

import (
	"math/rand"
	"sync"

	"github.com/wj-dominic/gonetlib/message"
	"github.com/wj-dominic/gonetlib/task"
)

// IPacketContextTask ...
type IPacketContextTask interface {
	Await(func(result interface{}, err error)) IPacketContextTask
	Start(...interface{})
}

type PacketContextTask[Out any] struct {
	wg  *sync.WaitGroup
	job task.Task[Out]
}

func (task *PacketContextTask[Out]) Await(job func(interface{}, error)) IPacketContextTask {
	if task.wg != nil && task.job != nil {
		task.wg.Add(1)
		task.job.Await(func(o Out, err error) {
			defer func() {
				task.wg.Done()
			}()

			job(o, err)
		})
	}

	return task
}

func (task *PacketContextTask[Out]) Start(i ...interface{}) {
	if task.job != nil {
		task.job.Start(i...)
	}
}

// IPacketContext ...
type IPacketContext interface {
	SetNode(INode)
	UnPack(*message.Message) error
	RunHandler(...uint8)
	Async(func(params ...interface{}) (interface{}, error), ...uint8) IPacketContextTask
	Wait()
}

// PacketContextTwoWay ...
type PacketContextTwoWay[TRequest any, TResponse any] struct {
	node     INode
	Request  TRequest
	Response TResponse
	wg       sync.WaitGroup
	handler  func(*PacketContextTwoWay[TRequest, TResponse]) error
}

func NewPacketContextTwoWay[TRequest any, TResponse any](handler func(*PacketContextTwoWay[TRequest, TResponse]) error) IPacketContext {
	return &PacketContextTwoWay[TRequest, TResponse]{
		node:    nil,
		handler: handler,
		wg:      sync.WaitGroup{},
	}
}

func (ctx *PacketContextTwoWay[TRequest, TResponse]) SetNode(node INode) {
	ctx.node = node
}

func (ctx *PacketContextTwoWay[TRequest, TResponse]) Send(message interface{}) {
	ctx.node.Send(message)
}

func (ctx *PacketContextTwoWay[TRequest, TResponse]) SendResponse() {
	ctx.node.Send(ctx.Response)
}

func (ctx *PacketContextTwoWay[TRequest, TResponse]) UnPack(packet *message.Message) error {
	packet.Pop(&ctx.Request)
	return nil
}

func (ctx *PacketContextTwoWay[TRequest, TResponse]) RunHandler(numOfThread ...uint8) {
	_numOfThread := uint8(rand.Intn(7) + 1)
	if len(numOfThread) > 0 {
		_numOfThread = numOfThread[0]
	}

	ctx.wg.Add(1)
	job := task.New(func(i ...interface{}) (error, error) {
		defer func() {
			ctx.wg.Done()
		}()

		err := ctx.handler(ctx)
		return err, nil
	}, _numOfThread)

	job.Start()
}

func (ctx *PacketContextTwoWay[TRequest, TResponse]) Async(job func(...interface{}) (interface{}, error), numOfThread ...uint8) IPacketContextTask {
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
		wg:  &ctx.wg,
		job: _job,
	}
}

func (ctx *PacketContextTwoWay[TRequest, TResponse]) Wait() {
	ctx.wg.Wait()
}

// PacketContextOneWay ...
type emptyResponse struct{}

type PacketContextOneWay[TRequest any] struct {
	PacketContextTwoWay[TRequest, emptyResponse]
}

func NewPacketContextOneWay[TRequest any](handler func(*PacketContextOneWay[TRequest]) error) IPacketContext {
	ctx := &PacketContextOneWay[TRequest]{
		PacketContextTwoWay[TRequest, emptyResponse]{
			node:    nil,
			handler: nil,
			wg:      sync.WaitGroup{},
		},
	}

	_handler := func(twoWayContext *PacketContextTwoWay[TRequest, emptyResponse]) error {
		return handler(ctx)
	}

	ctx.handler = _handler
	return ctx
}

func (ctx *PacketContextOneWay[TRequest]) SetNode(node INode) {
	ctx.PacketContextTwoWay.SetNode(node)
}

func (ctx *PacketContextOneWay[TRequest]) Send(message interface{}) {
	ctx.PacketContextTwoWay.Send(message)
}

func (ctx *PacketContextOneWay[TRequest]) UnPack(packet *message.Message) error {
	ctx.PacketContextTwoWay.UnPack(packet)
	return nil
}

func (ctx *PacketContextOneWay[TRequest]) RunHandler(numOfThread ...uint8) {
	ctx.PacketContextTwoWay.RunHandler(numOfThread...)
}

func (ctx *PacketContextOneWay[TRequest]) Async(job func(...interface{}) (interface{}, error), numOfThread ...uint8) IPacketContextTask {
	return ctx.PacketContextTwoWay.Async(job, numOfThread...)
}

func (ctx *PacketContextOneWay[TRequest]) Wait() {
	ctx.PacketContextTwoWay.Wait()
}
