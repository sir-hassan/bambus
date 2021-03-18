package sse

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sir-hassan/bambus/frontend"
	"io"
	"net/http"
	"net/http/httputil"
)

var _ frontend.Socket = &socket{}

type socket struct {
	conn     io.ReadWriteCloser
	buffered *bufio.ReadWriter
}

func NewSocket(w http.ResponseWriter, r *http.Request) (*socket, error) {
	// setup stream response.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.(http.Flusher).Flush()

	hj, ok := w.(http.Hijacker)
	if !ok {
		return nil, errors.New("couldn't hijack connection")
	}
	conn, buffered, err := hj.Hijack()
	if err != nil {
		return nil, errors.New("couldn't hijack connection")
	}

	return &socket{conn: conn, buffered: buffered}, nil
}

func (s *socket) Read() ([]byte, error) {
	data := make([]byte, 100)
	_, err := s.buffered.Read(data)
	return data, err
}

func (s *socket) Write(bytes []byte) error {
	chunkWriter := httputil.NewChunkedWriter(s.buffered)
	_, err := chunkWriter.Write([]byte(fmt.Sprintf("data: %s\n\n", bytes)))
	if err != nil {
		return err
	}
	return s.buffered.Flush()
}

func (s *socket) Close() error {
	return s.conn.Close()
}
