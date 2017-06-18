// Tideland Go Cells - Behaviors - Callback
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
	"github.com/tideland/golib/logger"

	"github.com/tideland/gocells/cells"
)

//--------------------
// CALLBACK BEHAVIOR
//--------------------

// Callbacker is a function called by the behavior when it receives an event.
type Callbacker func(cell cells.Cell, event cells.Event) error

// callbackBehavior is an event processor calling all stored functions
// if it receives an event.
type callbackBehavior struct {
	cell      cells.Cell
	callbacks []Callbacker
}

// NewCallbackBehavior creates a behavior with a number of callback functions.
// Each time an event is received those functions are called in the same order
// they have been passed.
func NewCallbackBehavior(callbacks ...Callbacker) cells.Behavior {
	if len(callbacks) == 0 {
		logger.Errorf("callback created without callback functions")
	}
	return &callbackBehavior{nil, callbacks}
}

// Init the behavior.
func (b *callbackBehavior) Init(cell cells.Cell) error {
	b.cell = cell
	return nil
}

// Terminate the behavior.
func (b *callbackBehavior) Terminate() error {
	return nil
}

// ProcessEvent calls a callback functions with the event data.
func (b *callbackBehavior) ProcessEvent(event cells.Event) error {
	for _, callback := range b.callbacks {
		if err := callback(b.cell, event); err != nil {
			logger.Errorf("callback terminated with error: %v", err)
			return err
		}
	}
	return nil
}

// Recover from an error.
func (b *callbackBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
