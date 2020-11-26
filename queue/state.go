package queue

import (
	"context"

	"github.com/fleetState/model"
)

// State is a state State struct
type State struct {
	ctx context.Context

	addWriteChan  chan model.VehicleDef
	addDeleteChan chan string

	popWriteChan  chan struct{}
	popDeleteChan chan struct{}

	stream Stream
}

// NewState creates new State instance
// it also runs 2 gorouteenes to listen from add channels
func NewState(ctx context.Context, s Stream) *State {
	state := &State{
		ctx:           ctx,
		addWriteChan:  make(chan model.VehicleDef, 1),
		addDeleteChan: make(chan string, 1),
		popWriteChan:  make(chan struct{}, 100),
		popDeleteChan: make(chan struct{}, 100),

		stream: s,
	}

	go state.popWrite()
	go state.popDelete()

	return state
}

// Write writes vid into write chan
func (s *State) Write(vehicle model.VehicleDef) {
	if !s.stream.Exist(vehicle.ID) {
		return
	}

	select {
	case <-s.ctx.Done():
		return
	case s.addWriteChan <- vehicle:
		return
	}
}

// popCreate poping from queue and starts go routines
func (s *State) popWrite() {
	for {
		select {
		case vehicle, ok := <-s.addWriteChan:
			if !ok {
				s.addWriteChan = nil
				return
			}
			s.popWriteChan <- struct{}{}

			go func() {
				s.stream.update(vehicle)
				<-s.popWriteChan
			}()
		}
	}
}

// Delete writes vid into write chan
func (s *State) Delete(vid string) {
	select {
	case <-s.ctx.Done():
		return
	case s.addDeleteChan <- vid:
		return
	}
}

// popDelete poping from queue and starts go routines
func (s *State) popDelete() {
	for {
		select {
		case vid, ok := <-s.addDeleteChan:
			if !ok {
				s.addWriteChan = nil
				return
			}
			s.popDeleteChan <- struct{}{}

			go func() {
				s.stream.delete(vid)
				<-s.popDeleteChan
			}()
		}
	}
}
