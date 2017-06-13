// Tideland Go Cells - Behaviors - Unit Tests - Broadcaster
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

// TestBroadcasterBehavior tests the broadcast behavior.
func TestBroadcasterBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("broadcaster-behavior")
	defer env.Stop()

	processor := func(accessor cells.EventSinkAccessor) error {
		sigc <- accessor.Len()
		return nil
	}

	env.StartCell("broadcast", behaviors.NewBroadcasterBehavior())
	env.StartCell("test-a", behaviors.NewCollectorBehavior(10, processor))
	env.StartCell("test-b", behaviors.NewCollectorBehavior(10, processor))
	env.Subscribe("broadcast", "test-a", "test-b")

	env.EmitNew("broadcast", "test", nil)
	env.EmitNew("broadcast", "test", nil)
	env.EmitNew("broadcast", "test", nil)

	env.EmitNew("test-a", cells.TopicProcess, nil)
	assert.Wait(sigc, 3, time.Second)
	env.EmitNew("test-b", cells.TopicProcess, nil)
	assert.Wait(sigc, 3, time.Second)
}

// EOF
