package core

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/sir-hassan/bambus/frontend"
	"time"
)

type DataPusher interface {
	Push(data []byte, socketsQueues []frontend.SocketQueue)
}

type simplePusher struct {
	logger      log.Logger
	jobsQueue   chan frontend.SocketQueue
	disConnects chan contact
}

func NewSimplePusher(logger log.Logger, disConnects chan contact) *simplePusher {
	p := &simplePusher{
		logger:      logger,
		disConnects: disConnects,
		jobsQueue:   make(chan frontend.SocketQueue, 100),
	}

	return p
}

func (p *simplePusher) Push(data []byte, socketsQueues []frontend.SocketQueue) {
	for _, sQueue := range socketsQueues {
		err := sQueue.Socket.Write(data)
		if err != nil {
			level.Error(p.logger).Log("msg", "socket write", "error", err)
			p.disConnects <- contact{socketQueue: sQueue}
			err = sQueue.Socket.Close()
			if err != nil {
				level.Error(p.logger).Log("msg", "socket close", "error", err)
			}
		}
	}
}

type queuingPusher struct {
	logger      log.Logger
	jobsQueue   chan frontend.SocketQueue
	disConnects chan contact
}

func NewQueuingPusher(logger log.Logger, disConnects chan contact) *queuingPusher {
	p := &queuingPusher{
		logger:      logger,
		disConnects: disConnects,
		jobsQueue:   make(chan frontend.SocketQueue, 100),
	}

	go func() {
		level.Debug(logger).Log("msg", "new pusher routine")
		newJobGate := make(chan frontend.SocketQueue)
		for queue := range p.jobsQueue {
			select {
			case newJobGate <- queue:
			default:
				level.Debug(logger).Log("msg", "new socket queue routine")
				go p.work(newJobGate, queue)
			}
		}
	}()

	return p
}

func (p *queuingPusher) Push(data []byte, queues []frontend.SocketQueue) {
	for _, queue := range queues {
		// check if its not in  a worker
		select {
		case queue.Seated <- struct{}{}:
			p.jobsQueue <- queue
		default:
		}
		queue.Queue <- data
	}
}

func (p *queuingPusher) work(gate chan frontend.SocketQueue, queue frontend.SocketQueue) {
	p.loopQueue(queue)
	for {
		select {
		case queue := <-gate:
			p.loopQueue(queue)
		case <-time.After(time.Second):
			return
		}
	}
}

func (p *queuingPusher) loopQueue(queue frontend.SocketQueue) {
	for {
		select {
		case data := <-queue.Queue:
			if data == nil {
				err := queue.Socket.Close()
				if err != nil {
					level.Error(p.logger).Log("msg", "socket close", "socket", queue.Socket, "error", err)
				}
				return
			}
			if queue.Socket == nil {
				level.Error(p.logger).Log("msg", "dropped message", "socket", queue.Socket, "message", data)
				continue
			}
			err := queue.Socket.Write(data)
			if err != nil {
				level.Error(p.logger).Log("msg", "socket write", "socket", queue.Socket, "error", err)
				p.disConnects <- contact{socketQueue: queue}
				queue.Socket = nil
				continue
			}
		case <-time.After(time.Second):
			<-queue.Seated
			return
		}
	}
}
