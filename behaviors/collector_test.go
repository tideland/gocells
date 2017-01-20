// Tideland Go Cells - Behaviors - Unit Tests - Collector
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

// TestCollectorBehavior tests the collector behavior.
func TestCollectorBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	env := cells.NewEnvironment("collector-behavior")
	defer env.Stop()

	env.StartCell("collector", behaviors.NewCollectorBehavior(10))

	for i := 0; i < 25; i++ {
		env.EmitNew("collector", "collect", i)
	}

	time.Sleep(100 * time.Millisecond)

	collected, err := env.Request("collector", cells.CollectedTopic, nil, cells.DefaultTimeout)
	assert.Nil(err)
	assert.Length(collected, 10)

	err = env.EmitNew("collector", cells.ResetTopic, nil)
	assert.Nil(err)

	collected, err = env.Request("collector", cells.CollectedTopic, nil, cells.DefaultTimeout)
	assert.Nil(err)
	assert.Length(collected, 0)
}

// EOF
