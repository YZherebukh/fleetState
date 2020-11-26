package config

import (
	"errors"
)

// package errors
var (
	ErrEmptyConfiguration = errors.New("empty configuration")
)

// Configuration is an interface for local config
type Configuration interface {
	Service() Service
	Logger() Logger
}
