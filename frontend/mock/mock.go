package mock

import (
	"github.com/sir-hassan/bambus/frontend"
)

var _ frontend.Socket = &Socket{}

type Socket struct {
	data []byte
}

func (s *Socket) Read() ([]byte, error) {
	return s.data, nil
}

func (s *Socket) Write(bytes []byte) error {
	s.data = append(s.data, bytes...)
	return nil
}

func (s *Socket) Close() error {
	panic("implement me")
}
