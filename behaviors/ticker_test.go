// Tideland Go Cells - Behaviors - Unit Tests - Ticker
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

// TestTickerBehavior tests the ticker behavior.
func TestTickerBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("ticker-behavior")
	defer env.Stop()

	processor := func(accessor cells.EventSinkAccessor) error {
		sigc <- accessor.Len()
		return nil
	}

	env.StartCell("ticker", behaviors.NewTickerBehavior(50*time.Millisecond))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10, processor))
	env.Subscribe("ticker", "collector")

	time.Sleep(125 * time.Millisecond)

	env.EmitNew("collector", cells.TopicProcess, nil)
	assert.Wait(sigc, 2, time.Minute)
}

// EOF
