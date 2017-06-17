// Tideland Go Cells - Behaviors - Unit Tests - Filter
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

// TestFilterBehavior tests the filter behavior.
func TestFilterBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sigc := audit.MakeSigChan()
	counter := 0
	env := cells.NewEnvironment("filter-behavior")
	defer env.Stop()

	filter := func(event cells.Event) (bool, error) {
		payload := event.Payload().String()
		return event.Topic() == payload, nil
	}
	conditioner := func(event cells.Event) bool {
		counter++
		return counter == 2
	}
	processor := func(cell cells.Cell, event cells.Event) error {
		sigc <- true
		return nil
	}

	env.StartCell("filter", behaviors.NewFilterBehavior(filter))
	env.StartCell("conditioner", behaviors.NewConditionBehavior(conditioner, processor))
	env.Subscribe("filter", "conditioner")

	env.EmitNew("filter", "a", "a")
	env.EmitNew("filter", "a", "b")
	env.EmitNew("filter", "a", "c")
	env.EmitNew("filter", "a", "d")
	env.EmitNew("filter", "b", "b")

	assert.Wait(sigc, true, time.Second)
}

// EOF
