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
	"time"

	"github.com/tideland/golib/audit"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestEvaluatorBehavior tests the evaluator behavior.
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
	filter := func(event cells.Event) (bool, error) {
		var payload float64
		err := event.Payload().Unmarshal(&payload)
		assert.Nil(err)
		return payload > 6.0, nil
	}
	processor := func(accessor cells.EventSinkAccessor) error {
		sigc <- accessor.Len()
		return nil
	}
	topics := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

	env.StartCell("evaluator", behaviors.NewEvaluatorBehavior(evaluator))
	env.StartCell("filter", behaviors.NewFilterBehavior(filter))
	env.StartCell("collector", behaviors.NewCollectorBehavior(10000, processor))
	env.Subscribe("evaluator", "filter")
	env.Subscribe("filter", "processor")

	for i := 0; i < 10000; i++ {
		topic := generator.OneStringOf(topics...)
		env.EmitNew("evaluator", topic, nil)
	}

	env.EmitNew("collector", cells.TopicProcess, nil)
	assert.Wait(sigc, 10, time.Minute)
}

// EOF
