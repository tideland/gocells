// Tideland Go Cells - Behaviors - Unit Tests - Counter
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

// TestCounterBehavior tests the counting of events.
func TestCounterBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	env := cells.NewEnvironment("counter-behavior")
	defer env.Stop()

	cf := func(id string, event cells.Event) []string {
		return event.Payload().GetDefault([]string{}).([]string)
	}
	env.StartCell("counter", behaviors.NewCounterBehavior(cf))

	env.EmitNew("counter", "count", []string{"a", "b"})
	env.EmitNew("counter", "count", []string{"a", "c", "d"})
	env.EmitNew("counter", "count", []string{"a", "d"})

	counters, err := behaviors.RequestCounterResults(env, "counter")
	assert.Nil(err)
	assert.Length(counters, 4, "four counted events")

	assert.Equal(counters["a"], int64(3))
	assert.Equal(counters["b"], int64(1))
	assert.Equal(counters["c"], int64(1))
	assert.Equal(counters["d"], int64(2))

	err = env.EmitNew("counter", cells.ResetTopic, nil)
	assert.Nil(err)

	counters, err = behaviors.RequestCounterResults(env, "counter")
	assert.Nil(err)
	assert.Empty(counters)
}

// EOF
