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

// filterMode describes if the filter works selecting or excluding.
type filterMode int

// Flags for the filt
const (
	selectFilter filterMode = iota
	excludeFilter
)

// Filter is a function type checking if an event shall be filtered.
type Filter func(event cells.Event) (bool, error)

// filterBehavior is a simple repeater using the filter
// function to check if an event shall be selected or excluded
// for re-emitting.
type filterBehavior struct {
	cell    cells.Cell
	mode    filterMode
	matches Filter
}

// NewSelectFilterBehavior creates a filter behavior based on the passed function.
// It re-emits every received event for which the filter function returns true.
func NewSelectFilterBehavior(matches Filter) cells.Behavior {
	return &filterBehavior{
		mode:    selectFilter,
		matches: matches,
	}
}

// NewExcludeFilterBehavior creates a filter behavior based on the passed function.
// It re-emits every received event for which the filter function returns false.
func NewExcludeFilterBehavior(matches Filter) cells.Behavior {
	return &filterBehavior{
		mode:    excludeFilter,
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
	switch b.mode {
	case selectFilter:
		// Select those who match.
		if ok {
			return b.cell.Emit(event)
		}
	case excludeFilter:
		// Exclude those who match, emit the others.
		if !ok {
			return b.cell.Emit(event)
		}
	}
	return nil
}

// Recover from an error.
func (b *filterBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
