// Tideland Go Cells - Behaviors - Unit Tests - Key/Value
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

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestKeyValueBehavior tests the key/value behavior.
func TestKeyValueBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	generator := audit.NewGenerator(audit.FixedRand())
	env := cells.NewEnvironment("key-value-behavior")
	defer env.Stop()

	env.StartCell("keyvalue", behaviors.NewKeyValueBehavior(5))

	topics := []string{"alpha", "beta", "gamma"}
	payloads := []int{1, 2, 3, 4, 5}

	for i := 0; i < 50; i++ {
		topic := generator.OneStringOf(topics...)
		payload := generator.OneIntOf(payloads...)
		env.EmitNew("keyvalue", topic, payload)
	}

	// TODO 2017-07-13 Mue Work in progress.
	assert.True(true)
}

// EOF
