// Tideland Go Cells - Behaviors - Unit Tests - Event Sequence
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

// TestSequenceBehavior tests the event sequence behavior.
func TestSequenceBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	generator := audit.NewGenerator(audit.FixedRand())
	env := cells.NewEnvironment("sequence-behavior")
	defer env.Stop()

	matches := func(event cells.Event, events *behaviors.EventDatas) (bool, bool) {
		switch event.Topic() {
		case "c":
			return events.Len() == 0, false
		case "d":
			last, _ := events.Last()
			return events.Len() == 1 && last.Topic == "c", false
		case "e":
			last, _ := events.Last()
			if events.Len() == 2 && last.Topic == "d" {
				return true, true
			}
		}
		return false, false
	}
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}

	env.StartCell("sequencer", behaviors.NewSequenceBehavior(matches))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10000))
	env.Subscribe("sequencer", "collector")

	for i := 0; i < 10000; i++ {
		topic := generator.OneStringOf(topics...)
		env.EmitNew("sequencer", topic, nil)
		generator.SleepOneOf(1*time.Millisecond, 2*time.Millisecond, 3*time.Millisecond)
	}

	collected, err := env.Request("collector", cells.CollectedTopic, nil, cells.DefaultTimeout)
	assert.Nil(err)
	events, ok := collected.(*behaviors.EventDatas)
	assert.True(ok)
	assert.True(events.Len() <= 10000)
	assert.Logf("Sequences: %d", events.Len())
	err = events.Do(func(index int, data *behaviors.EventData) error {
		assert.Equal(data.Topic, behaviors.EventSequenceTopic)
		return nil
	})
	assert.Nil(err)
}

// EOF
