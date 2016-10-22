// Tideland Go Cells - Behaviors - Unit Tests - Router
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
	"strings"
	"testing"
	"time"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestRouterBehavior tests the router behavior.
func TestRouterBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	env := cells.NewEnvironment("router-behavior")
	defer env.Stop()

	rf := func(emitterID, subscriberID string, event cells.Event) (bool, error) {
		ok := strings.Contains(event.Topic(), subscriberID)
		return ok, nil
	}
	env.StartCell("router", behaviors.NewRouterBehavior(rf))
	env.StartCell("test-1", behaviors.NewCollectorBehavior(10))
	env.StartCell("test-2", behaviors.NewCollectorBehavior(10))
	env.StartCell("test-3", behaviors.NewCollectorBehavior(10))
	env.StartCell("test-4", behaviors.NewCollectorBehavior(10))
	env.StartCell("test-5", behaviors.NewCollectorBehavior(10))
	env.Subscribe("router", "test-1", "test-2", "test-3", "test-4", "test-5")

	env.EmitNew("router", "test-1:test-2", "a")
	env.EmitNew("router", "test-1:test-2:test-3", "b")
	env.EmitNew("router", "test-3:test-4:test-5", "c")

	time.Sleep(100 * time.Millisecond)

	test := func(id string, length int) {
		collected, err := env.Request(id, cells.CollectedTopic, nil, cells.DefaultTimeout, nil)
		assert.Nil(err)
		assert.Length(collected, length)
	}

	test("test-1", 2)
	test("test-2", 2)
	test("test-3", 2)
	test("test-4", 1)
	test("test-5", 1)
}

// EOF
