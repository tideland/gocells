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
	matcher := func(accessor cells.EventSinkAccessor) cells.CriterionMatch {
		amatcher := func(index int, event cells.Event) (bool, error) {
			ok := event.Topic() == sequence[index]
			return ok, nil
		}
		matches, err := accessor.Match(amatcher)
		if err != nil || !matches {
			return cells.CriterionClear
		}
		if accessor.Len() == len(sequence) {
			return cells.CriterionDone
		}
		return cells.CriterionKeep
	}
	analyzer := func(accessor cells.EventSinkAccessor) (cells.Payload, error) {
		// TODO 2017-06-12 Mue Added analyzing to return better payload.
		return nil, nil
	}
	processor := func(accessor cells.EventSinkAccessor) error {
		// TODO 2017-06-12 Mue Signal aggregated collected payloads.
		sigc <- accessor.Len()
		return nil
	}
	topics := []string{"a", "b", "c", "d", "now"}

	env.StartCell("sequencer", behaviors.NewSequenceBehavior(matcher, analyzer))
	env.StartCell("collector", behaviors.NewCollectorBehavior(100, processor))
	env.Subscribe("sequencer", "collector")

	for i := 0; i < 1000; i++ {
		topic := generator.OneStringOf(topics...)
		env.EmitNew("sequencer", topic, nil)
		generator.SleepOneOf(0, 1*time.Millisecond, 2*time.Millisecond)
	}

	env.EmitNew("collector", cells.TopicProcess, nil)
	assert.Wait(sigc, 10, time.Minute)
}

// EOF
