package store

import (
	"sync"

	"github.com/fleetState/model"
)

// State is a model.VehicleDef in memory store struct
type State struct {
	mu      *sync.RWMutex
	vehicle map[string]model.VehicleDef
}

// New creates new model.Vehicle instance
func New() *State {
	return &State{
		mu:      &sync.RWMutex{},
		vehicle: make(map[string]model.VehicleDef),
	}
}

// All returns ids of all tracked vehicles
func (s *State) All() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resp := make([]string, 0, len(s.vehicle))

	for v := range s.vehicle {
		resp = append(resp, v)
	}

	return resp
}

// One returnes one model.VehicleDef by ID
func (s *State) One(id string) (model.VehicleDef, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.vehicle[id]
	return v, ok
}

// Update updates vehicle information by ID
func (s *State) Update(v model.VehicleDef) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.vehicle[v.ID] = v
}
