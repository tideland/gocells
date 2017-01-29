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

// TestComboBehavior tests the event combo behavior.
func TestComboBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	generator := audit.NewGenerator(audit.FixedRand())
	ctx := context.Background()
	env := cells.NewEnvironment("combo-behavior")
	defer env.Stop()

	matches := func(accessor cells.EventSinkAccessor) behaviors.CriterionMatch {
		combo := map[string]int{
			"a": 0,
			"b": 0,
			"c": 0,
			"d": 0,
		}
		matcher := func(index int, event cells.Event) (bool, error) {
			_, ok := combo[event.Topic()]
			if ok {
				combo[event.Topic()]++
			}
			return ok, nil
		}
		matches, err := accessor.Match(matcher)
		if err != nil || !matches {
			return behaviors.CriterionDropLast
		}
		for _, count := range combo {
			if count == 0 {
				return behaviors.CriterionKeep
			}
		}
		return behaviors.CriterionDone
	}
	mapper := func(id string, event cells.Event) (cells.Event, error) {
		sink, ok := event.Payload().Get(behaviors.EventComboEventsPayload, nil).(cells.EventSink)
		if !ok {
			assert.Fail("illegal payload")
		}
		topics := []string{}
		sink.Do(func(index int, event cells.Event) error {
			topics = append(topics, event.Topic())
			return nil
		})
		return cells.NewEvent(ctx, "topics", topics)
	}
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}

	env.StartCell("combiner", behaviors.NewComboBehavior(matches))
	env.StartCell("mapper", behaviors.NewMapperBehavior(mapper))
	env.StartCell("collector", behaviors.NewCollectorBehavior(100))
	env.Subscribe("combiner", "mapper")
	env.Subscribe("mapper", "collector")

	for i := 0; i < 1000; i++ {
		topic := generator.OneStringOf(topics...)
		env.EmitNew(ctx, "combiner", topic, nil)
		generator.SleepOneOf(0, 1*time.Millisecond, 2*time.Millisecond)
	}

	accessor, err := behaviors.RequestCollectedAccessor(env, "collector", cells.DefaultTimeout)
	assert.Nil(err)
	assert.NotEmpty(accessor)
	assert.Logf("Collected Combinations: %d", accessor.Len())
	accessor.Do(func(index int, event cells.Event) error {
		topics := event.Payload().GetDefault([]string{}).([]string)
		assert.True(len(topics) >= 4)
		assert.Contents("a", topics)
		assert.Contents("b", topics)
		assert.Contents("c", topics)
		assert.Contents("d", topics)
		return nil
	})
}

// EOF
