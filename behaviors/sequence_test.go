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

	sequence := []string{"a", "b", "now"}
	matches := func(events cells.EventDatas) behaviors.CriterionMatch {
		matcher := func(index int, data *cells.EventData) (bool, error) {
			ok := data.Topic == sequence[index]
			return ok, nil
		}
		matches, err := events.Match(matcher)
		if err != nil || !matches {
			return behaviors.CriterionFailed
		}
		if events.Len() == len(sequence) {
			return behaviors.CriterionDone
		}
		return behaviors.CriterionPartly
	}
	topics := []string{"a", "b", "c", "d", "now"}

	env.StartCell("sequencer", behaviors.NewSequenceBehavior(matches))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10000))
	env.Subscribe("sequencer", "collector")

	for i := 0; i < 1000; i++ {
		topic := generator.OneStringOf(topics...)
		env.EmitNew("sequencer", topic, nil)
		generator.SleepOneOf(0, 1*time.Millisecond, 2*time.Millisecond)
	}

	collectedRaw, err := env.Request("collector", cells.CollectedTopic, nil, cells.DefaultTimeout)
	assert.Nil(err)
	collected, ok := collectedRaw.(cells.EventDatas)
	assert.True(ok)
	assert.True(collected.Len() > 0)
	assert.Logf("Collected Sequences: %d", collected.Len())
	err = collected.Do(func(index int, data *cells.EventData) error {
		assert.Equal(data.Topic, behaviors.EventSequenceTopic)
		csequenceRaw, ok := data.Payload.Get(behaviors.EventSequenceEventsPayload)
		assert.True(ok)
		csequence, ok := csequenceRaw.(cells.EventDatas)
		assert.True(ok)
		assert.Length(csequence, 3)
		return csequence.Do(func(cindex int, cdata *cells.EventData) error {
			assert.Equal(cdata.Topic, sequence[cindex])
			return nil
		})
	})
	assert.Nil(err)
}

// EOF
