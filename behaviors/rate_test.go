// Tideland Go Cells - Behaviors - Unit Tests - Event Rate
//
// Copyright (C) 2010-2017 Frank Mueller / Oldenburg / Germany
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

// TestRateBehavior tests the event rate behavior.
func TestRateBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sigc := audit.MakeSigChan()
	generator := audit.NewGenerator(audit.FixedRand())
	env := cells.NewEnvironment("rate-behavior")
	defer env.Stop()

	matcher := func(event cells.Event) (bool, error) {
		return event.Topic() == "now", nil
	}
	processor := func(accessor cells.EventSinkAccessor) error {
		ok, err := accessor.Match(func(index int, event cells.Event) (bool, error) {
			return event.Topic() == behaviors.TopicRate, nil
		})
		sigc <- ok
		return err
	}
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}

	env.StartCell("rater", behaviors.NewRateBehavior(matcher, 100))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10000, processor))
	env.Subscribe("rater", "collector")

	for i := 0; i < 1000; i++ {
		topic := generator.OneStringOf(topics...)
		env.EmitNew("rater", topic, nil)
		generator.SleepOneOf(0, time.Millisecond, 2*time.Millisecond)
	}

	env.EmitNew("collector", cells.TopicProcess, nil)
	assert.Wait(sigc, true, 10 * time.Second)
}

// EOF
