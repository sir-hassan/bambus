package core

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/sir-hassan/bambus/frontend"
	"github.com/sir-hassan/bambus/frontend/mock"
	"sync"
	"testing"
	"time"
)

var _ log.Logger = &logRecorder{}

type logRecorder struct {
	lock    *sync.Mutex
	entries []interface{}
}

func (r *logRecorder) Log(keyvals ...interface{}) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.entries = append(r.entries, keyvals)
	return nil
}

func TestQueuingPusher(t *testing.T) {
	sQueues := make([]frontend.SocketQueue, 3, 3)
	for i, _ := range sQueues {
		sQueues[i] = frontend.NewSocketQueue(&mock.Socket{})
	}
	logger := &logRecorder{lock: new(sync.Mutex)}
	disConnects := make(chan contact, 10)

	pusher := NewQueuingPusher(logger, disConnects)
	pusher.Push([]byte("hello"), sQueues)

	time.Sleep(time.Millisecond)

	// asset that all sockets have been written.
	for i, _ := range sQueues {
		data, _ := sQueues[i].Socket.Read()
		if string(data) != "hello" {
			t.Errorf("sockets didn't get written")
		}
	}

	pusher.Push([]byte("another hello"), sQueues)
	time.Sleep(time.Millisecond)

	// asset logger records.
	expectedEntries := []string{
		"[level debug msg new pusher routine]",
		"[level debug msg new socket queue routine]",
		"[level debug msg new socket queue routine]",
		"[level debug msg new socket queue routine]",
	}
	for i, entry := range logger.entries {
		if expectedEntries[i] != fmt.Sprintf("%s", entry) {
			t.Errorf("expected logs didn't match, exp: %s, got: %s", expectedEntries[i], fmt.Sprintf("%s", entry))
		}
	}
}
