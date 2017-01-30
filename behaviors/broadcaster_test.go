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

// TestBroadcasterBehavior tests the broadcast behavior.
func TestBroadcasterBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	ctx := context.Background()
	env := cells.NewEnvironment("broadcaster-behavior")
	defer env.Stop()

	var wg sync.WaitGroup

	env.StartCell("broadcast", behaviors.NewBroadcasterBehavior())
	env.StartCell("test-a", behaviors.NewCollectorBehavior(10))
	env.StartCell("test-b", behaviors.NewCollectorBehavior(10))
	env.StartCell("signaller", behaviors.NewSimpleProcessorBehavior(func(cell cells.Cell, event cells.Event) error {
		wg.Done()
		return nil
	}))
	env.Subscribe("broadcast", "test-a", "test-b")
	env.Subscribe("test-b", "signaller")

	wg.Add(3)

	env.EmitNew(ctx, "broadcast", "test", "a")
	env.EmitNew(ctx, "broadcast", "test", "b")
	env.EmitNew(ctx, "broadcast", "test", "c")

	wg.Wait()

	accessor, err := behaviors.RequestCollectedAccessor(env, "test-a", time.Second)
	assert.Nil(err)
	assert.Length(accessor, 3)

	accessor, err = behaviors.RequestCollectedAccessor(env, "test-b", time.Second)
	assert.Nil(err)
	assert.Length(accessor, 3)
}

// EOF
