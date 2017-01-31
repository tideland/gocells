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

// waiterBehavior implements the waiter behavior.
type waiterBehavior struct {
	cell   cells.Cell
	waiter cells.PayloadWaiter
}

// NewWaiterBehavior creates a behavior where the cell stores the payload
// of the first received event in the passed waiter.
func NewWaiterBehavior(waiter cells.PayloadWaiter) cells.Behavior {
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
	if b.waiter == nil {
		return errors.New(ErrMissingPayloadWaiter, errorMessages)
	}
	b.waiter.Set(event.Payload())
	return nil
}

// Recover from an error.
func (b *waiterBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
