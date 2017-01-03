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
		env.EmitNew("rater", topic, nil)
		time.Sleep(time.Duration(rand.Intn(3)) * time.Millisecond)
	}

	collected, err := env.Request("collector", cells.CollectedTopic, nil, cells.DefaultTimeout)
	assert.Nil(err)
	events := collected.([]behaviors.EventData)
	assert.True(len(events) <= 10000)
	for _, event := range events {
		assert.Equal(event.Topic, "event-rate!")
		hi, ok := event.Payload.GetDuration(behaviors.EventRateHighPayload)
		assert.True(ok)
		avg, ok := event.Payload.GetDuration(behaviors.EventRateAveragePayload)
		assert.True(ok)
		lo, ok := event.Payload.GetDuration(behaviors.EventRateLowPayload)
		assert.True(ok)
		assert.True(lo <= avg)
		assert.True(avg <= hi)
	}
}

// EOF
