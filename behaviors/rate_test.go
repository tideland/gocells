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
	"context"
	"math/rand"
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
	ctx := context.Background()
	env := cells.NewEnvironment("rate-behavior")
	defer env.Stop()

	matches := func(event cells.Event) bool {
		return event.Topic() == "now"
	}
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}

	env.StartCell("rater", behaviors.NewRateBehavior(matches, 100))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10000))
	env.Subscribe("rater", "collector")

	for i := 0; i < 10000; i++ {
		topic := topics[rand.Intn(len(topics))]
		env.EmitNew(ctx, "rater", topic, nil)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Millisecond)
	}

	accessor, err := behaviors.RequestCollectedAccessor(env, "collector", cells.DefaultTimeout)
	assert.Nil(err)
	assert.True(accessor.Len() <= 10000)
	err = accessor.Do(func(index int, event cells.Event) error {
		assert.Equal(event.Topic(), "event-rate!")
		hi := event.Payload().GetDuration(behaviors.PayloadRateHigh, -1)
		avg := event.Payload().GetDuration(behaviors.PayloadRateAverage, -1)
		lo := event.Payload().GetDuration(behaviors.PayloadRateLow, -1)
		assert.True(lo <= avg)
		assert.True(avg <= hi)
		return nil
	})
	assert.Nil(err)
}

// EOF
