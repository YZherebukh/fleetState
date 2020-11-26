//go:generate mockgen -source ../logger/logger.go  -destination ../logger/mock/mock_logger.go

package logger

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/fleetState/config"
	cont "github.com/fleetState/context"
)

type logger struct {
	log *log.Logger
}

// Logger is a logger interface
type Logger interface {
	Errorf(ctx context.Context, format string, v ...interface{})
	Infof(ctx context.Context, format string, v ...interface{})
	Warningf(ctx context.Context, format string, v ...interface{})
	Fatalf(ctx context.Context, format string, v ...interface{})
}

// New creates a new logger
func New(cfg config.Logger) Logger {
	l := logrus.New()
	level, err := logrus.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		level = logrus.ErrorLevel
	}
	l.Level = level

	lf, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		log.Errorf("Failed to open log file: %v", err)
		return nil
	}
	l.Out = lf

	return logger{log: l}
}

// Errorf logs with the Error severity.
// Arguments are handled in the manner of fmt.Printf.
// firs parameter is context. Logger will try go get processID coried by context
// so all logs for same request will have same processID
func (l logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.log.Errorf(fmt.Sprintf("processID:%s, %s", cont.ProcessID(ctx), format), v...)
}

// Infof logs with the Info severity.
// Arguments are handled in the manner of fmt.Printf.
// firs parameter is context. Logger will try go get processID coried by context
// so all logs for same request will have same processID
func (l logger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.log.Infof(fmt.Sprintf("processID:%s, %s", cont.ProcessID(ctx), format), v...)
}

// Warningf logs with the Warning severity.
// Arguments are handled in the manner of fmt.Printf.
// firs parameter is context. Logger will try go get processID coried by context
// so all logs for same request will have same processID
func (l logger) Warningf(ctx context.Context, format string, v ...interface{}) {
	l.log.Warningf(fmt.Sprintf("processID:%s, %s", cont.ProcessID(ctx), format), v...)
}

// Fatalf logs with the Fatal severity, and ends with os.Exit(1).
// Arguments are handled in the manner of fmt.Printf.
// firs parameter is context. Logger will try go get processID coried by context
// so all logs for same request will have same processID
func (l logger) Fatalf(ctx context.Context, format string, v ...interface{}) {
	l.log.Fatalf(fmt.Sprintf("processID:%s, %s", cont.ProcessID(ctx), format), v...)
}
