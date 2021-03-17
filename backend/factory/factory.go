package factory

import (
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/sir-hassan/bambus/backend"
	"github.com/sir-hassan/bambus/backend/redis"
	"github.com/sir-hassan/bambus/config"
)

func CreateTubeCreator(logger log.Logger, cc config.ComponentConfig) (backend.TubeCreator, error) {
	if cc.IsEmpty() {
		return nil, errors.New("config is empty")
	}
	switch cc.Type {
	case "redis":
		return func() (backend.Tube, error) {
			return redis.NewTube(logger, cc.Settings)
		}, nil
	default:
		return nil, fmt.Errorf("invalid backend type: %s", cc.Type)
	}
}
