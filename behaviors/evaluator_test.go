// Tideland Go Cells - Behaviors - Unit Tests - Evaluator
//
// Copyright (C) 2010-2017 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors_test

//--------------------
// IMPORTS
//--------------------

import (
	"strconv"
	"testing"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestEvaluatorBehavior tests the evaluator behavior. Scenario
// is to wait until the average evaluation has been larger than
// 6.0 for 3 times.
func TestEvaluatorBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sigc := audit.MakeSigChan()
	generator := audit.NewGenerator(audit.FixedRand())
	env := cells.NewEnvironment("evaluator-behavior")
	defer env.Stop()

	evaluator := func(event cells.Event) (float64, error) {
		i, err := strconv.Atoi(event.Topic())
		assert.Nil(err)
		return float64(i), nil
	}
	stopper := func(cell cells.Cell, event cells.Event) error {
		// TODO 2017-06-03 Mue Find criterion for stopping.
	}
	topics := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	env.StartCell("evaluator", behaviors.NewEvaluatorBehavior(evaluator))
	env.StartCell("stopper", behaviors.NewSimpleBehavior(stopper))
	env.Subscribe("evaluator", "stopper")

	go func() {
		for i := 0; i < 10000; i++ {
			topic := generator.OneStringOf(topics...)
			env.EmitNew(ctx, "evaluator", topic, nil)
		}
	}()
}

// EOF
