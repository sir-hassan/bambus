package core

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/sir-hassan/bambus/backend"
	"github.com/sir-hassan/bambus/frontend"
	"time"
)

type Hub struct {
	logger      log.Logger
	connects    chan contact
	tubeCreator backend.TubeCreator
}

func NewHub(logger log.Logger, tubeCreator backend.TubeCreator) *Hub {
	return &Hub{
		logger:      logger,
		connects:    make(chan contact),
		tubeCreator: tubeCreator,
	}
}

func (h Hub) Plug(soc frontend.Socket, channels []string) func() {

	sQueue := frontend.NewSocketQueue(soc)

	sr := contact{
		socketQueue: sQueue,
		channels:    channels,
		passUnPlug:  make(chan func()),
	}

	for {
		select {
		case h.connects <- sr:
			return <-sr.passUnPlug
		case <-time.After(time.Second):
			fmt.Println("creating hub process...")
			tube, err := h.tubeCreator()
			if err != nil {
				fmt.Println("redis dial error.")
			}
			NewHubProcess(h.logger, tube, h.connects)
		}
	}
}
