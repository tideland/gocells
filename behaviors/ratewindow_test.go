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
	generator := audit.NewGenerator(audit.FixedRand())
	env := cells.NewEnvironment("rate-window-behavior")
	defer env.Stop()

	matches := func(event cells.Event) bool {
		return event.Topic() == "now"
	}
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}
	duration := 10 * time.Millisecond

	env.StartCell("windower", behaviors.NewRateWindowBehavior(matches, 5, duration))
	env.StartCell("collector", behaviors.NewCollectorBehavior(100))
	env.Subscribe("windower", "collector")

	for i := 0; i < 10; i++ {
		// Slow loop.
		for j := 0; j < 100; j++ {
			topic := generator.OneStringOf(topics...)
			env.EmitNew("windower", topic, nil)
			time.Sleep(1)
		}
		// Fast loop.
		for j := 0; j < 100; j++ {
			topic := generator.OneStringOf(topics...)
			env.EmitNew("windower", topic, nil)
		}
	}

	collected, err := env.Request("collector", cells.CollectedTopic, nil, cells.DefaultTimeout)
	assert.Nil(err)
	events := collected.([]behaviors.EventData)
	assert.Logf("Window Events: %d", len(events))
	assert.True(len(events) >= 1)

	for _, event := range events {
		count, ok := event.Payload.GetInt(behaviors.EventRateWindowCountPayload)
		assert.True(ok)
		assert.Equal(count, 5)
		first, ok := event.Payload.GetTime(behaviors.EventRateWindowFirstTimePayload)
		assert.True(ok)
		last, ok := event.Payload.GetTime(behaviors.EventRateWindowLastTimePayload)
		assert.True(ok)
		difference := last.Sub(first)
		assert.True(difference <= duration)
	}
}

// EOF
