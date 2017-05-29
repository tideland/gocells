// Tideland Go Cells - Behaviors - Waiter
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
	"github.com/tideland/golib/errors"

	"github.com/tideland/gocells/cells"
)

//--------------------
// WAITER BEHAVIOR
//--------------------

// WaiterFunc describes the function called after the first event.
type WaiterFunc func(event cells.Event) error

// waiterBehavior implements the waiter behavior.
type waiterBehavior struct {
	cell   cells.Cell
	waiter WaiterFunc
}

// NewWaiterBehavior creates a behavior where the cell calls the waiter
// function for the first received event.
func NewWaiterBehavior(waiter WaiterFunc) cells.Behavior {
	return &waiterBehavior{
		waiter: waiter,
	}
}

// Init the behavior.
func (b *waiterBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *waiterBehavior) Terminate() error {
	return nil
}

// ProcessEvent calls the simple processor function.
func (b *waiterBehavior) ProcessEvent(event cells.Event) error {
	if b.waiter != nil {
		err := b.waiter(event)
		b.waiter == nil
		return err
	}
	return nil
}

// Recover from an error.
func (b *waiterBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
