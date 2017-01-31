// Tideland Go Cells - Behaviors - Unit Tests - Waiter
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

	"github.com/tideland/golib/audit"

	"time"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestWaiterBehavior tests the waiter behavior.
func TestWaiterBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	ctx := context.Background()
	env := cells.NewEnvironment("waiter-behavior")
	defer env.Stop()

	waiter := cells.NewPayloadWaiter()
	env.StartCell("waiter", behaviors.NewWaiterBehavior(waiter))

	env.EmitNew(ctx, "waiter", "foo", "1")
	env.EmitNew(ctx, "waiter", "bar", "2")
	env.EmitNew(ctx, "waiter", "baz", "3")

	waitCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	payload, err := waiter.Wait(waitCtx)
	assert.Nil(err)
	assert.Equal(payload.GetDefault("-"), "1")
}

// EOF
