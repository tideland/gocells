// Tideland Go Cells - Behaviors - Unit Tests - Filter
//
// Copyright (C) 2010-2016 Frank Mueller / Tideland / Oldenburg / Germany
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

// TestFilterBehavior tests the filter behavior.
func TestFilterBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	env := cells.NewEnvironment("filter-behavior")
	defer env.Stop()

	ff := func(id string, event cells.Event) bool {
		dp, ok := event.Payload().Get(cells.DefaultPayload)
		if !ok {
			return false
		}
		payload := dp.(string)
		return event.Topic() == payload
	}
	env.StartCell("filter", behaviors.NewFilterBehavior(ff))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10))
	env.Subscribe("filter", "collector")

	env.EmitNew("filter", "a", "a")
	env.EmitNew("filter", "a", "b")
	env.EmitNew("filter", "b", "b")

	time.Sleep(100 * time.Millisecond)

	collected, err := env.Request("collector", cells.CollectedTopic, nil, cells.DefaultTimeout, nil)
	assert.Nil(err)
	assert.Length(collected, 2, "two collected events")
}

// EOF
