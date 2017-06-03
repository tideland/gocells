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
		return cells.Values{
			cells.PayloadDefault: counters,
		}.Payload()
	}
	counter := func(event cells.Event) []string {
		increments, ok := event.Payload().GetStringSlice(cells.PayloadDefault)
		if !ok {
			return nil
		}
		return increments
	}
	processor := func(accessor cells.EventSinkAccessor) error {
		if accessor.Len() == 3 {
			last, err := accessor.PeekLast()
			assert.Nil(err)
			values := cells.Values{}
			last.Payload().Do(func(key, value string) error {
				values[key] = value
				return nil
			})
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

	assert.Wait(sigc, cells.Values{
		"a": "3",
		"b": "1",
		"c": "1",
		"d": "2",
	}, time.Second)
}

// EOF
