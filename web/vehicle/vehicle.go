//go:generate mockgen -source ../vehicle/vehicle.go  -destination ../vehicle/mock/mock_vehicle.go

package vehicle

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/fleetState/model"
	"github.com/fleetState/queue"
	"github.com/fleetState/web"

	"github.com/fleetState/logger"
	"github.com/gorilla/mux"
)

// Middleware is a middleware interface
type Middleware interface {
	SetContextHeader(next http.HandlerFunc) http.HandlerFunc
	SetContextHeaderWithTimeout(next http.HandlerFunc) http.HandlerFunc
}

// package const
const (
	vID = "vid"
)

// Handler is a web events handler struct
type Handler struct {
	ctx        context.Context
	router     *mux.Router
	log        logger.Logger
	middleware Middleware
	resp       *web.Response
	vehicle    *model.Vehicle
	state      queue.State
	stream     queue.Stream
}

// NewHandler creates new vehicle handler instancce
func NewHandler(ctx context.Context, router *mux.Router, l logger.Logger, m Middleware, r *web.Response,
	v *model.Vehicle, state queue.State, stream queue.Stream) {
	h := Handler{
		ctx:        ctx,
		router:     router,
		log:        l,
		middleware: m,
		resp:       r,
		vehicle:    v,
		state:      state,
		stream:     stream,
	}

	apiV1 := h.router.PathPrefix("/v1").Subrouter()

	apiV1.HandleFunc("/vehicle", h.middleware.SetContextHeaderWithTimeout(http.HandlerFunc(h.All))).
		Methods(http.MethodGet)

	apiV1.HandleFunc("/vehicle/{vid}", h.middleware.SetContextHeaderWithTimeout(http.HandlerFunc(h.Update))).
		Methods(http.MethodPost)

	apiV1.HandleFunc("/vehicle/{vid}/stream", http.HandlerFunc(h.Stream)).
		Methods(http.MethodGet)
}

// All handles GET all vehicle state requests
func (h *Handler) All(w http.ResponseWriter, r *http.Request) {
	select {
	case <-h.ctx.Done():
		return
	default:
		vehicles := h.vehicle.All()
		if len(vehicles) == 0 {
			h.resp.SendError(w, r, http.StatusNoContent, http.StatusText(http.StatusNoContent))
			return
		}

		var resp = make([]struct {
			VID string `json:"vehicleID"`
		}, len(vehicles))

		for i := range vehicles {
			resp[i].VID = vehicles[i]
		}

		h.resp.SendJSON(w, r, http.StatusOK, resp)
	}
}

// Update handles POST Update vehicle state requests
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	select {
	case <-h.ctx.Done():
		return
	default:
		ctx := r.Context()

		vars := mux.Vars(r)
		id, ok := vars[vID]

		if !ok {
			h.log.Warningf(ctx, "get vehicleID failed. error: %s", emptyID)
			h.resp.SendError(w, r, http.StatusBadRequest, emptyID)
			return
		}

		var req struct {
			Lat  float64 `json:"lat"`
			Long float64 `json:"lon"`
		}

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			h.log.Errorf(ctx, "failed to read request body, err: %s", err.Error())
			h.resp.SendError(w, r, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		defer func() {
			_ = r.Body.Close()
		}()

		h.state.Write(h.vehicle.Update(id, req.Lat, req.Long, time.Now().UTC()))

		h.resp.SendStatus(w, r, http.StatusOK)
	}
}

// Stream handles GET one vehicle state stream requests
func (h *Handler) Stream(w http.ResponseWriter, r *http.Request) {
	select {
	case <-h.ctx.Done():
		return
	default:
		ctx := r.Context()
		vars := mux.Vars(r)

		id, ok := vars[vID]
		if !ok {
			h.log.Warningf(ctx, "get vehicleID failed. error: %s", emptyID)
			h.resp.SendError(w, r, http.StatusBadRequest, emptyID)
			return
		}

		h.stream.Create(id)
		defer h.state.Delete(id)

		flusher, ok := w.(http.Flusher)
		if !ok {
			h.log.Errorf(ctx, "failed to stream vehicle state")
			h.resp.SendError(w, r, http.StatusHTTPVersionNotSupported, http.StatusText(http.StatusHTTPVersionNotSupported))
			return
		}

		encoder := json.NewEncoder(w)

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-h.stream.Read(id):
				if !ok {
					return
				}
				flusher.Flush()

				var resp struct {
					ID          string  `json:"vehicleID"`
					Latitude    float64 `json:"latitude"`
					Longitude   float64 `json:"longitude"`
					Speed       int     `json:"speed"`
					Measurement string  `json:"measurement"`
				}

				resp.ID = v.ID
				resp.Latitude = v.Latitude
				resp.Longitude = v.Longitude
				resp.Speed = v.Speed.Num
				resp.Measurement = v.Speed.Measurement

				err := encoder.Encode(&resp)
				if err != nil {
					h.log.Errorf(ctx, "failed to stream vehicle state, err: %s", err.Error())
					h.resp.SendError(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
					return
				}

				flusher.Flush()
			}
		}
	}
}
