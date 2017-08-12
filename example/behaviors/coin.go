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
	"strings"

	"github.com/tideland/golib/identifier"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// BEHAVIORS
//--------------------

// MakeRouter returns a behavior routing individual coin events to
// their spubscribed spown points.
func MakeRouter() cells.Behavior {
	// Router directly routes to spawn points, not to subscribers.
	router := func(cell cells.Cell, event cells.Event) error {
		var coin Coin
		if err := event.Payload().Unmarshal(&coin); err != nil {
			return err
		}
		cellID := identifier.JoinedIdentifier("coin", strings.ToLower(coin.Symbol))
		return cell.Environment().Emit(cellID, event)
	}
	return behaviors.NewSimpleProcessorBehavior(router)
}

// EOF
