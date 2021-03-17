// Package backend defines abstractions that bambus core logic use to interact
// with pubsub brokers.
// The main one is the Tube interface which is basically a connection to the
// between bambus and the pubsub broker.

package backend

type Message struct {
	Channel string
	Data    []byte
}

type Tube interface {
	Subscribe(channels []string) error
	UnSubscribe(channels []string) error
	Feed() <-chan *Message
}

type TubeCreator func() (Tube, error)
