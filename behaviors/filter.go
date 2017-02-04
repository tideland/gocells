// Tideland Go Cells - Behaviors - Filter
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// IMPORTS
//--------------------

import "github.com/tideland/gocells/cells"

//--------------------
// FILTER BEHAVIOR
//--------------------

// Filter is a function type checking if an event shall be filtered.
type Filter func(event cells.Event) (bool, error)

// filterBehavior is a simple repeater using the filter
// function to check if an event shall be emitted.
type filterBehavior struct {
	cell    cells.Cell
	matches Filter
}

// NewFilterBehavior creates a filter behavior based on the passed function.
// It emits every received event for which the filter function returns true.
func NewFilterBehavior(matches Filter) cells.Behavior {
	return &filterBehavior{
		matches: matches,
	}
}

// Init the behavior.
func (b *filterBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *filterBehavior) Terminate() error {
	return nil
}

// ProcessEvent emits the event when the filter func returns true.
func (b *filterBehavior) ProcessEvent(event cells.Event) error {
	ok, err := b.matches(event)
	if err != nil {
		return err
	}
	if ok {
		b.cell.Emit(event)
	}
	return nil
}

// Recover from an error.
func (b *filterBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
