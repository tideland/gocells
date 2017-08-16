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
	"time"

	"github.com/tideland/golib/identifier"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// BEHAVIORS
//--------------------

// MakeRouter returns a behavior routing individual coin events to
// their subscribed spown points.
func MakeRouter() cells.Behavior {
	router := func(emitterID, subscriberID string, event cells.Event) (bool, error) {
		return event.Topic() == subscriberID, nil
	}
	return behaviors.NewRouterBehavior(router)
}

// MakeCoinRateWindow returns a rate window behavior for one coin
// looking for raises.
func MakeCoinRateWindow() cells.Behavior {
	current := 0.0
	criterion := func(event cells.Event) (bool, error) {
		var coin Coin
		if err := event.Payload().Unmarshal(&coin); err != nil {
			return false, err
		}
		raised := coin.PriceBTC > current
		current = coin.PriceBTC
		return raised, nil
	}
	return behaviors.NewRateWindowBehavior(criterion, 3, time.Minute)
}

//--------------------
// ENVIRONMENT
//--------------------

// SetupCoinEnvironment creates the environment for one coin.
func SetupCoinEnvironment(env cells.Environment, symbol string) error {
	// Broadcaster as spawn cell.
	spawnCellID := identifier.JoinedIdentifier("coin", symbol)
	if env.HasCell(spawnCellID) {
		return nil
	}
	env.StartCell(spawnCellID, behaviors.NewBroadcasterBehavior())
	// Coin rate window behavior.
	rateCellID := identifier.JoinedIdentifier("coin-rate-window", symbol)
	env.StartCell(rateCellID, MakeCoinRateWindow())
	// Subscriptions.
	env.Subscribe("router", spawnCellID)
	env.Subscribe(spawnCellID, rateCellID)
	env.Subscribe(rateCellID, "logger")
	// TODO(mue): More to come.
	return nil
}

// EOF
