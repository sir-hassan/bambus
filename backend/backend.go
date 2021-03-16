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
