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

// FSMProcessor is the signature of a function or method which processes
// an event and returns the following status or an error.
type FSMProcessor func(cell cells.Cell, event cells.Event) FSMStatus

// FSMStatus describes the current status of a finite state machine.
// It also contains a reference to the current process function.
type FSMStatus struct {
	Info    string
	Process FSMProcessor
	Error   error
}

// Done returns true if the status contains no processor anymore.
func (s FSMStatus) Done() bool {
	return s.Process == nil || s.Error != nil
}

// FSMInfo contains information about the current status of the FSM.
type FSMInfo struct {
	Info  string
	Done  bool
	Error error
}

// fsmBehavior runs the finite state machine.
type fsmBehavior struct {
	cell   cells.Cell
	status FSMStatus
}

// NewFSMBehavior creates a finite state machine behavior based on the
// passed initial status. The process function is called with the event
// and has to return the next status, which can be the same or a different
// one.
func NewFSMBehavior(status FSMStatus) cells.Behavior {
	return &fsmBehavior{
		status: status,
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
	// Check if done.
	if b.status.Done() {
		return nil
	}
	// Process event and determine next status.
	b.status = b.status.Process(b.cell, event)
	// Emit information.
	b.cell.EmitNew(cells.TopicStatus, FSMInfo{
		Info:  b.status.Info,
		Done:  b.status.Done(),
		Error: b.status.Error,
	})
	return b.status.Error
}

// Recover from an error.
func (b *fsmBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
