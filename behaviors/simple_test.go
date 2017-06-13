// Tideland Go Cells - Behaviors - Unit Tests - Simple
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
	"sync"
	"testing"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestSimpleBehavior tests the simple processor behavior.
func TestSimpleBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	env := cells.NewEnvironment("simple-procesor-behavior")
	defer env.Stop()

	topics := []string{}
	var wg sync.WaitGroup
	spf := func(c cells.Cell, event cells.Event) error {
		topics = append(topics, event.Topic())
		wg.Done()
		return nil
	}
	env.StartCell("simple", behaviors.NewSimpleProcessorBehavior(spf))

	wg.Add(3)
	env.EmitNew("simple", "foo", "")
	env.EmitNew("simple", "bar", "")
	env.EmitNew("simple", "baz", "")

	wg.Wait()
	assert.Length(topics, 3)
}

// EOF
