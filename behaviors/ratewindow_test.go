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
	"context"
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
	ctx := context.Background()
	generator := audit.NewGenerator(audit.FixedRand())
	env := cells.NewEnvironment("rate-window-behavior")
	defer env.Stop()

	matches := func(event cells.Event) bool {
		return event.Topic() == "now"
	}
	boringTopics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	interestingTopics := []string{"a", "b", "c", "d", "now"}
	duration := 25 * time.Millisecond

	sink := cells.NewEventSink(100)

	env.StartCell("windower", behaviors.NewRateWindowBehavior(matches, 5, duration))
	env.StartCell("collector", behaviors.NewCollectorBehavior(sink))
	env.Subscribe("windower", "collector")

	for i := 0; i < 10; i++ {
		// Slow loop.
		for j := 0; j < 100; j++ {
			topic := generator.OneStringOf(boringTopics...)
			env.EmitNew(ctx, "windower", topic, nil)
			time.Sleep(1)
		}
		// Fast loop.
		for j := 0; j < 10; j++ {
			topic := generator.OneStringOf(interestingTopics...)
			env.EmitNew(ctx, "windower", topic, nil)
		}
	}

	accessor, err := behaviors.RequestCollectedAccessor(ctx, env, "collector", cells.DefaultTimeout)
	assert.Nil(err)
	assert.True(accessor.Len() >= 1)
	assert.Logf("Window Events: %d", accessor.Len())

	err = accessor.Do(func(index int, event cells.Event) error {
		count := event.Payload().GetInt(behaviors.EventRateWindowCountPayload, -1)
		assert.Equal(count, 5)
		first := event.Payload().GetTime(behaviors.EventRateWindowFirstTimePayload, time.Time{})
		last := event.Payload().GetTime(behaviors.EventRateWindowLastTimePayload, time.Time{})
		difference := last.Sub(first)
		assert.True(difference <= duration)
		return nil
	})
	assert.Nil(err)
}

// EOF
