package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fleetState/logger"
)

// Response is a web request response struct
type Response struct {
	log logger.Logger
}

// NewResponse creates a new Response instance
func NewResponse(log logger.Logger) *Response {
	return &Response{log: log}
}

// SendJSON send response with JSON body
func (resp *Response) SendJSON(w http.ResponseWriter, r *http.Request, status int, response interface{}) {
	data, err := json.Marshal(response)
	if err != nil {
		resp.log.Errorf(r.Context(), "marshal response failed. error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(fmt.Sprintf("common.renderJSON: err:%s", err)))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(data)

	resp.log.Infof(r.Context(), "request %s:%s success. status: %d", r.Method, r.URL.String(), status)
}

// SendStatus sends status code as a response to thhp request
func (resp *Response) SendStatus(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	resp.log.Infof(r.Context(), "request %s:%s success. status: %d", r.Method, r.URL.String(), status)
}

// SendError sends a defined string as an error message
// with appropriate headers to the HTTP response
func (resp *Response) SendError(w http.ResponseWriter, r *http.Request, status int, message string) {
	resp.log.Errorf(r.Context(), "request %s,%s failed. status: %s, message: %s %s",
		r.Method, r.URL.String(), status, http.StatusText(status), message)

	http.Error(w, message, status)
}
