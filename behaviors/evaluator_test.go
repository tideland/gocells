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
	env := cells.NewEnvironment("evaluator-behavior")
	defer env.Stop()

	topics := []string{"1", "2", "1", "1", "5", "2", "3", "1", "4", "9"}
	evaluator := func(event cells.Event) (float64, error) {
		i, err := strconv.Atoi(event.Topic())
		assert.Nil(err)
		return float64(i), nil
	}
	processor := func(accessor cells.EventSinkAccessor) (cells.Payload, error) {
		event, ok := accessor.PeekLast()
		assert.True(ok)
		sigc <- event
		return nil, nil
	}

	env.StartCell("evaluator", behaviors.NewEvaluatorBehavior(evaluator))
	env.StartCell("collector", behaviors.NewCollectorBehavior(1000, processor))
	env.Subscribe("evaluator", "collector")

	for _, topic := range topics {
		env.EmitNew("evaluator", topic, nil)
	}
	time.Sleep(time.Second)

	env.EmitNew("collector", cells.TopicProcess, cells.PayloadClear)
	assert.WaitTested(sigc, func(value interface{}) error {
		event, ok := value.(cells.Event)
		assert.True(ok)
		var evaluation behaviors.Evaluation
		err := event.Payload().Unmarshal(&evaluation)
		assert.Equal(evaluation.Count, 10)
		assert.Equal(evaluation.MinRating, 1.0)
		assert.Equal(evaluation.MaxRating, 9.0)
		assert.Equal(evaluation.AvgRating, 2.9)
		return err
	}, time.Second)
}

// EOF
