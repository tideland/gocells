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

// ProcessCallback is a function called by the behavior when it receives an event.
type ProcessCallback func(topic string, payload cells.Payload) error

// callbackBehavior is an event processor calling all stored functions
// if it receives an event.
type callbackBehavior struct {
	cell      cells.Cell
	callbacks []ProcessCallback
}

// NewCallbackBehavior creates a behavior with a number of callback functions.
// Each time an event is received those functions are called in the same order
// they have been passed.
func NewCallbackBehavior(cbfs ...ProcessCallback) cells.Behavior {
	if len(cbfs) == 0 {
		logger.Errorf("callback created without callback functions")
	}
	return &callbackBehavior{nil, cbfs}
}

// Init the behavior.
func (b *callbackBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *callbackBehavior) Terminate() error {
	return nil
}

// ProcessEvent calls a callback functions with the event data.
func (b *callbackBehavior) ProcessEvent(event cells.Event) error {
	for _, callback := range b.callbacks {
		if err := callback(event.Topic(), event.Payload()); err != nil {
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
