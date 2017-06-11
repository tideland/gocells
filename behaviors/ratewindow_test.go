// Tideland Go Cells - Behaviors - Unit Tests - Event Rate Window
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

// TestRateWindowBehavior tests the event rate window behavior.
func TestRateWindowBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sigc := audit.MakeSigChan()
	generator := audit.NewGenerator(audit.FixedRand())
	env := cells.NewEnvironment("rate-window-behavior")
	defer env.Stop()

	matcher := func(event cells.Event) (bool, error) {
		return event.Topic() == "now", nil
	}
	processor := func(accessor cells.EventSinkAccessor) error {
		sigc := accessor.Len()
		return nil
	}
	boringTopics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	interestingTopics := []string{"a", "b", "c", "d", "now"}
	duration := 25 * time.Millisecond

	env.StartCell("windower", behaviors.NewRateWindowBehavior(matcher, 5, duration))
	env.StartCell("collector", behaviors.NewCollectorBehavior(100, processor))
	env.Subscribe("windower", "collector")

	for i := 0; i < 10; i++ {
		// Slow loop.
		for j := 0; j < 100; j++ {
			topic := generator.OneStringOf(boringTopics...)
			env.EmitNew("windower", topic, nil)
			time.Sleep(1)
		}
		// Fast loop.
		for j := 0; j < 10; j++ {
			topic := generator.OneStringOf(interestingTopics...)
			env.EmitNew("windower", topic, nil)
		}
	}

	env.EmitNew("collector", cells.TopicProcess, nil)
	assert.Wait(sigc, 10, time.Minute)
}

// EOF
