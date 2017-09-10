// Tideland Go Cells - Example - Behaviors - Model
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// IMPORTS
//--------------------

//--------------------
// POPULATION
//--------------------

// Population contains the population of a country or region
// in one year.
type Population struct {
	CountryName string
	CountryCode string
	Year        string
	Value       int
}

// Populations is the set of all populations.
type Popoulations []Population

// EOF
