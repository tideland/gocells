// Tideland Go Cells - Example - Behaviors - Global usable ones
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
	"strconv"
	"time"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// BEHAVIORS
//--------------------

// MakeLogger returns a logger behavior.
func MakeLogger() cells.Behavior {
	return behaviors.NewLoggerBehavior()
}

// MakeRawCoinsConverter returns a behavior converting raw coins
// into correct typed ones.
func MakeRawCoinsConverter() cells.Behavior {
	processor := func(cell cells.Cell, event cells.Event) error {
		var rawCoins RawCoins
		err := event.Payload().Unmarshal(&rawCoins)
		if err != nil {
			return err
		}
		// Two helpers for trusted conversions.
		atof := func(a string) float64 {
			f, err := strconv.ParseFloat(a, 64)
			if err != nil {
				return 0.0
			}
			return f
		}
		atoi := func(a string) int {
			return int(atof(a))
		}
		// Convert and emit the coins.
		var coins Coins
		for _, rawCoin := range rawCoins {
			if rawCoin.MarketCapUSD == "" {
				continue
			}
			coin := Coin{
				Symbol:      rawCoin.Symbol,
				Name:        rawCoin.Name,
				PriceUSD:    atof(rawCoin.PriceUSD),
				TotalSupply: atoi(rawCoin.TotalSupply),
				LastUpdated: time.Unix(int64(atoi(rawCoin.LastUpdated)), 0),
			}
			coins = append(coins, coin)
		}
		cell.EmitNew("coins", coins)
		return nil
	}
	return behaviors.NewSimpleProcessorBehavior(processor)
}

// EOF
