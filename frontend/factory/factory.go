package factory

import (
	"fmt"
	"github.com/sir-hassan/bambus/frontend"
	"github.com/sir-hassan/bambus/frontend/sse"
	"github.com/sir-hassan/bambus/frontend/websocket"
	"net/http"
)

func CreateSocketCreator(socketType string) (frontend.SocketCreator, error) {
	switch socketType {
	case "websocket":
		return func(w http.ResponseWriter, r *http.Request) (frontend.Socket, error) {
			return websocket.NewSocket(w, r)
		}, nil
	case "sse":
		return func(w http.ResponseWriter, r *http.Request) (frontend.Socket, error) {
			return sse.NewSocket(w, r)
		}, nil
	default:
		return nil, fmt.Errorf("invalid socket type: %s", socketType)
	}
}
