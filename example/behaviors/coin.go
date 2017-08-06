// Tideland Go Cells - Example - Behaviors - Working on individual Coins
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
	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
	"github.com/tideland/golib/identifier"
)

//--------------------
// BEHAVIORS
//--------------------

// MakeCoinEntryPoint returns a behavior checking the environment
// for a behavior running as entry point for operations on individual
// coins.
func MakeCoinEntryPoint() cells.Behavior {
	processor := func(cell cells.Cell, event cells.Event) error {
		coinID := event.Payload().String()
		cellID := identifier.JoinedIdentifier("coin", coinID)
		if !cell.Environment().HasCell(cellID) {
			cell.Environment().StartCell(cellID, behaviors.NewBroadcasterBehavior())
		}
		return nil
	}
	return behaviors.NewSimpleProcessorBehavior(processor)
}

// EOF
