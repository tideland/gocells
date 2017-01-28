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

// TestSequenceBehavior tests the event sequence behavior.
func TestSequenceBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	generator := audit.NewGenerator(audit.FixedRand())
	ctx := context.Background()
	env := cells.NewEnvironment("sequence-behavior")
	defer env.Stop()

	sequence := []string{"a", "b", "now"}
	matches := func(accessor cells.EventSinkAccessor) behaviors.CriterionMatch {
		matcher := func(index int, event cells.Event) (bool, error) {
			ok := event.Topic() == sequence[index]
			return ok, nil
		}
		matches, err := accessor.Match(matcher)
		if err != nil || !matches {
			return behaviors.CriterionClear
		}
		if accessor.Len() == len(sequence) {
			return behaviors.CriterionDone
		}
		return behaviors.CriterionKeep
	}
	topics := []string{"a", "b", "c", "d", "now"}

	env.StartCell("sequencer", behaviors.NewSequenceBehavior(matches))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10000))
	env.Subscribe("sequencer", "collector")

	for i := 0; i < 1000; i++ {
		topic := generator.OneStringOf(topics...)
		env.EmitNew(ctx, "sequencer", topic, nil)
		generator.SleepOneOf(0, 1*time.Millisecond, 2*time.Millisecond)
	}

	accessor, err := behaviors.RequestCollectedAccessor(env, "collector", cells.DefaultTimeout)
	assert.Nil(err)
	assert.NotEmpty(accessor)
	assert.Logf("Collected Sequences: %d", accessor.Len())
	err = accessor.Do(func(index int, event cells.Event) error {
		assert.Equal(event.Topic(), behaviors.EventSequenceTopic)
		csequenceRaw := event.Payload().Get(behaviors.EventSequenceEventsPayload, nil)
		csequence, ok := csequenceRaw.(cells.EventSink)
		assert.True(ok)
		assert.Length(csequence, 3)
		return csequence.Do(func(cindex int, cevent cells.Event) error {
			assert.Equal(cevent.Topic(), sequence[cindex])
			return nil
		})
	})
	assert.Nil(err)
}

// EOF
