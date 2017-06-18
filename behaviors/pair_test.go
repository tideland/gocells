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

// TestPairBehavior tests the event pair behavior.
func TestPairBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sigc := audit.MakeSigChan()
	generator := audit.NewGenerator(audit.FixedRand())
	env := cells.NewEnvironment("pair-behavior")
	defer env.Stop()

	matches := func(event cells.Event, data cells.Payload) (cells.Payload, bool) {
		if event.Topic() == "now" {
			now := time.Now().Unix()
			payload, _ := cells.NewPayload(now)
			return payload, true
		}
		return nil, false
	}
	mkfilter := func(positive bool) behaviors.Filter {
		topic := behaviors.TopicPairTimeout
		if positive {
			topic = behaviors.TopicPair
		}
		return func(event cells.Event) (bool, error) {
			return event.Topic() == topic, nil
		}
	}
	mkcondition := func() (behaviors.ConditionTester, behaviors.ConditionProcessor) {
		counter := 0
		return func(event cells.Event) bool {
				counter++
				return counter == 25
			}, func(cell cells.Cell, event cells.Event) error {
				sigc <- counter
				return nil
			}
	}
	topics := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "now"}
	duration := time.Millisecond

	env.StartCell("pairer", behaviors.NewPairBehavior(matches, duration))
	env.StartCell("positive-filter", behaviors.NewFilterBehavior(mkfilter(true)))
	env.StartCell("negative-filter", behaviors.NewFilterBehavior(mkfilter(false)))
	env.StartCell("positive-condition", behaviors.NewConditionBehavior(mkcondition()))
	env.StartCell("negative-condition", behaviors.NewConditionBehavior(mkcondition()))
	env.Subscribe("pairer", "positive-filter", "negative-filter")
	env.Subscribe("positive-filter", "positive-condition")
	env.Subscribe("negative-filter", "negative-condition")

	for i := 0; i < 5000; i++ {
		topic := generator.OneStringOf(topics...)
		env.EmitNew("pairer", topic, nil)
		generator.SleepOneOf(0, time.Millisecond, 2*time.Millisecond)
	}

	assert.Wait(sigc, 25, 5*time.Second)
	assert.Wait(sigc, 25, 5*time.Second)
}

// EOF
