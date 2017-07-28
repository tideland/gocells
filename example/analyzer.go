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
)

//--------------------
// ANALYZER
//--------------------

// Analyzer provides the testable core application of the
// example.
type Analyzer struct {
	ctx context.Context
}

// NewAnalyer creates a new analyzer instance.
func NewAnalyer(ctx context.Context) *Analyzer {
	logger.Infof("starting the Tideland Go Cells example analyzer")
	return &Analyzer{
		ctx: ctx,
	}
}

// Run performs the analyzing.
func (a *Analyzer) Run() error {
	logger.Infof("running the Tideland Go Cells example analyzer")

	return nil
}

// Cleanup tells the analyzer to remove temporary data,
// e.g. files.
func (a *Analyzer) Cleanup() {}

// EOF
