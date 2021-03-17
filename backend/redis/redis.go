package redis

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gomodule/redigo/redis"
	"github.com/sir-hassan/bambus/backend"
	"github.com/sir-hassan/bambus/config"
)

type ConnConfig struct {
	Network string
	Address string
	Options []redis.DialOption
}

type tube struct {
	logger  log.Logger
	conn    redis.PubSubConn
	msgChan chan *backend.Message
}

func NewTube(logger log.Logger, settings config.Settings) (*tube, error) {
	// todo: parse from setting.
	cfg := ConnConfig{
		Network: "tcp",
		Address: "localhost" + ":6379",
		Options: nil,
	}
	rc, err := redis.Dial(cfg.Network, cfg.Address, cfg.Options...)
	return &tube{
		logger:  logger,
		conn:    redis.PubSubConn{Conn: rc},
		msgChan: make(chan *backend.Message),
	}, err
}

var _ backend.Tube = &tube{}

func (t *tube) Subscribe(channels []string) error {
	channelsArray := make([]interface{}, len(channels))
	for i, v := range channels {
		channelsArray[i] = v
	}
	go t.start()
	return t.conn.Subscribe(channelsArray...)
}

func (t *tube) UnSubscribe(channels []string) error {
	if len(channels) == 0 {
		return nil
	}
	channelsArray := make([]interface{}, len(channels))
	for i, v := range channels {
		channelsArray[i] = v
	}
	return t.conn.Unsubscribe(channelsArray...)
}

func (t *tube) Feed() <-chan *backend.Message {
	return t.msgChan
}

func (t *tube) start() {
	for {
		data := t.conn.Receive()
		switch v := data.(type) {
		case redis.Message:
			t.msgChan <- &backend.Message{Channel: v.Channel, Data: v.Data}
		case redis.Subscription:
			level.Info(t.logger).Log("msg", "redis subscription feedback", "channel", v.Channel, "kind", v.Kind)
		case error:
			level.Error(t.logger).Log("msg", "redis receive", "error", v)
			return
		default:
			panic("unknown redis receive message type.")
		}
	}
}
