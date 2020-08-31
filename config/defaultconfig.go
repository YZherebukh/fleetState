package config

import (
	"fmt"
	"os"
	"sync"
)

// package const
const (
	port        = "port"
	versionAPI  = "versionAPI"
	logFilePath = "logFilePath"
	logLevel    = "logLevel"
)

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

// config is a struct with concurent safe public method to access a config
type config struct {
	mu      *sync.RWMutex
	service Service
	logger  Logger
}

// new creates a new config instace
func new() *config {
	return &config{
		mu: &sync.RWMutex{},
	}
}

// New initiates a new Configuration instance
func New() (Configuration, error) {
	config := new()

	err := config.setService(os.Getenv(port), os.Getenv(port))
	if err != nil {
		return nil, fmt.Errorf("set up config failed. error %s", err.Error())
	}

	err = config.setLogger(os.Getenv(logFilePath), os.Getenv(logLevel))
	if err != nil {
		return nil, fmt.Errorf("set up config failed. error %s", err.Error())
	}

	return config, nil
}

// Service returns a copy of Service config
func (c *config) Service() Service {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return Service{
		Port:       c.service.Port,
		VersionAPI: c.service.VersionAPI,
	}
}

// Logger returns a copy of Logger config
func (c *config) Logger() Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return Logger{
		FilePath: c.logger.FilePath,
		Level:    c.logger.Level,
	}
}

// setService sets Service config
func (c *config) setService(port, versionAPI string) error {
	if port == "" {
		return fmt.Errorf("%w, service port is missing", ErrEmptyConfiguration)
	}

	if versionAPI == "" {
		return fmt.Errorf("%w, version API is missing", ErrEmptyConfiguration)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.service = Service{
		Port:       port,
		VersionAPI: versionAPI,
	}
	return nil
}

// setLogger sets Logger config
func (c *config) setLogger(path, level string) error {
	if path == "" {
		return fmt.Errorf("%w, log filepath is missing", ErrEmptyConfiguration)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger = Logger{
		FilePath: path,
		Level:    level,
	}
	return nil
}
