// Tideland Go Cells - Behaviors - Unit Tests - Status
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

// TestStatusBehavior tests the status behavior.
func TestStatusBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("status-behavior")
	defer env.Stop()

	callback := func(cell cells.Cell, event cells.Event) error {
		switch event.Topic() {
		case cells.TopicStatus:
			cellID := event.Payload().String()
			cell.Environment().EmitNew(cellID, "callback-status", "active")
		default:
			assert.Logf("doing something: %v", event.Topic())
		}
		return nil
	}
	statusProcessor := func(event cells.Event) error {
		sigc <- event.Topic()
		return nil
	}

	env.StartCell("callback", behaviors.NewCallbackBehavior(callback))
	env.StartCell("status", behaviors.NewStatusBehavior(statusProcessor))

	env.EmitNew("callback", "one", nil)
	env.EmitNew("callback", "two", nil)
	env.EmitNew("callback", cells.TopicStatus, "status")
	env.EmitNew("callback", "three", nil)

	assert.Wait(sigc, "callback-status", time.Second)
}

// EOF
