// Tideland Go Cells - Example - Behaviors - Working on Slices of Coins
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
	"strings"

	"github.com/tideland/golib/identifier"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// BEHAVIORS
//--------------------

// MakeCoinsSpawnPointer returns a behavior checking the environment
// for behaviors running as spawn points for operations on individual
// coins.
func MakeCoinsSpawnPointer() cells.Behavior {
	processor := func(cell cells.Cell, event cells.Event) error {
		var coins Coins
		if err := event.Payload().Unmarshal(&coins); err != nil {
			return err
		}
		for _, coin := range coins {
			cellID := identifier.JoinedIdentifier("coin", strings.ToLower(coin.Symbol))
			if !cell.Environment().HasCell(cellID) {
				cell.Environment().StartCell(cellID, behaviors.NewBroadcasterBehavior())
			}
		}
		return nil
	}
	return behaviors.NewSimpleProcessorBehavior(processor)
}

// MakeCoinsSplitter returns a behavior splitting a list
// of coins into individual emits.
func MakeCoinsSplitter() cells.Behavior {
	processor := func(cell cells.Cell, event cells.Event) error {
		var coins Coins
		if err := event.Payload().Unmarshal(&coins); err != nil {
			return err
		}
		for _, coin := range coins {
			cell.EmitNew("coin", coin)
		}
		return nil
	}
	return behaviors.NewSimpleProcessorBehavior(processor)
}

// EOF
