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
	sigc := audit.MakeSigChan()
	generator := audit.NewGenerator(audit.FixedRand())
	env := cells.NewEnvironment("sequence-behavior")
	defer env.Stop()

	sequence := []string{"a", "b", "now"}
	sequencer := func(accessor cells.EventSinkAccessor) cells.CriterionMatch {
		analyzer := cells.NewEventSinkAnalyzer(accessor)
		matcher := func(index int, event cells.Event) (bool, error) {
			ok := event.Topic() == sequence[index]
			return ok, nil
		}
		matches, err := analyzer.Match(matcher)
		if err != nil || !matches {
			return cells.CriterionClear
		}
		if accessor.Len() == len(sequence) {
			return cells.CriterionDone
		}
		return cells.CriterionKeep
	}
	analyzer := func(accessor cells.EventSinkAccessor) (cells.Payload, error) {
		first, ok := accessor.PeekFirst()
		assert.True(ok)
		return first.Payload(), nil
	}
	processor := func(accessor cells.EventSinkAccessor) (cells.Payload, error) {
		var indexes []int
		err := accessor.Do(func(_ int, event cells.Event) error {
			var index int
			event.Payload().Unmarshal(&index)
			indexes = append(indexes, index)
			return nil
		})
		assert.Nil(err)
		sigc <- indexes
		return nil, nil
	}
	topics := []string{"a", "b", "c", "d", "now"}

	env.StartCell("sequencer", behaviors.NewSequenceBehavior(sequencer, analyzer))
	env.StartCell("collector", behaviors.NewCollectorBehavior(100, processor))
	env.Subscribe("sequencer", "collector")

	for i := 0; i < 1000; i++ {
		topic := generator.OneStringOf(topics...)
		env.EmitNew("sequencer", topic, i)
		generator.SleepOneOf(0, 1*time.Millisecond, 2*time.Millisecond)
	}

	env.EmitNew("collector", cells.TopicProcess, nil)
	assert.Wait(sigc, []int{155, 269, 287, 298, 523, 888}, time.Minute)
}

// EOF
