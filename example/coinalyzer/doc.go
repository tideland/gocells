// Tideland Go Cells - Example - Documentation
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package coinalyzer implements the main program of the
// Tideland Go Cells example. The compile command line
// programm polls the feed at https://coinmarketcap.com/api/.
// This data is passed to the cells to analyze them.
//
// The analyzis focusses on price in USD and total supply, simply
// to show quicker changes. There are multiple categories of cells:
//
// 1. The general ones to do jobs like converting the raw data
// into our reduced coin format, to split lists of coins into
// individual ones, or to log events.
//
// 2. Those which are working on a number of coins, e.g. to
// calculate average values or simply count them.
//
// 3. Sometimes we're only interested in individual coins, to
// see how they evolve or how their changes match to other
// changes.
package main

//EOF
