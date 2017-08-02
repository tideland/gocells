// Tideland Go Cells - Example - Totally Working Behaviors
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package main

//--------------------
// IMPORTS
//--------------------

import (
	"time"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// BEHAVIORS
//--------------------

// MakeAveragePercentChange1hEvaluator returns a behavior evaluating the
// average of all percent changes of one hour.
func MakeAveragePercentChange1hEvaluator() cells.Behavior {
	evaluator := func(event cells.Event) (float64, error) {
		var coins Coins
		if err := event.Payload().Unmarshal(&coins); err != nil {
			return 0.0, err
		}
		total := 0.0
		for _, coin := range coins {
			total += coin.PercentChange1h
		}
		return total / float64(len(coins)), nil
	}
	return behaviors.NewLimitedEvaluatorBehavior(evaluator, 360*24*7)
}

// MakeAveragePercentChange1hRater returns a behavior checking the
// rate of the average percent changes of a minute.
func MakeAveragePercentChange1hRater(up bool) cells.Behavior {
	var last float64
	criterion := func(event cells.Event) (bool, error) {
		var evaluation behaviors.Evaluation
		var matches bool
		if err := event.Payload().Unmarshal(&evaluation); err != nil {
			return false, err
		}
		switch {
		case up && evaluation.AvgRating > last:
			matches = true
		case !up && evaluation.AvgRating < last:
			matches = true
		}
		last = evaluation.AvgRating
		return matches, nil
	}
	// Time buffer a bit longer than a minute.
	return behaviors.NewRateWindowBehavior(criterion, 6, 65*time.Second)
}

// EOF
