// Tideland Go Cells - Example - Coin Poller
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
	"time"

	"github.com/tideland/golib/loop"
)

//--------------------
// CONSTANTS
//--------------------

const (
	// PollInterval controls the duration between two polls,
	// CoinMarketCap wants it to be limited.
	PollInterval = 10 * time.Second

	// PollURL references the CoinMarketCap JSON API.
	PollURL = "https://api.coinmarketcap.com/v1/ticker/"
)

//--------------------
// COIN POLLER
//--------------------

// CoinPoller polls the CoinMarketCap JSON API as input of
// the cells.
type CoinPoller struct {
	ctx  context.Context
	loop loop.Loop
}

// NewCoinPoller creates and starts a new coin poller goroutine.
func NewCoinPoller(ctx context.Context) *CoinPoller {
	cp := &CoinPoller{
		ctx: ctx,
	}
	cp.loop = loop.Go(cp.backendLoop, "coin poller")
	return cp
}

// Stop tells the CoinPoller to stop working.
func (cp *CoinPoller) Stop() error {
	return cp.loop.Stop()
}

// backendLoop implements the gorouting of the coin poller.
func (cp *CoinPoller) backendLoop(l loop.Loop) error {
	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-l.ShallStop():
			return nil
		case <-ticker.C:
			if err := cp.poll(); err != nil {
				return err
			}
		}
	}
}

// poll requests the coin values and pushes them into the
// cells.
func (cp *CoinPoller) poll() error {
	return nil
}

// EOF
