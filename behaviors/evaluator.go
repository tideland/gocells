// Tideland Go Cells - Behaviors - Evaluator
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/gocells/cells"
)

//--------------------
// EVALUATOR BEHAVIOR
//--------------------

// Evaluator is a function returning a rate for each received event.
type Evaluator func(event cells.Event) (float64, error)

// evaluatorBehavior implements the evaluator behavior.
type evaluatorBehavior struct {
	cell     cells.Cell
	evaluate Evaluator
	count    int
	minRate  float64
	maxRate  float64
	avgRate  float64
}

// NewEvaluatorBehavior creates a behavior rating received events based
// on the passed function. This function returns a rate. Their minimum,
// maximum, average, and number of events are emitted. A "reset!" topic
// sets all values to zero again.
func NewEvaluatorBehavior(evaluator Evaluator) cells.Behavior {
	return &evaluatorBehavior{
		evaluate: evaluator,
		count:    0,
		minRate:  0.0,
		maxRate:  0.0,
		avgRate:  0.0,
	}
}

// Init the behavior.
func (b *evaluatorBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

// Terminate the behavior.
func (b *evaluatorBehavior) Terminate() error {
	return nil
}

// ProcessEvent calls the simple processor function.
func (b *evaluatorBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case cells.TopicReset:
		b.count = 0
		b.minRate = 0.0
		b.maxRate = 0.0
		b.avgRate = 0.0
	default:
		rate, err := b.evaluate(event)
		if err != nil {
			return err
		}
		// Calculate values.
		if b.count == 0 {
			b.count = 1
			b.minRate = rate
			b.maxRate = rate
			b.avgRate = rate
		} else {
			b.avgRate = (b.avgRate*float64(b.count) + rate) / float64(b.count+1)
			b.count = b.count + 1
			if rate > b.maxRate {
				b.maxRate = rate
			}
			if rate < b.minRate {
				b.minRate = rate
			}
		}
		// Emit value.
		b.cell.EmitNew(event.Context(), TopicEvaluation, cells.PayloadValues{
			PayloadEvaluationCount: b.count,
			PayloadEvaluationAvg:   b.avgRate,
			PayloadEvaluationMax:   b.maxRate,
			PayloadEvaluationMin:   b.minRate,
		})
	}
	return nil
}

// Recover from an error.
func (b *evaluatorBehavior) Recover(err interface{}) error {
	b.count = 0
	b.minRate = 0.0
	b.maxRate = 0.0
	b.avgRate = 0.0
	return nil
}

// EOF
