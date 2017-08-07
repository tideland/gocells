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

	evaluator := func(event cells.Event) (float64, error) {
		i, err := strconv.Atoi(event.Topic())
		return float64(i), err
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

	// Standard evaluating.
	topics := []string{"1", "2", "1", "1", "3", "2", "3", "1", "3", "9"}
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
		assert.Equal(evaluation.AvgRating, 2.6)
		assert.Equal(evaluation.MedRating, 2.0)
		return err
	}, time.Second)

	// Reset and check with only one value.
	env.EmitNew("evaluator", cells.TopicReset, nil)
	env.EmitNew("evaluator", "4711", nil)
	time.Sleep(time.Second)

	env.EmitNew("collector", cells.TopicProcess, cells.PayloadClear)
	assert.WaitTested(sigc, func(value interface{}) error {
		event, ok := value.(cells.Event)
		assert.True(ok)
		var evaluation behaviors.Evaluation
		err := event.Payload().Unmarshal(&evaluation)
		assert.Equal(evaluation.Count, 1)
		assert.Equal(evaluation.MinRating, 4711.0)
		assert.Equal(evaluation.MaxRating, 4711.0)
		assert.Equal(evaluation.AvgRating, 4711.0)
		assert.Equal(evaluation.MedRating, 4711.0)
		return err
	}, time.Second)

	// Crash evaluating.
	topics = []string{"1", "2", "3", "4", "crash", "1", "2", "1", "2", "1"}
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
		assert.Equal(evaluation.Count, 5)
		assert.Equal(evaluation.MinRating, 1.0)
		assert.Equal(evaluation.MaxRating, 2.0)
		assert.Equal(evaluation.AvgRating, 1.4)
		assert.Equal(evaluation.MedRating, 1.0)
		return err
	}, time.Second)
}

// TestLimitedEvaluatorBehavior tests the limited evaluator behavior.
func TestLimitedEvaluatorBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	sigc := audit.MakeSigChan()
	env := cells.NewEnvironment("limited-evaluator-behavior")
	defer env.Stop()

	evaluator := func(event cells.Event) (float64, error) {
		i, err := strconv.Atoi(event.Topic())
		return float64(i), err
	}
	processor := func(accessor cells.EventSinkAccessor) (cells.Payload, error) {
		event, ok := accessor.PeekLast()
		assert.True(ok)
		sigc <- event
		return nil, nil
	}

	env.StartCell("evaluator", behaviors.NewMovingEvaluatorBehavior(evaluator, 5))
	env.StartCell("collector", behaviors.NewCollectorBehavior(1000, processor))
	env.Subscribe("evaluator", "collector")

	// Standard evaluating.
	topics := []string{"1", "2", "1", "1", "9", "2", "3", "1", "3", "2"}
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
		assert.Equal(evaluation.Count, 5)
		assert.Equal(evaluation.MinRating, 1.0)
		assert.Equal(evaluation.MaxRating, 3.0)
		assert.Equal(evaluation.AvgRating, 2.2)
		assert.Equal(evaluation.MedRating, 2.0)
		return err
	}, time.Second)
}

// EOF
