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

// TestAggregatorBehavior tests the aggregator behavior. Scenario
// is simply to count the lengths of the random topic until it
// reached the value 100.
func TestAggregatorBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	generator := audit.NewGenerator(audit.FixedRand())
	ctx := context.Background()
	env := cells.NewEnvironment("aggregator-behavior")
	defer env.Stop()

	aggregate := func(value interface{}, event cells.Event) (interface{}, error) {
		current, ok := value.(int)
		if !ok {
			current = 0
		}
		current += len(event.Topic())
		return current, nil
	}
	matches := func(event cells.Event) (bool, error) {
		length := event.Payload().GetInt(behaviors.PayloadAggregatorValue)
		return length > 100, nil
	}
	waiter := cells.NewPayloadWaiter()

	env.StartCell("aggregator", behaviors.NewAggregatorBehavior(aggregate))
	env.StartCell("filter", behaviors.NewFilterBehavior(matches))
	env.StartCell("waiter", behaviors.NewWaiterBehavior(waiter))
	env.Subscribe("aggregator", "filter")
	env.Subscribe("filter", "waiter")

	go func() {
		for i := 0; i < 199; i++ {
			topic := generator.Word()
			env.EmitNew(ctx, "aggregator", topic, nil)
		}
	}()

	waitCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	payload, err := waiter.Wait(waitCtx)
	assert.Nil(err)
	length := payload.GetInt(behaviors.PayloadAggregatorValue, 0)
	assert.True(length > 100)
}

// EOF
