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

	processor := func(index int, event cells.Event) error {
		if index == 2 {
			wg.Done()
		}
	}

	env.StartCell("broadcast", behaviors.NewBroadcasterBehavior())
	env.StartCell("test-a", behaviors.NewCollectorBehavior(10, processor))
	env.StartCell("test-b", behaviors.NewCollectorBehavior(10, processor))
	env.StartCell("signaller", behaviors.NewSimpleProcessorBehavior(func(cell cells.Cell, event cells.Event) error {
		wg.Done()
		return nil
	}))
	env.Subscribe("broadcast", "test-a", "test-b")
	env.Subscribe("test-b", "signaller")

	wg.Add(3)

	env.EmitNew("broadcast", "test", nil)
	env.EmitNew("broadcast", "test", nil)
	env.EmitNew("broadcast", "test", nil)

	wg.Wait()

	env.EmitNew("broadcast", cells.TopicCollected, nil)
}

// EOF
