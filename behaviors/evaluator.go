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
	"sort"

	"github.com/tideland/gocells/cells"
)

//--------------------
// CONSTANTS
//--------------------

const (
	// TopicEvaluation labes an event as emitted evaluation.
	TopicEvaluation = "evaluation"
)

//--------------------
// EVALUATOR BEHAVIOR
//--------------------

// Evaluator is a function returning a rating for each received event.
type Evaluator func(event cells.Event) (float64, error)

// Evaluation contains the aggregated result of all evaluations.
type Evaluation struct {
	Count     int
	MinRating float64
	MaxRating float64
	AvgRating float64
	MedRating float64
}

// evaluatorBehavior implements the evaluator behavior.
type evaluatorBehavior struct {
	cell       cells.Cell
	evaluate   Evaluator
	maxRatings int
	ratings    []float64
	evaluation Evaluation
}

// NewEvaluatorBehavior creates a behavior evaluating received events based
// on the passed function. This function returns a rating. Their minimum,
// maximum, average, median, and number of events are emitted. The number
// of ratings for the median calculation is unlimited. Choose
// NewLimitedEvaluatorBehavior() to create the behavior with a limit.
//
// A "reset" topic sets all values to zero again.
func NewEvaluatorBehavior(evaluator Evaluator) cells.Behavior {
	return NewLimitedEvaluatorBehavior(evaluator, 0)
}

// NewLimitedEvaluatorBehavior creates the evaluator behavior with a
// limit for median calculation.
func NewLimitedEvaluatorBehavior(evaluator Evaluator, limit int) cells.Behavior {
	return &evaluatorBehavior{
		evaluate:   evaluator,
		maxRatings: limit,
		evaluation: Evaluation{},
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
		b.ratings = nil
		b.evaluation = Evaluation{}
	default:
		rating, err := b.evaluate(event)
		if err != nil {
			return err
		}
		b.ratings = append(b.ratings, rating)
		if b.maxRatings > 0 && len(b.ratings) > b.maxRatings {
			b.ratings = b.ratings[1:]
		}
		numOfRatings := len(b.ratings)
		sort.Float64s(b.ratings)
		// Calculate values.
		if b.evaluation.Count == 0 {
			b.evaluation.Count = 1
			b.evaluation.MinRating = rating
			b.evaluation.MaxRating = rating
			b.evaluation.AvgRating = rating
			b.evaluation.MedRating = rating
		} else {
			totalRating := b.evaluation.AvgRating*float64(b.evaluation.Count) + rating
			b.evaluation.Count = b.evaluation.Count + 1
			b.evaluation.AvgRating = totalRating / float64(b.evaluation.Count)
			if numOfRatings%2 == 0 {
				// Even, have to calculate.
				middle := numOfRatings / 2
				b.evaluation.MedRating = (b.ratings[middle-1] + b.ratings[middle]) / 2
			} else {
				// Odd, can take the middle.
				b.evaluation.MedRating = b.ratings[numOfRatings/2]
			}
			if rating > b.evaluation.MaxRating {
				b.evaluation.MaxRating = rating
			}
			if rating < b.evaluation.MinRating {
				b.evaluation.MinRating = rating
			}
		}
		// Emit value.
		b.cell.EmitNew(TopicEvaluation, b.evaluation)
	}
	return nil
}

// Recover from an error.
func (b *evaluatorBehavior) Recover(err interface{}) error {
	b.ratings = nil
	b.evaluation = Evaluation{}
	return nil
}

// EOF
