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
	"testing"
	"time"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestPairBehavior tests the event rate behavior.
func TestPairBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
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
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}

	env.StartCell("pairer", behaviors.NewPairBehavior(matches, time.Second))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10000))
	env.Subscribe("pairer", "collector")

	for i := 0; i < 10000; i++ {
		topic := generator.OneStringOf(topics...)
		env.EmitNew("pairer", topic, nil)
		pause := time.Duration(generator.OneIntOf(0, 1, 2)) * time.Millisecond
		time.Sleep(pause)
	}

	collected, err := env.Request("collector", cells.CollectedTopic, nil, cells.DefaultTimeout)
	assert.Nil(err)
	events := collected.([]behaviors.EventData)
	assert.True(len(events) >= 1)
	assert.Logf("Events: %d", len(events))
}

// EOF
