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
	"github.com/tideland/golib/logger"
)

//--------------------
// ANALYZER
//--------------------

// Analyzer provides the testable core application of the
// example.
type Analyzer struct {
}

// NewAnalyer creates a new analyzer instance.
func NewAnalyer() *Analyzer {
	logger.Infof("Starting the Tideland Go Cells example analyzer ...")
	return &Analyzer{}
}

// Run performs the analyzing.
func (a *Analyzer) Run() error {
	return nil
}

// EOF
