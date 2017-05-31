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
// and has to return the next state, which can be the same or a different
// one. In case of nil the state will be transferred into a generic end
// state, if an error is returned the state is a generic error state.
func NewFSMBehavior(state FSMState) cells.Behavior {
	return &fsmBehavior{
		state: state,
		done:  false,
	}
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
	if b.done {
		return nil
	}
	// Determine next state.
	state, err := b.state(b.cell, event)
	if err != nil {
		b.done = true
		b.err = err
	} else if state == nil {
		b.done = true
	}
	b.state = state
	// Emit status.
	b.cell.EmitNew(cells.TopicStatus, cells.PayloadValues{
		cells.PayloadDone:  b.done,
		cells.PayloadError: b.err,
	}.Payload())
	return nil
}

// Recover from an error.
func (b *fsmBehavior) Recover(err interface{}) error {
	b.done = true
	b.err = cells.NewCannotRecoverError(b.cell.ID(), err)
	return nil
}

// EOF
