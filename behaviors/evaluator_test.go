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
	"context"
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

// TestEvaluatorBehavior tests the evaluator behavior. Scenario
// is to wait until the average evaluation has been larger than
// 6.0 for 3 times.
func TestEvaluatorBehavior(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	generator := audit.NewGenerator(audit.FixedRand())
	ctx := context.Background()
	env := cells.NewEnvironment("evaluator-behavior")
	defer env.Stop()

	evaluate := func(event cells.Event) (float64, error) {
		i, err := strconv.Atoi(event.Topic())
		assert.Nil(err)
		return float64(i), nil
	}
	matches := func(accessor cells.EventSinkAccessor) cells.CriterionMatch {
		ok, err := accessor.Match(func(index int, event cells.Event) (bool, error) {
			avg, _ := event.Payload().GetFloat64(behaviors.PayloadEvaluationAverage)
			return avg > 6.0, nil
		})
		assert.Nil(err)
		if !ok {
			return cells.CriterionDropLast
		}
		if accessor.Len() < 3 {
			return cells.CriterionKeep
		}
		return cells.CriterionDone
	}
	topics := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	waiter := cells.NewPayloadWaiter()

	env.StartCell("evaluator", behaviors.NewEvaluatorBehavior(evaluate))
	env.StartCell("combo", behaviors.NewComboBehavior(matches))
	env.StartCell("waiter", behaviors.NewWaiterBehavior(waiter))
	env.Subscribe("evaluator", "combo")
	env.Subscribe("combo", "waiter")

	go func() {
		for i := 0; i < 10000; i++ {
			topic := generator.OneStringOf(topics...)
			env.EmitNew(ctx, "evaluator", topic, nil)
		}
	}()

	waitCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	payload, err := waiter.Wait(waitCtx)
	assert.Nil(err)
	accessor, ok := payload.Get(behaviors.PayloadComboEvents, nil).(cells.EventSinkAccessor)
	assert.True(ok)
	assert.Length(accessor, 3)
	accessor.Match(func(index int, event cells.Event) (bool, error) {
		avg := event.Payload().GetFloat64(behaviors.PayloadEvaluationAverage, 0.0)
		assert.True(avg > 6.0)
		assert.Logf("avg: %f", avg)
		return false, nil
	})
}

// EOF
