// Tideland Go Cells - Behaviors - Condition
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
// CONDITION BEHAVIOR
//--------------------

// ConditionTester checks if an event matches a wanted state.
type ConditionTester func(event cells.Event) bool

// ConditionProcessor handles the matching event.
type ConditionProcessor func(cell cells.Cell, event cells.Event) error

// conditionBehavior implements the condition behavior.
type conditionBehavior struct {
	cell    cells.Cell
	test    ConditionTester
	process ConditionProcessor
}

// NewConditionBehavior creates a behavior testing of a cell
// fullfills a given condition. If the test returns true the
// processor is called.
func NewConditionBehavior(tester ConditionTester, processor ConditionProcessor) cells.Behavior {
	return &conditionBehavior{
		test:    tester,
		process: processor,
	}
}

// Init the behavior.
func (b *conditionBehavior) Init(cell cells.Cell) error {
	b.cell = cell
	return nil
}

// Terminate the behavior.
func (b *conditionBehavior) Terminate() error {
	return nil
}

// ProcessEvent callsthe one-timer.
func (b *conditionBehavior) ProcessEvent(event cells.Event) error {
	if b.test(event) {
		return b.process(b.cell, event)
	}
	return nil
}

// Recover from an error.
func (b *conditionBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
