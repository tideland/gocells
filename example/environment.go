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

	// Establish initial subscriptions.
	env.Subscribe("raw-coins", "coins-splitter")

	return env, nil
}

// EOF
