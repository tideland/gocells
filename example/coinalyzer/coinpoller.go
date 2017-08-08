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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tideland/golib/logger"
	"github.com/tideland/golib/loop"

	"github.com/tideland/gocells/cells"
	"github.com/tideland/gocells/example/behaviors"
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
	cfg  Configuration
	env  cells.Environment
	loop loop.Loop
}

// NewCoinPoller creates and starts a new coin poller goroutine.
func NewCoinPoller(cfg Configuration, env cells.Environment) *CoinPoller {
	cp := &CoinPoller{
		cfg: cfg,
		env: env,
	}
	cp.loop = loop.Go(cp.backendLoop, "cells-example-poller")
	return cp
}

// Stop tells the CoinPoller to stop working.
func (cp *CoinPoller) Stop() {
	cp.loop.Kill(nil)
}

// Wait waits until the coin poller stopped.
func (cp *CoinPoller) Wait() error {
	return cp.loop.Wait()
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
				logger.Errorf("cannot retrieve coin values: %v", err)
			}
		}
	}
}

// poll requests the coin values and pushes them into the
// cells.
func (cp *CoinPoller) poll() error {
	// Retrieve the current values.
	logger.Infof("polling from %s", PollURL)
	resp, err := http.Get(PollURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var rawCoins behaviors.RawCoins
	err = json.Unmarshal(body, &rawCoins)
	if err != nil {
		return err
	}
	// Pass the values to the cell environment.
	return cp.env.EmitNew("raw-coins", "raw-coins", rawCoins)
}

// EOF
