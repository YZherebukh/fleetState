package store

import (
	"sync"

	"github.com/fleetState/model"
)

type state struct {
	mu      *sync.RWMutex
	vehicle map[string]model.VehicleDef
}

// New creates new model.Vehicle instance
func New() model.State {
	return &state{
		mu:      &sync.RWMutex{},
		vehicle: make(map[string]model.VehicleDef),
	}
}

// All returns ids of all tracked vehicles
func (s *state) All() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resp := make([]string, 0, len(s.vehicle))

	for v := range s.vehicle {
		resp = append(resp, v)
	}

	return resp
}

// Update updates vehicle information by ID
func (s *state) One(id string) (model.VehicleDef, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.vehicle[id]
	return v, ok
}

// Update updates vehicle information by ID
func (s *state) Update(v model.VehicleDef) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.vehicle[v.ID] = v
}
