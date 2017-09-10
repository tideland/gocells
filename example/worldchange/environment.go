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
	"sync"

	"github.com/tideland/gocells/cells"
	"github.com/tideland/gocells/example/behaviors"
)

//--------------------
// ENVIRONMENT
//--------------------

// InitEnvironment creates a new cells environment and
// its behaviors and subscriptions.
func InitEnvironment(cfg Configuration, wg sync.WaitGroup) (cells.Environment, error) {
	env := cells.NewEnvironment("world-change-analyzer")

	// Start initial cells.
	env.StartCell("logger", behaviors.MakeLogger())
	env.StartCell("ticker", behaviors.MakeTicker())
	env.StartCell("eod", behaviors.MakeEndOfData(wg))

	// Establish initial subscriptions.

	return env, nil
}

// EOF
