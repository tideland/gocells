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

	"context"

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

	sink := cells.NewEventSink(10)
	env.StartCell("collector", behaviors.NewCollectorBehavior(sink))

	for i := 0; i < 25; i++ {
		env.EmitNew("collector", "collect", i)
	}

	ctx := context.WithTimeout(context.Background(), time.Second)
	waiter := cells.NewPayloadWaiter()
	env.EmitNew("collector", cells.CollectedTopic, waiter)
	payload, err := waiter.Wait(ctx)
	assert.Nil(err)
	accessor := payload.Default(nil).(cells.EventSinkAccessor)
	assert.NotNil(accessor)
	assert.Length(accessor, sink.Len())

	env.EmitNew("collector", cells.ResetTopic, nil)

	ctx = context.WithTimeout(context.Background(), time.Second)
	waiter = cells.NewPayloadWaiter()
	env.EmitNew("collector", cells.CollectedTopic, waiter)
	payload, err = waiter.Wait(ctx)
	assert.Nil(err)
	accessor = payload.Default(nil).(cells.EventSinkAccessor)
	assert.NotNil(accessor)
	assert.Empty(accessor)
}

// EOF
