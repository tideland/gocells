// Tideland Go Cells - Behaviors - Status
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
// STATUS BEHAVIOR
//--------------------

// StatusProcessor is called when a cells transmits its status to a
// status cell.
type StatusProcessor func(event cells.Event) error

// statusBehavior implements the combo behavior.
type statusBehavior struct {
	cell    cells.Cell
	process StatusProcessor
}

// NewStatusBehavior returns a behavior for the receiving and processing
// of cell status event. Many behaviors react on the topic "status" and interpret
// the payload as the ID of a status cell. After emitting the status to the
// status cell that one can process the status.
func NewStatusBehavior(processor StatusProcessor) cells.Behavior {
	return &statusBehavior{
		process: processor,
	}
}

// Init implements the cells.Behavior interface.
func (b *statusBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate implements the cells.Behavior interface.
func (b *statusBehavior) Terminate() error {
	return nil
}

// ProcessEvent implements the cells.Behavior interface.
func (b *statusBehavior) ProcessEvent(event cells.Event) error {
	return b.process(event)
}

// Recover implements the cells.Behavior interface.
func (b *statusBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
