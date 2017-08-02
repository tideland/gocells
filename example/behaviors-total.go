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
	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// BEHAVIORS
//--------------------

// MakeAveragePercentChange1hEvaluator returns a behavior evaluating the
// average of all  percent changes of one hour.
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

// EOF
