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
	cell          cells.Cell
	evaluate      Evaluator
	maxRatings    int
	ratings       []float64
	sortedRatings []float64
	evaluation    Evaluation
}

// NewEvaluatorBehavior creates a behavior evaluating received events based
// on the passed function. This function returns a rating. Their minimum,
// maximum, average, median, and number of events are emitted. The number
// of ratings for the median calculation is unlimited. So think about
// choosing NewMovingEvaluatorBehavior() to create the behavior with a
// limit and reduce memory usage.
//
// A "reset" topic sets all values to zero again.
func NewEvaluatorBehavior(evaluator Evaluator) cells.Behavior {
	return NewMovingEvaluatorBehavior(evaluator, 0)
}

// NewMovingEvaluatorBehavior creates the evaluator behavior with a
// moving rating window for calculation.
func NewMovingEvaluatorBehavior(evaluator Evaluator, limit int) cells.Behavior {
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
		b.sortedRatings = nil
		b.evaluation = Evaluation{}
	default:
		// Evaluate event and collect rating.
		rating, err := b.evaluate(event)
		if err != nil {
			return err
		}
		b.ratings = append(b.ratings, rating)
		if b.maxRatings > 0 && len(b.ratings) > b.maxRatings {
			b.ratings = b.ratings[1:]
		}
		if len(b.sortedRatings) < len(b.ratings) {
			// Let it grow up to the needed size.
			b.sortedRatings = append(b.sortedRatings, 0.0)
		}
		// Evaluate ratings.
		b.evaluateRatings()
		b.cell.EmitNew(TopicEvaluation, b.evaluation)
	}
	return nil
}

// Recover from an error.
func (b *evaluatorBehavior) Recover(err interface{}) error {
	b.ratings = nil
	b.sortedRatings = nil
	b.evaluation = Evaluation{}
	return nil
}

// evaluateRatings evaluates the collected ratings.
func (b *evaluatorBehavior) evaluateRatings() {
	copy(b.sortedRatings, b.ratings)
	sort.Float64s(b.sortedRatings)
	// Count.
	b.evaluation.Count = len(b.sortedRatings)
	// Average.
	totalRating := 0.0
	for _, rating := range b.sortedRatings {
		totalRating += rating
	}
	b.evaluation.AvgRating = totalRating / float64(b.evaluation.Count)
	// Median.
	if b.evaluation.Count%2 == 0 {
		// Even, have to calculate.
		middle := b.evaluation.Count / 2
		b.evaluation.MedRating = (b.sortedRatings[middle-1] + b.sortedRatings[middle]) / 2
	} else {
		// Odd, can take the middle.
		b.evaluation.MedRating = b.sortedRatings[b.evaluation.Count/2]
	}
	// Minimum and maximum.
	b.evaluation.MinRating = b.sortedRatings[0]
	b.evaluation.MaxRating = b.sortedRatings[len(b.sortedRatings)-1]
}

// EOF
