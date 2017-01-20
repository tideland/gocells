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
	boringTopics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	interestingTopics := []string{"a", "b", "c", "d", "now"}
	duration := 25 * time.Millisecond

	env.StartCell("windower", behaviors.NewRateWindowBehavior(matches, 5, duration))
	env.StartCell("collector", behaviors.NewCollectorBehavior(100))
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

	collected, err := env.Request("collector", cells.CollectedTopic, nil, cells.DefaultTimeout)
	assert.Nil(err)
	events, ok := collected.(*cells.EventDatas)
	assert.True(ok)
	assert.Logf("Window Events: %d", events.Len())
	assert.True(events.Len() >= 1)

	err = events.Do(func(index int, data *cells.EventData) error {
		count, ok := data.Payload.GetInt(behaviors.EventRateWindowCountPayload)
		assert.True(ok)
		assert.Equal(count, 5)
		first, ok := data.Payload.GetTime(behaviors.EventRateWindowFirstTimePayload)
		assert.True(ok)
		last, ok := data.Payload.GetTime(behaviors.EventRateWindowLastTimePayload)
		assert.True(ok)
		difference := last.Sub(first)
		assert.True(difference <= duration)
		return nil
	})
	assert.Nil(err)
}

// EOF
