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
	"sync"
	"testing"

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
	env := cells.NewEnvironment("filter-behavior")
	defer env.Stop()

	var wg sync.WaitGroup
	filter := func(event cells.Event) (bool, error) {
		var payload string
		err := event.Payload().Unmarshal(&payload)
		assert.Nil(err)
		return event.Topic() == payload, nil
	}
	processor := func(c cells.Cell, event cells.Event) error {
		var payload string
		err := event.Payload().Unmarshal(&payload)
		assert.Nil(err)
		assert.Equal(event.Topic(), payload)
		wg.Done()
		return nil
	}
	env.StartCell("filter", behaviors.NewFilterBehavior(filter))
	env.StartCell("simple", behaviors.NewSimpleProcessorBehavior(processor))
	env.Subscribe("filter", "simple")

	wg.Add(2)
	env.EmitNew("filter", "a", "a")
	env.EmitNew("filter", "a", "b")
	env.EmitNew("filter", "a", "c")
	env.EmitNew("filter", "a", "d")
	env.EmitNew("filter", "b", "b")

	wg.Wait()
}

// EOF
