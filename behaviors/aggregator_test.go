// Tideland Go Cells - Behaviors - Unit Tests - Aggregator
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

// TestAggregatorBehavior tests the aggregator behavior. Scenario
// is simply to count the lengths of the random topic until it
// reached the value 100.
func TestAggregatorBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	generator := audit.NewGenerator(audit.FixedRand())
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("aggregator-behavior")
	defer env.Stop()

	aggregate := func(payload cells.Payload, event cells.Event) (cells.Payload, error) {
		var topics []string
		err := payload.Unmarshal(&topics)
		if err != nil {
			return nil, err
		}
		topics = append(topics, event.Topic())
		return cells.NewPayload(topics), nil
	}
	match := func(event cells.Event) (bool, error) {
		var topics []string
		err := payload.Unmarshal(&topics)
		if err != nil {
			return false, err
		}
		return len(topics) > 19, nil
	}
	wait := func(event cells.Event) error {
		var topics []string
		err := payload.Unmarshal(&topics)
		if err != nil {
			return err
		}
		sigc <- len(topics)
		return nil
	}

	env.StartCell("aggregator", behaviors.NewAggregatorBehavior(aggregate))
	env.StartCell("filter", behaviors.NewFilterBehavior(match))
	env.StartCell("waiter", behaviors.NewWaiterBehavior(wait))
	env.Subscribe("aggregator", "filter")
	env.Subscribe("filter", "waiter")

	go func() {
		for i := 0; i < 50; i++ {
			topic := generator.Word()
			env.EmitNew("aggregator", topic, nil)
		}
	}()

	assert.Wait(sigc, 20, time.Minute)
}

// EOF
