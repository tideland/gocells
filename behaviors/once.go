// Tideland Go Cells - Behaviors - Once
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
// ONCE BEHAVIOR
//--------------------

// OneTimer describes the function called after the first event.
type OneTimer func(cell cells.Cell, event cells.Event) error

// onceBehavior implements the oneTimer behavior.
type onceBehavior struct {
	cell     cells.Cell
	oneTimer OneTimer
}

// NewOnceBehavior creates a behavior where the cell calls the one-timer
// function for the first received event. Afterwards it will never be called
// again.
func NewOnceBehavior(oneTimer OneTimer) cells.Behavior {
	return &onceBehavior{
		oneTimer: oneTimer,
	}
}

// Init the behavior.
func (b *onceBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *onceBehavior) Terminate() error {
	return nil
}

// ProcessEvent callsthe one-timer.
func (b *onceBehavior) ProcessEvent(event cells.Event) error {
	if b.oneTimer != nil {
		err := b.oneTimer(b.cell, event)
		b.oneTimer = nil
		return err
	}
	return nil
}

// Recover from an error.
func (b *onceBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
