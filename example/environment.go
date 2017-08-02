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

	"github.com/tideland/gocells/behaviors"

	"github.com/tideland/gocells/cells"
)

//--------------------
// ENVIRONMENT
//--------------------

// InitEnvironment creates a new cells environment and
// its behaviors and subscriptions.
func InitEnvironment(ctx context.Context) (cells.Environment, error) {
	env := cells.NewEnvironment("cells-example")

	// Start initial cells.
	env.StartCell("raw-coins", MakeRawCoinsConverter())
	env.StartCell("coins-splitter", MakeCoinsSplitter())
	env.StartCell("average-pc1h", MakeAveragePercentChange1hEvaluator())
	env.StartCell("average-pc1h-up", MakeAveragePercentChange1hRater(true))
	env.StartCell("average-pc1h-down", MakeAveragePercentChange1hRater(false))
	env.StartCell("logger", behaviors.NewLoggerBehavior())

	// Establish initial subscriptions.
	env.Subscribe("raw-coins", "coins-splitter")
	env.Subscribe("raw-coins", "average-pc1h")
	env.Subscribe("average-pc1h", "average-pc1h-up", "average-pc1h-down")
	env.Subscribe("average-pc1h-up", "logger")
	env.Subscribe("average-pc1h-down", "logger")

	return env, nil
}

// EOF
