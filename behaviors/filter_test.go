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

	ff := func(id string, event cells.Event) bool {
		payload, ok := event.Payload().GetDefault(nil).(string)
		if !ok {
			return false
		}
		return event.Topic() == payload
	}
	sink := cells.NewEventSink(10)
	env.StartCell("filter", behaviors.NewFilterBehavior(ff))
	env.StartCell("collector", behaviors.NewCollectorBehavior(sink))
	env.Subscribe("filter", "collector")

	env.EmitNew("filter", "a", "a")
	env.EmitNew("filter", "a", "b")
	env.EmitNew("filter", "b", "b")

	accessor, err := behaviors.RequestCollectedAccessor(env, "collector")
	assert.Nil(err)
	assert.Length(accessor, 2)
}

// EOF
