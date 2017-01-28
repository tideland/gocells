// Tideland Go Cells - Behaviors - Unit Tests - Filter
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestFilterBehavior tests the filter behavior.
func TestFilterBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	ctx := context.Background()
	env := cells.NewEnvironment("filter-behavior")
	defer env.Stop()

	var wg sync.WaitGroup
	ff := func(id string, event cells.Event) bool {
		payload, ok := event.Payload().GetDefault(nil).(string)
		if !ok {
			return false
		}
		return event.Topic() == payload
	}
	sf := func(c cells.Cell, event cells.Event) error {
		wg.Done()
		return nil
	}
	env.StartCell("filter", behaviors.NewFilterBehavior(ff))
	env.StartCell("simple", behaviors.NewSimpleProcessorBehavior(sf))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10))
	env.Subscribe("filter", "simple", "collector")

	wg.Add(2)
	env.EmitNew(ctx, "filter", "a", "a")
	env.EmitNew(ctx, "filter", "a", "b")
	env.EmitNew(ctx, "filter", "b", "b")

	wg.Wait()
	accessor, err := behaviors.RequestCollectedAccessor(env, "collector", time.Second)
	assert.Nil(err)
	assert.Length(accessor, 2)
}

// EOF
