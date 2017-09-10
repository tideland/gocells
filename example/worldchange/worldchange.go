// Tideland Go Cells - Example - World Change
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

	"github.com/tideland/golib/logger"

	"github.com/tideland/gocells/cells"
)

//--------------------
// WORLD CHANGE
//--------------------

// WorldChange provides the testable core application of the
// example.
type WorldChange struct {
	cfg Configuration
	wg  sync.WaitGroup
	env cells.Environment
}

// NewWorldChange creates a new world change analyzer instance.
func NewWorldChange(cfg Configuration) *WorldChange {
	logger.Infof("starting the Tideland Go Cells World Change Analyzer")
	var wg sync.WaitGroup
	env, err := InitEnvironment(cfg, wg)
	if err != nil {
		logger.Fatalf("cannot init environment: %v", err)
		return nil
	}
	return &WorldChange{
		cfg: cfg,
		wg:  wg,
		env: env,
	}
}

// Run performs the analyzing.
func (w *WorldChange) Run() error {
	w.wg.Wait()
	return nil
}

// Cleanup tells the analyzer to remove temporary data,
// e.g. files.
func (w *WorldChange) Cleanup() {
	w.env.Stop()
}

// EOF
