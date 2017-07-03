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
	"fmt"
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
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("round-robin-behavior")
	defer env.Stop()

	processor := func(accessor cells.EventSinkAccessor) (cells.Payload, error) {
		sigc <- accessor.Len()
		return nil, nil
	}

	env.StartCell("round-robin", behaviors.NewRoundRobinBehavior())
	env.StartCell("round-robin-1", behaviors.NewCollectorBehavior(10, processor))
	env.StartCell("round-robin-2", behaviors.NewCollectorBehavior(10, processor))
	env.StartCell("round-robin-3", behaviors.NewCollectorBehavior(10, processor))
	env.StartCell("round-robin-4", behaviors.NewCollectorBehavior(10, processor))
	env.StartCell("round-robin-5", behaviors.NewCollectorBehavior(10, processor))
	env.Subscribe("round-robin", "round-robin-1", "round-robin-2", "round-robin-3", "round-robin-4", "round-robin-5")

	for i := 0; i < 25; i++ {
		err := env.EmitNew("round-robin", "round", i)
		assert.Nil(err)
	}

	for i := 1; i < 6; i++ {
		cellID := fmt.Sprintf("round-robin-%d", i)
		env.EmitNew(cellID, cells.TopicProcess, nil)
		assert.Wait(sigc, 5, time.Second)
	}
}

// EOF
