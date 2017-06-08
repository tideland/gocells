// Tideland Go Cells - Behaviors - Unit Tests - Once
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

// TestOnceBehavior tests the once behavior.
func TestOnceBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("once-behavior")
	defer env.Stop()

	oneTimer := func(cell cells.Cell, event cell.Event) error {
		sigc <- event.Topic()
		err := cell.EmitNew(cell.ID(), event.Payload())
		return err
	}
	env.StartCell("first", behaviors.NewOnceBehavior(oneTimer))
	env.StartCell("second", behaviors.NewOnceBehavior(oneTimer))
	env.Subscribe("first", "second")

	env.EmitNew("first", "foo", "1")
	env.EmitNew("first", "bar", "2")

	audit.Wait(sigc, "foo", time.Second)
	audit.Wait(sigc, "first", time.Second)
}

// EOF
