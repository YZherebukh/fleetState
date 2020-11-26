package def

import (
	"fmt"
	"os"
	"sync"

	"github.com/fleetState/config"
)

// package const
const (
	port        = "port"
	versionAPI  = "versionAPI"
	logFilePath = "logFilePath"
	logLevel    = "logLevel"
)

// DefaultConfig is a struct with concurent safe public method to access a config
type DefaultConfig struct {
	mu      *sync.RWMutex
	service config.Service
	logger  config.Logger
}

// New initiates a new Configuration instance
func New() (*DefaultConfig, error) {
	config := &DefaultConfig{
		mu: &sync.RWMutex{},
	}

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
func (c *DefaultConfig) Service() config.Service {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return config.Service{
		Port:       c.service.Port,
		VersionAPI: c.service.VersionAPI,
	}
}

// Logger returns a copy of Logger config
func (c *DefaultConfig) Logger() config.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return config.Logger{
		FilePath: c.logger.FilePath,
		Level:    c.logger.Level,
	}
}

// setService sets Service config
func (c *DefaultConfig) setService(port, versionAPI string) error {
	if port == "" {
		return fmt.Errorf("%w, service port is missing", config.ErrEmptyConfiguration)
	}

	if versionAPI == "" {
		return fmt.Errorf("%w, version API is missing", config.ErrEmptyConfiguration)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.service = config.Service{
		Port:       port,
		VersionAPI: versionAPI,
	}
	return nil
}

// setLogger sets Logger config
func (c *DefaultConfig) setLogger(path, level string) error {
	if path == "" {
		return fmt.Errorf("%w, log filepath is missing", config.ErrEmptyConfiguration)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger = config.Logger{
		FilePath: path,
		Level:    level,
	}
	return nil
}
