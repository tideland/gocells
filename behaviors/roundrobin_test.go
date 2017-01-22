// Tideland Go Cells - Behaviors - Unit Tests - Round-Robin
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
	"time"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestRoundRobinBehavior tests the round robin behavior.
func TestRoundRobinBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	env := cells.NewEnvironment("round-robin-behavior")
	defer env.Stop()

	env.StartCell("round-robin", behaviors.NewRoundRobinBehavior())
	env.StartCell("round-robin-1", behaviors.NewCollectorBehavior(10))
	env.StartCell("round-robin-2", behaviors.NewCollectorBehavior(10))
	env.StartCell("round-robin-3", behaviors.NewCollectorBehavior(10))
	env.StartCell("round-robin-4", behaviors.NewCollectorBehavior(10))
	env.StartCell("round-robin-5", behaviors.NewCollectorBehavior(10))
	env.Subscribe("round-robin", "round-robin-1", "round-robin-2", "round-robin-3", "round-robin-4", "round-robin-5")

	time.Sleep(100 * time.Millisecond)

	// Just 23 to let two cells receive less events.
	for i := 0; i < 25; i++ {
		err := env.EmitNew("round-robin", "round", i)
		assert.Nil(err)
	}

	time.Sleep(100 * time.Millisecond)

	test := func(id string) int {
		collected, err := env.Request(id, cells.CollectedTopic, nil, cells.DefaultTimeout)
		assert.Nil(err)
		events, ok := collected.(cells.EventDatas)
		assert.True(ok)
		assert.Length(events, 5)
		return events.Len()
	}

	l1 := test("round-robin-1")
	l2 := test("round-robin-2")
	l3 := test("round-robin-3")
	l4 := test("round-robin-4")
	l5 := test("round-robin-5")

	assert.Equal(l1+l2+l3+l4+l5, 25)
}

// EOF
