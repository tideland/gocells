// Tideland Go Cells - Behaviors - Finite State Machine
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// IMPORTS
//--------------------

import (
	"context"

	"github.com/tideland/gocells/cells"
)

//--------------------
// FSM BEHAVIOR
//--------------------

// FSMState is the signature of a function or method which processes
// an event and returns the following state or an error.
type FSMState func(cell cells.Cell, event cells.Event) (FSMState, error)

// fsmStatus contains information about the current status of the FSM.
type fsmStatus struct {
	done bool
	err  error
}

// fsmBehavior runs the finite state machine.
type fsmBehavior struct {
	cell  cells.Cell
	state FSMState
	done  bool
	err   error
}

// NewFSMBehavior creates a finite state machine behavior based on the
// passed initial state function. The function is called with the event
// has to return the next state, which can be the same one. In case of
// nil the stae will be transfered into a generic end state, if an error
// is returned the state is a generic error state.
func NewFSMBehavior(state FSMState) cells.Behavior {
	return &fsmBehavior{nil, state, false, nil}
}

// Init the behavior.
func (b *fsmBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *fsmBehavior) Terminate() error {
	return nil
}

// ProcessEvent executes the state function and stores
// the returned new state.
func (b *fsmBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case cells.StatusTopic:
		waiter, ok := event.Payload().GetWaiter(cells.DefaultPayload)
		if ok {
			response := &fsmStatus{
				done: b.done,
				err:  b.err,
			}
			waiter.Set(cells.NewPayload(response))
		}
	default:
		if b.done {
			return nil
		}
		state, err := b.state(b.cell, event)
		if err != nil {
			b.done = true
			b.err = err
		} else if state == nil {
			b.done = true
		}
		b.state = state
	}
	return nil
}

// Recover from an error.
func (b *fsmBehavior) Recover(err interface{}) error {
	b.done = true
	b.err = cells.NewCannotRecoverError(b.cell.ID(), err)
	return nil
}

// RequestFSMStatus retrieves the status of a FSM cell.
func RequestFSMStatus(env cells.Environment, id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cells.DefaultTimeout)
	defer cancel()
	waiter := cells.NewPayloadWaiter()
	err := env.EmitNewContext(ctx, id, cells.StatusTopic, waiter)
	if err != nil {
		return false, err
	}
	payload, err := waiter.Wait(ctx)
	if err != nil {
		return false, err
	}
	status := payload.Default(nil).(*fsmStatus)
	if status == nil {
		return false, cells.NewInvalidResponseError(payload)
	}
	return status.done, status.err
}

// EOF
