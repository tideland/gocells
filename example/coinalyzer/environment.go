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
	"github.com/tideland/gocells/cells"
	"github.com/tideland/gocells/example/behaviors"
)

//--------------------
// ENVIRONMENT
//--------------------

// InitEnvironment creates a new cells environment and
// its behaviors and subscriptions.
func InitEnvironment(cfg Configuration) (cells.Environment, error) {
	env := cells.NewEnvironment("cells-example")

	// Start initial cells.
	env.StartCell("logger", behaviors.MakeLogger())
	env.StartCell("raw-coins-converter", behaviors.MakeRawCoinsConverter())
	env.StartCell("coins-spawn-pointer", behaviors.MakeCoinsSpawnPointer())
	env.StartCell("coins-averager", behaviors.MakeCoinsAverager())
	env.StartCell("coins-splitter", behaviors.MakeCoinsSplitter())
	env.StartCell("router", behaviors.MakeRouter())

	// Establish initial subscriptions.
	env.Subscribe("raw-coins-converter", "coins-spawn-pointer", "coins-averager", "coins-splitter")
	env.Subscribe("coins-splitter", "router")

	return env, nil
}

// EOF
