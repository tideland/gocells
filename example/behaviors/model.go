// Tideland Go Cells - Example - Behaviors - Model
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
)

//--------------------
// COIN
//--------------------

// RawCoin contains one raw coin object of the API.
type RawCoin struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	Rank             string `json:"rank"`
	PriceUSD         string `json:"price_usd"`
	PriceBTC         string `json:"price_btc"`
	Volume24hUSD     string `json:"24h_volume_usd"`
	MarketCapUSD     string `json:"market_cap_usd"`
	AvailableSupply  string `json:"available_supply"`
	TotalSupply      string `json:"total_supply"`
	PercentChange1h  string `json:"percent_change_1h"`
	PercentChange24h string `json:"percent_change_24h"`
	PercentChange7d  string `json:"percent_change_7d"`
	LastUpdated      string `json:"last_updated"`
}

// RawCoins contains a list of raw coins
type RawCoins []RawCoin

// Coin is one converted raw coin reduced to the minimum.
type Coin struct {
	Symbol      string
	Name        string
	PriceUSD    float64
	TotalSupply int
	LastUpdated time.Time
}

// Coins contains a list of coins
type Coins []Coin

// EOF
