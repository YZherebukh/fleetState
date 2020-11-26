//go:generate mockgen -source ../config/config.go  -destination ../config/mock/mock_config.go

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

// Service is a struct with service configuration
type Service struct {
	Port       string
	VersionAPI string
}

// Logger is a struct with logger configuration
type Logger struct {
	FilePath string
	Level    string
}
