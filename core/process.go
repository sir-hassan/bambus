package core

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/sir-hassan/bambus/backend"
)

type hubProcess struct {
	logger      log.Logger
	tube        backend.Tube
	connects    chan contact
	disConnects chan contact
	pusher      DataPusher

	contactsTable *contactsTable
}

func NewHubProcess(logger log.Logger, tube backend.Tube, connects chan contact) {
	p := &hubProcess{
		logger:        logger,
		tube:          tube,
		connects:      connects,
		disConnects:   make(chan contact),
		contactsTable: newContactsTable(),
	}
	p.pusher = NewQueuingPusher(logger, p.disConnects)

	go p.start()
}

func (r *hubProcess) start() {
	for {
		select {
		case sr := <-r.disConnects:
			r.contactsTable.Remove(sr.socketQueue)
			garbageChannels := r.contactsTable.GarbageCollect()
			err := r.tube.UnSubscribe(garbageChannels)
			if err != nil {
				level.Error(r.logger).Log("msg", "backend unsubscribe", "err", err)
			}
		case sr := <-r.connects:
			sr.passUnPlug <- func() {
				r.disConnects <- sr
			}
			r.contactsTable.Add(sr.socketQueue, sr.channels)
			// todo: shouldn't resubscribe to the same channels again.
			err := r.tube.Subscribe(sr.channels)
			if err != nil {
				level.Error(r.logger).Log("msg", "backend subscribe", "err", err)
			}
		case msg := <-r.tube.Feed():
			socs := r.contactsTable.GetSocketsInChannel(msg.Channel)
			r.pusher.Push(msg.Data, socs)
		}
	}
}
