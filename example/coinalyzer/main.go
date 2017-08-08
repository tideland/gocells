// Tideland Go Cells - Example - Main
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
// INITIALIZATION
//--------------------

// initConfiguration prepares the configuration based on the arguments.
func initConfiguration() (Configuration, error) {
	cfg := Configuration{}
	return cfg, nil
}

//--------------------
// MAIN
//--------------------

// main is, guess what, the main programm. Currently in a very
// temporary state.
func main() {
	cfg, err := initConfiguration()
	if err != nil {
		logger.Fatalf("cannot run example: %v", err)
	}
	analyzer := NewAnalyzer(cfg)
	err = analyzer.Run()
	if err != nil {
		logger.Errorf("analyzer stopped with error: %v", err)
	}
	analyzer.Cleanup()
}

// EOF
