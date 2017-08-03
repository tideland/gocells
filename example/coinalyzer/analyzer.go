// Tideland Go Cells - Example - Analyzer
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

	"github.com/tideland/golib/logger"

	"github.com/tideland/gocells/cells"
)

//--------------------
// ANALYZER
//--------------------

// Analyzer provides the testable core application of the
// example.
type Analyzer struct {
	ctx context.Context
	env cells.Environment
}

// NewAnalyer creates a new analyzer instance.
func NewAnalyer(ctx context.Context) *Analyzer {
	logger.Infof("starting the Tideland Go Cells Coinalyzer")
	env, err := InitEnvironment(ctx)
	if err != nil {
		logger.Fatalf("cannot init environment: %v", err)
		return nil
	}
	return &Analyzer{
		ctx: ctx,
		env: env,
	}
}

// Run performs the analyzing.
func (a *Analyzer) Run() error {
	cp := NewCoinPoller(a.ctx, a.env)

	return cp.Wait()
}

// Cleanup tells the analyzer to remove temporary data,
// e.g. files.
func (a *Analyzer) Cleanup() {
	a.env.Stop()
}

// EOF
