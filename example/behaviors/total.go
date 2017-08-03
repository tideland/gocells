// Tideland Go Cells - Example - Behaviors - Working on the Total
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
	"time"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// BEHAVIORS
//--------------------

// MakeAvgPercentChange1h returns a behavior calculating the
// average PercentChange1h of all coins.
func MakeAvgPercentChange1h() cells.Behavior {
	mapper := func(event cells.Event) (cells.Event, error) {
		var coins Coins
		if err := event.Payload().Unmarshal(&coins); err != nil {
			return nil, err
		}
		total := 0.0
		for _, coin := range coins {
			total += coin.PercentChange1h
		}
		avg := total / float64(len(coins))
		return cells.NewEvent("avg-percent-change-1h", avg)
	}
	return behaviors.NewMapperBehavior(mapper)
}

// MakeTopPercentChange1hCoins returns a behavior receiving the
// average PercentChange1h of all coins as well as all coins.
// It emits a list of coins with higher than average changes.
func MakeTopPercentChange1hCoins() cells.Behavior {
	var avgPercentChange1h float64
	processor := func(cell cells.Cell, event cells.Event) error {
		switch event.Topic() {
		case "avg-percent-change-1h":
			if err := event.Payload().Unmarshal(&avgPercentChange1h); err != nil {
				return err
			}
		default:
			var coins Coins
			var topCoins Coins
			if err := event.Payload().Unmarshal(&coins); err != nil {
				return err
			}
			for _, coin := range coins {
				if coin.PercentChange1h > avgPercentChange1h {
					topCoins = append(topCoins, coin)
				}
			}
			cell.EmitNew("top-coins", topCoins)
		}
		return nil
	}
	return behaviors.NewSimpleProcessorBehavior(processor)
}

// MakeAvgMarketCapEvaluator returns a behavior evaluating the
// average of all market capitalizations of one hour.
func MakeAvgMarketCapEvaluator() cells.Behavior {
	evaluator := func(event cells.Event) (float64, error) {
		var coins Coins
		if err := event.Payload().Unmarshal(&coins); err != nil {
			return 0.0, err
		}
		total := 0.0
		for _, coin := range coins {
			total += coin.MarketCapUSD
		}
		return total / float64(len(coins)), nil
	}
	return behaviors.NewLimitedEvaluatorBehavior(evaluator, 360*24*7)
}

// MakeAvgMarketCapRater returns a behavior checking the change
// rate of the average market capitalization of a minute.
func MakeAvgMarketCapRater(up bool) cells.Behavior {
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
