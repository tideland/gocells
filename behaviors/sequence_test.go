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

	matches := func(event cells.Event, datas *behaviors.EventDatas) (bool, bool) {
		sequence := []string{"a", "e", "now"}
		matcher := func(index int, data *behaviors.EventData) (bool, error) {
			ok := data.Topic == sequence[index]
			return ok, nil
		}
		matches, err := datas.Match(matcher)
		if err != nil || !matches {
			return false, false
		}
		if datas.Len() == len(sequence)-1 && event.Topic() == sequence[len(sequence)-1] {
			return true, true
		}
		return true, false
	}
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}

	env.StartCell("sequencer", behaviors.NewSequenceBehavior(matches))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10000))
	env.Subscribe("sequencer", "collector")

	for i := 0; i < 10000; i++ {
		topic := generator.OneStringOf(topics...)
		env.EmitNew("sequencer", topic, nil)
		generator.SleepOneOf(0, 1*time.Millisecond, 2*time.Millisecond)
	}

	collectedRaw, err := env.Request("collector", cells.CollectedTopic, nil, cells.DefaultTimeout)
	assert.Nil(err)
	collected, ok := collectedRaw.(*behaviors.EventDatas)
	assert.True(ok)
	assert.True(collected.Len() > 0)
	assert.Logf("Collected Sequences: %d", collected.Len())
	err = collected.Do(func(index int, data *behaviors.EventData) error {
		assert.Equal(data.Topic, behaviors.EventSequenceTopic)
		sequenceRaw, ok := data.Payload.Get(behaviors.EventSequenceEventsPayload)
		assert.True(ok)
		sequence, ok := sequenceRaw.(*behaviors.EventDatas)
		assert.True(ok)
		assert.Length(sequence, 3)
		return nil
	})
	assert.Nil(err)
}

// EOF
