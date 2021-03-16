package dummy

import (
	"fmt"
	"github.com/sir-hassan/bambus/backend"
	"time"
)

type dummyTube struct {
	messagesChan    chan *backend.Message
	subscribeChan   chan []string
	unSubscribeChan chan []string

	channels []string
}

var _ backend.Tube = &dummyTube{}

func NewDummyTube() *dummyTube {
	tube := &dummyTube{
		messagesChan:    make(chan *backend.Message),
		subscribeChan:   make(chan []string),
		unSubscribeChan: make(chan []string),
		channels:        make([]string, 0),
	}

	go func() {
		i := 0
		for {
			select {
			case channels := <-tube.subscribeChan:
			outer:
				for _, c := range channels {
					for _, channel := range tube.channels {
						if c == channel {
							continue outer
						}
					}
					tube.channels = append(tube.channels, c)
				}

			//case unSub := <- tube.unSubscribeChan:

			case <-time.After(10 * time.Millisecond):
			}

			if len(tube.channels) == 0 {
				continue
			}

			c := tube.channels[i%len(tube.channels)]
			newMsg := &backend.Message{
				Channel: c,
				Data:    []byte(fmt.Sprintf("data from channel: %s, time: %s", c, time.Now().String())),
			}

			select {
			case tube.messagesChan <- newMsg:
				i++
			case <-time.After(10 * time.Millisecond):
			}

		}
	}()

	return tube
}

func (t *dummyTube) Subscribe(channels []string) error {
	t.subscribeChan <- channels
	return nil
}

func (t *dummyTube) UnSubscribe(channels []string) error {
	//t.unSubscribeChan <- channels
	return nil
}

func (t *dummyTube) Feed() <-chan *backend.Message {
	return t.messagesChan
}
