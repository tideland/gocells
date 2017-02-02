// Tideland Go Cells - Behaviors - Unit Tests - Event Pair
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

// TestPairBehavior tests the event pair behavior.
func TestPairBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	ctx := context.Background()
	generator := audit.NewGenerator(audit.FixedRand())
	env := cells.NewEnvironment("pair-behavior")
	defer env.Stop()

	matches := func(event cells.Event, data interface{}) (interface{}, bool) {
		if event.Topic() == "now" {
			now := time.Now().Unix()
			return now, true
		}
		return nil, false
	}
	filterFuncBuilder := func(positive bool) behaviors.FilterFunc {
		var topic string
		if positive {
			topic = behaviors.TopicPair
		} else {
			topic = behaviors.TopicPairTimeout
		}
		return func(id string, event cells.Event) bool {
			return event.Topic() == topic
		}
	}
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}
	duration := time.Millisecond

	env.StartCell("pairer", behaviors.NewPairBehavior(matches, duration))
	env.StartCell("positive-filter", behaviors.NewFilterBehavior(filterFuncBuilder(true)))
	env.StartCell("negative-filter", behaviors.NewFilterBehavior(filterFuncBuilder(false)))
	env.StartCell("positive-collector", behaviors.NewCollectorBehavior(1000))
	env.StartCell("negative-collector", behaviors.NewCollectorBehavior(1000))
	env.Subscribe("pairer", "positive-filter", "negative-filter")
	env.Subscribe("positive-filter", "positive-collector")
	env.Subscribe("negative-filter", "negative-collector")

	for i := 0; i < 5000; i++ {
		topic := generator.OneStringOf(topics...)
		env.EmitNew(ctx, "pairer", topic, nil)
		generator.SleepOneOf(0, time.Millisecond, 2*time.Millisecond)
	}

	accessor, err := behaviors.RequestCollectedAccessor(env, "positive-collector", cells.DefaultTimeout)
	assert.Nil(err)
	assert.True(accessor.Len() >= 1)
	assert.Logf("Positive Events: %d", accessor.Len())

	err = accessor.Do(func(index int, event cells.Event) error {
		first := event.Payload().GetTime(behaviors.PayloadPairFirstTime, time.Time{})
		second := event.Payload().GetTime(behaviors.PayloadPairSecondTime, time.Time{})
		difference := second.Sub(first)
		assert.False(first.IsZero())
		assert.False(second.IsZero())
		assert.True(difference <= duration)
		return nil
	})

	accessor, err = behaviors.RequestCollectedAccessor(env, "negative-collector", cells.DefaultTimeout)
	assert.Nil(err)
	assert.True(accessor.Len() >= 1)
	assert.Logf("Negative Events: %d", accessor.Len())

	err = accessor.Do(func(index int, event cells.Event) error {
		first := event.Payload().GetTime(behaviors.PayloadPairFirstTime, time.Time{})
		timeout := event.Payload().GetTime(behaviors.PayloadPairTimeout, time.Time{})
		difference := timeout.Sub(first)
		assert.False(first.IsZero())
		assert.False(timeout.IsZero())
		assert.True(difference > duration)
		return nil
	})
	assert.Nil(err)
}

// EOF
