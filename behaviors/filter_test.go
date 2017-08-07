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
	selects := 0
	excludes := 0
	env := cells.NewEnvironment("filter-behavior")
	defer env.Stop()

	filter := func(event cells.Event) (bool, error) {
		payload := event.Payload().String()
		return event.Topic() == payload, nil
	}
	selectConditioner := func(event cells.Event) bool {
		selects++
		return selects == 2
	}
	excludesConditioner := func(event cells.Event) bool {
		excludes++
		return excludes == 4
	}
	processor := func(cell cells.Cell, event cells.Event) error {
		sigc <- true
		return nil
	}

	env.StartCell("select", behaviors.NewSelectFilterBehavior(filter))
	env.StartCell("selects", behaviors.NewConditionBehavior(selectConditioner, processor))
	env.StartCell("exclude", behaviors.NewExcludeFilterBehavior(filter))
	env.StartCell("excludes", behaviors.NewConditionBehavior(excludesConditioner, processor))

	env.Subscribe("select", "selects")
	env.Subscribe("exclude", "excludes")

	env.EmitNew("select", "a", "a")
	env.EmitNew("select", "a", "b")
	env.EmitNew("select", "a", "c")
	env.EmitNew("select", "a", "d")
	env.EmitNew("select", "b", "b")
	env.EmitNew("select", "b", "a")

	assert.Wait(sigc, true, time.Second)

	env.EmitNew("exclude", "a", "a")
	env.EmitNew("exclude", "a", "b")
	env.EmitNew("exclude", "a", "c")
	env.EmitNew("exclude", "a", "d")
	env.EmitNew("exclude", "b", "b")
	env.EmitNew("exclude", "b", "a")

	assert.Wait(sigc, true, time.Second)
}

// EOF
