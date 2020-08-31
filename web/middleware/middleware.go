package middleware

import (
	"context"
	"net/http"
	"time"

	cont "github.com/fleetState/contextvalue"
	"github.com/fleetState/logger"
)

// Middleware is a middleware interface
type Middleware interface {
	SetContextHeader(next http.HandlerFunc) http.HandlerFunc
	SetContextHeaderWithTimeout(next http.HandlerFunc) http.HandlerFunc
}

// New creates new Middleware
func New(log logger.Logger) Middleware {
	return &middleware{log: log}
}

type middleware struct {
	log logger.Logger
}

// SetContextHeader is a middlware, that is setting a requestID into r.Context()
// if one is missing, and sets a timeout request to 1 minute
func (m *middleware) SetContextHeader(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := cont.ProcessID(ctx)
		if len(id) == 0 {
			newCtx, err := cont.SetProcessID(ctx)
			if err != nil {
				m.log.Warningf(ctx, "failed to set up request id %v", err.Error())
			}

			m.log.Infof(newCtx, "request: %s:%s received.", r.Method, r.URL.String())
			next.ServeHTTP(w, r.WithContext(newCtx))
			return
		}

		m.log.Infof(ctx, "request %s:%s received.", r.Method, r.URL.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// SetContextHeaderWithTimeout is like SetContextHeader but it also adds one minute timeout
func (m *middleware) SetContextHeaderWithTimeout(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Minute)
		defer cancel()

		id := cont.ProcessID(ctx)
		if len(id) == 0 {
			newCtx, err := cont.SetProcessID(ctx)
			if err != nil {
				m.log.Warningf(ctx, "failed to set up request id %v", err.Error())
			}

			m.log.Infof(newCtx, "request: %s:%s received.", r.Method, r.URL.String())
			next.ServeHTTP(w, r.WithContext(newCtx))
			return
		}

		m.log.Infof(ctx, "request %s:%s received.", r.Method, r.URL.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
