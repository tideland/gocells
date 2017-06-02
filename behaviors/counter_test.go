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
		// TODO 2017-06-02 Mue Analzye collected counters.
		return nil
	}

	env.StartCell("counter", behaviors.NewCounterBehavior(cf))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10, processor))
	env.Subscribe("counter", "collector")

	env.EmitNew("counter", "count", counters("a", "b"))
	env.EmitNew("counter", "count", counters("a", "c", "d"))
	env.EmitNew("counter", "count", counters("a", "d"))

	env.EmitNew("collector", cells.TopicProcess, nil)

	// TODO 2017-06-02 Mue Check resetting the counters.
}

// EOF
