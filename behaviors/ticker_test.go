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
	env := cells.NewEnvironment("ticker-behavior")
	defer env.Stop()

	env.StartCell("ticker", behaviors.NewTickerBehavior(50*time.Millisecond))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10))
	env.Subscribe("ticker", "collector")

	time.Sleep(125 * time.Millisecond)

	accessor, err := behaviors.RequestCollectedAccessor(env, "collector", cells.DefaultTimeout)
	assert.Nil(err)
	assert.Length(accessor, 2)
}

// EOF
