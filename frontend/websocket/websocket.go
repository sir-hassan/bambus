package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/sir-hassan/bambus/frontend"
	"net/http"
	"sync"
)

func ignoreOrigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{CheckOrigin: ignoreOrigin}

type socket struct {
	conn *websocket.Conn
	lock *sync.Mutex
}

var _ frontend.Socket = &socket{}

func NewSocket(w http.ResponseWriter, r *http.Request) (*socket, error) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &socket{conn: wsConn, lock: new(sync.Mutex)}, nil
}

func (s *socket) Close() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.conn.Close()
}

func (s *socket) Read() ([]byte, error) {
	_, message, err := s.conn.ReadMessage()
	return message, err
}

func (s *socket) Write(data []byte) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.conn.WriteMessage(websocket.TextMessage, data)
}
