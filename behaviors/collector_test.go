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

// TestCollectorBehavior tests the collector behavior.
func TestCollectorBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	ctx := context.Background()
	env := cells.NewEnvironment("collector-behavior")
	defer env.Stop()

	sink := cells.NewEventSink(10)
	env.StartCell("collector", behaviors.NewCollectorBehavior(sink))

	for i := 0; i < 25; i++ {
		env.EmitNew(ctx, "collector", "collect", i)
	}

	accessor, err := behaviors.RequestCollectedAccessor(ctx, env, "collector", time.Second)
	assert.Nil(err)
	assert.Length(accessor, sink.Len())

	env.EmitNew(ctx, "collector", cells.ResetTopic, nil)

	accessor, err = behaviors.RequestCollectedAccessor(ctx, env, "collector", time.Second)
	assert.Nil(err)
	assert.Empty(accessor)
}

// EOF
