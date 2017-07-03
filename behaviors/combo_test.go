// Tideland Go Cells - Behaviors - Unit Tests - Event Combination
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

// TestComboBehavior tests the event combo behavior.
func TestComboBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	generator := audit.NewGenerator(audit.FixedRand())
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("combo-behavior")
	defer env.Stop()

	matcher := func(accessor cells.EventSinkAccessor) (cells.CriterionMatch, cells.Payload) {
		combo := map[string]int{
			"a": 0,
			"b": 0,
			"c": 0,
			"d": 0,
		}
		matches, err := accessor.Match(func(index int, event cells.Event) (bool, error) {
			_, ok := combo[event.Topic()]
			if ok {
				combo[event.Topic()]++
			}
			return ok, nil
		})
		if err != nil || !matches {
			return cells.CriterionDropLast, nil
		}
		for _, count := range combo {
			if count == 0 {
				return cells.CriterionKeep, nil
			}
		}
		payload, err := cells.NewPayload(combo)
		assert.Nil(err)
		return cells.CriterionDone, payload
	}
	processor := func(accessor cells.EventSinkAccessor) (cells.Payload, error) {
		ok, err := accessor.Match(func(index int, event cells.Event) (bool, error) {
			var payload map[string]int
			if err := event.Payload().Unmarshal(&payload); err != nil {
				return false, err
			}
			if len(payload) != 4 {
				return false, nil
			}
			for key := range payload {
				if payload[key] == 0 {
					return false, nil
				}
			}
			return true, nil
		})
		sigc <- ok
		return nil, err
	}

	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}

	env.StartCell("combiner", behaviors.NewComboBehavior(matcher))
	env.StartCell("collector", behaviors.NewCollectorBehavior(100, processor))
	env.Subscribe("combiner", "collector")

	for i := 0; i < 1000; i++ {
		topic := generator.OneStringOf(topics...)
		env.EmitNew("combiner", topic, nil)
	}

	env.EmitNew("collector", cells.TopicProcess, nil)
	assert.Wait(sigc, true, time.Minute)
}

// EOF
