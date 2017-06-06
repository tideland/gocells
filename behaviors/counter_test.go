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

	counters := func(counters ...string) cells.Payload {
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
	processor := func(accessor cells.EventSinkAccessor) error {
		if accessor.Len() == 3 {
			last, ok := accessor.PeekLast()
			assert.True(ok)
			var values map[string]uint
			err := last.Payload().Unmarshal(&values)
			assert.Nil(err)
			sigc <- values
		}
		return nil
	}

	env.StartCell("counter", behaviors.NewCounterBehavior(counter))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10, processor))
	env.Subscribe("counter", "collector")

	env.EmitNew("counter", "count", counters("a", "b"))
	env.EmitNew("counter", "count", counters("a", "c", "d"))
	env.EmitNew("counter", "count", counters("a", "d"))

	env.EmitNew("collector", cells.TopicProcess, nil)

	assert.Wait(sigc, map[string]uint{
		"a": 3,
		"b": 1,
		"c": 1,
		"d": 2,
	}, time.Second)
}

// EOF
