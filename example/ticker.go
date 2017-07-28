// Tideland Go Cells - Example - Ticker
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
)

//--------------------
// CONSTANTS
//--------------------

const (
	// PollInterval controls the duration between two polls,
	// CoinMarketCap wants it to be limited.
	PollInterval = 10 * time.Second
)

//--------------------
// TICKER
//--------------------

// Ticker polls the CoinMarketCap JSON API as input of
// the cells.
type Ticker struct {
}

// EOF
