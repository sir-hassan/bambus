package config

import (
	"errors"
	"fmt"
)

type ComponentConfig struct {
	Type     string   `json:"type"`
	Settings Settings `json:"settings"`
}

func (cc ComponentConfig) IsEmpty() bool {
	if cc.Type == "" || cc.Settings == nil || len(cc.Settings) == 0 {
		return true
	}
	return false
}

type AppConfig struct {
	Authenticator ComponentConfig `json:"authenticator"`
	Backend       ComponentConfig `json:"backend"`
	Frontend      Settings        `json:"frontend"`
}

type Settings map[string]interface{}

type MissingKeyError struct {
	error
	key string
}

func (err MissingKeyError) String() string {
	return fmt.Sprintf("field(%s) is required", err.key)
}

func (s Settings) getStringField(key string) (string, error) {
	if key == "" {
		return "", errors.New("empty key string is not allowed")
	}
	val, ok := s[key]
	if !ok {
		return "", MissingKeyError{key: key}
	}
	stringVal, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("field(%s) is not string type", key)
	}
	return stringVal, nil
}

type StringOption func(key string, val string, err error) (string, error)

func (s Settings) GetStringField(key string, options ...StringOption) (string, error) {
	val, err := s.getStringField(key)
	for _, option := range options {
		val, err = option(key, val, err)
	}
	return val, err
}

func WithDefaultString(defaultString string) StringOption {
	return func(key, val string, err error) (string, error) {
		if err != nil && errors.As(err, &MissingKeyError{}) {
			return defaultString, nil
		}
		return val, err
	}
}

func NoEmptyString(key, val string, err error) (string, error) {
	if err == nil && val == "" {
		return "", fmt.Errorf("field(%s) is empty", key)
	}
	return val, err
}
