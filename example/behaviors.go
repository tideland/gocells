// Tideland Go Cells - Example - Cells Behaviors or Processing Functions
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
	"strconv"
	"time"

	"github.com/tideland/gocells/cells"
)

//--------------------
// PROCESSING FUNCTIONS
//--------------------

// SplitRawCoins splits the raw coins into individual ones and
// maps them to working coins.
func SplitRawCoins(cell cells.Cell, event cells.Event) error {
	var rawCoins RawCoins
	err := event.Payload().Unmarshal(&rawCoins)
	if err != nil {
		return err
	}
	// Two helpers for trusted conversions.
	atoi := func(a string) int {
		i, err := strconv.Atoi(a)
		if err != nil {
			return 0
		}
		return i
	}
	atof := func(a string) float64 {
		f, err := strconv.ParseFloat(a, 64)
		if err != nil {
			return 0.0
		}
		return f
	}
	// Convert and emit the coins.
	for _, rawCoin := range rawCoins {
		coin := Coin{
			ID:               rawCoin.ID,
			Name:             rawCoin.Name,
			Symbol:           rawCoin.Symbol,
			Rank:             atoi(rawCoin.Rank),
			PriceUSD:         atof(rawCoin.PriceUSD),
			PriceBTC:         atof(rawCoin.PriceBTC),
			Volume24hUSD:     atof(rawCoin.Volume24hUSD),
			MarketCapUSD:     atof(rawCoin.MarketCapUSD),
			AvailableSupply:  atoi(rawCoin.AvailableSupply),
			TotalSupply:      atoi(rawCoin.TotalSupply),
			PercentChange1h:  atof(rawCoin.PercentChange1h),
			PercentChange24h: atof(rawCoin.PercentChange24h),
			PercentChange7d:  atof(rawCoin.PercentChange7d),
			LastUpdated:      time.Unix(int64(atoi(rawCoin.LastUpdated)), 0),
		}
		cell.EmitNew("coin", coin)
	}
	return nil
}

// EOF
