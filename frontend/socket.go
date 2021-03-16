package frontend

import "net/http"

type Socket interface {
	Read() ([]byte, error)
	Write([]byte) error
	Close() error
}

type SocketCreator func(w http.ResponseWriter, r *http.Request) (Socket, error)

type SocketQueue struct {
	Socket Socket
	Queue  chan []byte
	Seated chan struct{}
}

func NewSocketQueue(socket Socket) SocketQueue {
	return SocketQueue{
		Socket: socket,
		Queue:  make(chan []byte, 100),
		Seated: make(chan struct{}, 1),
	}
}
