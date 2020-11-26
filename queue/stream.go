package queue

import (
	"context"
	"sync"

	"github.com/fleetState/model"
)

// Stream is a Stream struct
type Stream struct {
	mu      *sync.RWMutex
	chanMap map[string]chan model.VehicleDef
}

// NewStream creates new Stream instance
func NewStream(ctx context.Context) *Stream {
	return &Stream{
		mu:      &sync.RWMutex{},
		chanMap: make(map[string]chan model.VehicleDef),
	}
}

// Exist is checking if vehicle sate hah being streamed
func (s *Stream) Exist(vid string) bool {
	_, ok := s.chanMap[vid]
	return ok
}

// Create creates new chan model.VehicleDef in state stream map
func (s *Stream) Create(vid string) {
	if _, ok := s.chanMap[vid]; ok {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.chanMap[vid] = make(chan model.VehicleDef)
	return
}

// Update is pushing updated information of model.VehicleDef to a streaming channel
func (s *Stream) update(v model.VehicleDef) {
	c, ok := s.chanMap[v.ID] // if channel with vehicle id exists in this map
	if !ok {                 // that means that information about vehicle is streaming
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	c <- v
}

// Read is returning chan model.VehicleDef by vid to be able to read information from it
func (s *Stream) Read(vid string) chan model.VehicleDef {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.chanMap[vid] // if channel with vehicle id exists in this map
	if !ok {                // that means that information about vehicle can be streamed
		return nil
	}

	return c
}

// Delete deletes chan model.VehicleDef from state stream map
// should be used, if streaming is not required any more
func (s *Stream) delete(vid string) {
	_, ok := s.chanMap[vid]
	if !ok {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.chanMap, vid)
}
