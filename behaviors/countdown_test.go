// Tideland Go Cells - Behaviors - Unit Tests - Countdown
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

// TestCountdownBehavior tests the countdown of events.
func TestCountdownBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("countdown-behavior")
	defer env.Stop()

	zeroer := func(accessor cells.EventSinkAccessor) (cells.Event, int, error) {
		t := accessor.Len()
		event, err := cells.NewEvent("zero", t)
		return event, t - 1, err
	}
	tester := func(event cells.Event) bool {
		return event.Topic() == "zero"
	}
	processor := func(cell cells.Cell, event cells.Event) error {
		sigc <- event.Topic()
		return nil
	}

	env.StartCell("countdowner", behaviors.NewCountdownBehavior(5, zeroer))
	env.StartCell("conditioner", behaviors.NewConditionBehavior(tester, processor))
	env.Subscribe("countdowner", "conditioner")

	countdown := func(t int) {
		assert.Logf("countdown with T = %d", t)
		for i := 0; i < t; i++ {
			err := env.EmitNew("countdowner", "count", i)
			assert.Nil(err)
		}
		assert.Wait(sigc, "zero", time.Second)
	}

	countdown(5)
	countdown(4)
	countdown(3)
	countdown(2)
	countdown(1)
}

// EOF
