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
	"sync"
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

// MakeTicker returns a ticker for the clocking of the
// emitting of the data.
func MakeTicker() cells.Behavior {
	return behaviors.NewTickerBehavior(100 * time.Millisecond)
}

// MakeEndOFData returns a behavior reacting when all
// data is done.
func MakeEndOfData(wg sync.WaitGroup) cells.Behavior {
	wg.Add(1)
	return nil
}

// EOF
