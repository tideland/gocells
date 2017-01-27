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

// TestBroadcasterBehavior tests the broadcast behavior.
func TestBroadcasterBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	ctx := context.Background()
	env := cells.NewEnvironment("broadcaster-behavior")
	defer env.Stop()

	sinkA := cells.NewEventSink(10)
	sinkB := cells.NewEventSink(10)

	env.StartCell("broadcast", behaviors.NewBroadcasterBehavior())
	env.StartCell("test-a", behaviors.NewCollectorBehavior(sinkA))
	env.StartCell("test-b", behaviors.NewCollectorBehavior(sinkB))
	env.Subscribe("broadcast", "test-a", "test-b")

	env.EmitNew(ctx, "broadcast", "test", "a")
	env.EmitNew(ctx, "broadcast", "test", "b")
	env.EmitNew(ctx, "broadcast", "test", "c")

	accessor, err := behaviors.RequestCollectedAccessor(ctx, env, "test-a", time.Second)
	assert.Nil(err)
	assert.Length(accessor, 3)

	accessor, err = behaviors.RequestCollectedAccessor(ctx, env, "test-b", time.Second)
	assert.Nil(err)
	assert.Length(accessor, 3)
}

// EOF
