// Tideland Go Cells - Behaviors - Unit Tests - Counter
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

// TestCounterBehavior tests the counting of events.
func TestCounterBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("counter-behavior")
	defer env.Stop()

	mkcounters := func(counters ...string) cells.Payload {
		payload, err := cells.NewPayload(counters)
		assert.Nil(err)
		return payload
	}
	counter := func(event cells.Event) []string {
		var increments []string
		err := event.Payload().Unmarshal(&increments)
		assert.Nil(err)
		return increments
	}
	conditioner := func(event cells.Event) bool {
		var values map[string]uint
		err := event.Payload().Unmarshal(&values)
		assert.Nil(err)
		return values["a"] == 3 && values["b"] == 1 && values["c"] == 1 && values["d"] == 2
	}
	processor := func(cell cells.Cell, event cells.Event) error {
		sigc <- true
		return nil
	}

	env.StartCell("counter", behaviors.NewCounterBehavior(counter))
	env.StartCell("conditioner", behaviors.NewConditionBehavior(conditioner, processor))
	env.Subscribe("counter", "conditioner")

	env.EmitNew("counter", "count", mkcounters("a", "b"))
	env.EmitNew("counter", "count", mkcounters("a", "c", "d"))
	env.EmitNew("counter", "count", mkcounters("a", "d"))

	assert.Wait(sigc, true, time.Second)
}

// EOF
