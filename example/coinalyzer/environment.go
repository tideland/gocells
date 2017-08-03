// Tideland Go Cells - Example - Cells Environment
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
	"context"

	"github.com/tideland/gocells/cells"
	"github.com/tideland/gocells/example/behaviors"
)

//--------------------
// ENVIRONMENT
//--------------------

// InitEnvironment creates a new cells environment and
// its behaviors and subscriptions.
func InitEnvironment(ctx context.Context) (cells.Environment, error) {
	env := cells.NewEnvironment("cells-example")

	// Start initial cells.
	env.StartCell("raw-coins", behaviors.MakeRawCoinsConverter())
	env.StartCell("coins-splitter", behaviors.MakeCoinsSplitter())
	env.StartCell("avg-percent-change-1h", behaviors.MakeAvgPercentChange1h())
	env.StartCell("top-coins", behaviors.MakeTopPercentChange1hCoins())
	env.StartCell("avg-marketcap", behaviors.MakeAvgMarketCapEvaluator())
	env.StartCell("avg-marketcap-up", behaviors.MakeAvgMarketCapRater(true))
	env.StartCell("avg-marketcap-down", behaviors.MakeAvgMarketCapRater(false))
	env.StartCell("logger", behaviors.MakeLogger())

	// Establish initial subscriptions.
	env.Subscribe("raw-coins", "coins-splitter")

	// PercentChange1h analysis.
	env.Subscribe("raw-coins", "avg-percent-change-1h", "top-coins")
	env.Subscribe("avg-percent-change-1h", "top-coins")
	env.Subscribe("top-coins", "logger")

	// MarketCap analysis.
	env.Subscribe("raw-coins", "avg-marketcap")
	env.Subscribe("avg-marketcap", "avg-marketcap-up", "avg-marketcap-down")
	env.Subscribe("avg-marketcap-up", "logger")
	env.Subscribe("avg-marketcap-down", "logger")

	return env, nil
}

// EOF
