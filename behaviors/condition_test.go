// Tideland Go Cells - Behaviors - Unit Tests - Condition
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestConditionBehavior tests the condition behavior.
func TestConditionBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	generator := audit.NewGenerator(audit.FixedRand())
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("condition-behavior")
	defer env.Stop()

	tester := func(event cells.Event) bool {
		return event.Topic() == "now"
	}
	processor := func(cell cells.Cell, event cells.Event) error {
		sigc <- cell.ID() + " " + event.Topic()
		return nil
	}

	env.StartCell("condition", behaviors.NewConditionBehavior(tester, processor))

	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}
	go func() {
		for i := 0; i < 50; i++ {
			topic := generator.OneStringOf(topics...)
			env.EmitNew("condition", topic, nil)
		}
	}()

	assert.Wait(sigc, "condition now", time.Second)
}

// EOF
