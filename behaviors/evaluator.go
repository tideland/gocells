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
// CONSTANTS
//--------------------

const (
	// TopicEvaluation labes an event as emitted evaluation.
	TopicEvaluation = "evaluation"

	// PayloadEvaluationAverage contains the average evaluated value.
	PayloadEvaluationAverage = "evaluation:average"

	// PayloadEvaluationCount contains the number of evaluated events.
	PayloadEvaluationCount = "evaluation:count"

	// PayloadEvaluationMax contains the maximum evaluated value.
	PayloadEvaluationMax = "evaluation:max"

	// PayloadEvaluationMin contains the minimum evaluated value.
	PayloadEvaluationMin = "evaluation:min"
)

//--------------------
// EVALUATOR BEHAVIOR
//--------------------

// Evaluator is a function returning a rating for each received event.
type Evaluator func(event cells.Event) (float64, error)

// evaluatorBehavior implements the evaluator behavior.
type evaluatorBehavior struct {
	cell      cells.Cell
	evaluate  Evaluator
	count     int
	minRating float64
	maxRating float64
	avgRating float64
}

// NewEvaluatorBehavior creates a behavior evaluating received events based
// on the passed function. This function returns a rating. Their minimum,
// maximum, average, and number of events are emitted. A "reset!" topic
// sets all values to zero again.
func NewEvaluatorBehavior(evaluator Evaluator) cells.Behavior {
	return &evaluatorBehavior{
		evaluate:  evaluator,
		count:     0,
		minRating: 0.0,
		maxRating: 0.0,
		avgRating: 0.0,
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

// ProcessEvent evaluates the event.
func (b *evaluatorBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case cells.TopicReset:
		b.count = 0
		b.minRating = 0.0
		b.maxRating = 0.0
		b.avgRating = 0.0
	default:
		rating, err := b.evaluate(event)
		if err != nil {
			return err
		}
		// Calculate values.
		if b.count == 0 {
			b.count = 1
			b.minRating = rating
			b.maxRating = rating
			b.avgRating = rating
		} else {
			b.avgRating = (b.avgRating*float64(b.count) + rating) / float64(b.count+1)
			b.count = b.count + 1
			if rating > b.maxRating {
				b.maxRating = rating
			}
			if rating < b.minRating {
				b.minRating = rating
			}
		}
		// Emit value.
		b.cell.EmitNew(event.Context(), TopicEvaluation, cells.PayloadValues{
			PayloadEvaluationCount:   b.count,
			PayloadEvaluationAverage: b.avgRating,
			PayloadEvaluationMax:     b.maxRating,
			PayloadEvaluationMin:     b.minRating,
		})
	}
	return nil
}

// Recover from an error.
func (b *evaluatorBehavior) Recover(err interface{}) error {
	b.count = 0
	b.minRating = 0.0
	b.maxRating = 0.0
	b.avgRating = 0.0
	return nil
}

// EOF
