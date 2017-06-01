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
	generator := audit.NewGenerator(audit.FixedRand())
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("collector-behavior")
	defer env.Stop()

	processor := func(accessor cells.EventSinkAccessor) error {
		sigc <- accessor.Len()
		return nil
	}

	env.StartCell("collector", behaviors.NewCollectorBehavior(10, processor))

	for _, word := range generator.Words(25) {
		env.EmitNew("collector", "collect", cells.NewDefaultPayload(word))
	}

	env.EmitNew("collector", cells.TopicProcess, nil)

	audit.Wait(sigc, 10, time.Minute)

	env.EmitNew("collector", cells.TopicReset, nil)

	audit.Wait(sigc, 0, time.Minute)
}

// EOF
