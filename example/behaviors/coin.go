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
	"github.com/tideland/golib/logger"

	"github.com/tideland/gocells/behaviors"
	"github.com/tideland/gocells/cells"
)

//--------------------
// BEHAVIORS
//--------------------

// MakeRouter returns a behavior routing individual coin events to
// their spubscribed spown points.
func MakeRouter() cells.Behavior {
	router := func(cell cells.Cell, event cells.Event) error {
		println("ROUTING")
		var coin Coin
		if err := event.Payload().Unmarshal(&coin); err != nil {
			return err
		}
		cellID := identifier.JoinedIdentifier("coin", coin.Symbol)
		logger.Infof("ROUTING TO %v", cellID)
		return cell.Environment().Emit(cellID, event)
	}
	return behaviors.NewSimpleProcessorBehavior(router)
}

// MakeCoinRateWindow returns a rate window behavior for one coin
// looking for raises.
func MakeCoinRateWindow() cells.Behavior {
	price := 0.0
	criterion := func(event cells.Event) (bool, error) {
		var coin Coin
		if err := event.Payload().Unmarshal(&coin); err != nil {
			return false, err
		}
		logger.Infof("COIN: %v / PRICE: %v / PRICE USD: %v", coin.Symbol, price, coin.PriceUSD)
		if coin.PriceUSD > price {
			price = coin.PriceUSD
			return true, nil
		}
		return false, nil
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
	if !env.HasCell(spawnCellID) {
		env.StartCell(spawnCellID, behaviors.NewBroadcasterBehavior())
	}
	// Coin rate window behavior.
	rateCellID := identifier.JoinedIdentifier("coin-rate-window", symbol)
	env.StartCell(rateCellID, MakeCoinRateWindow())
	env.Subscribe(spawnCellID, rateCellID)
	env.Subscribe(rateCellID, "logger")
	// TODO(mue): More to come.
	return nil
}

// EOF
