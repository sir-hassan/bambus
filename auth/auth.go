package auth

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sir-hassan/bambus/config"
	"net/http"
)

type Authenticator interface {
	Auth(req *http.Request) ([]string, error)
}

func AuthenticatorCreator(cc config.ComponentConfig) (Authenticator, error) {
	if cc.IsEmpty() {
		return nil, errors.New("config is empty")
	}

	switch cc.Type {
	case "dummy":
		return newDummyAuthenticator(cc.Settings)
	default:
		return nil, fmt.Errorf("invalid authenticator type: %s", cc.Type)
	}
}

type dummyAuthenticatorConfig struct {
	channelPrefix string
}

type DummyAuthenticator struct {
	config dummyAuthenticatorConfig
}

func newDummyAuthenticator(settings config.Settings) (*DummyAuthenticator, error) {
	cfg := dummyAuthenticatorConfig{}

	val, err := settings.GetStringField(
		"channelPrefix",
		config.WithDefaultString("qq"),
		config.NoEmptyString,
	)
	if err != nil {
		return nil, err
	}
	cfg.channelPrefix = val
	return &DummyAuthenticator{config: cfg}, nil
}

var _ Authenticator = DummyAuthenticator{}

func (d DummyAuthenticator) Auth(req *http.Request) ([]string, error) {
	if !bytes.HasSuffix([]byte(req.URL.Path), []byte("/dummy")) {
		return nil, nil
	}

	prefix := d.config.channelPrefix
	channels := []string{prefix + "1", prefix + "2", prefix + "3", prefix + "4"}
	return channels, nil
}
