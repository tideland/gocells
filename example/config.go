// Tideland Go Cells - Example - Configuration
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
)

//--------------------
// CONFIGURATION
//--------------------

// Configuration contains the settings controlling the analyzer.
type Configuration struct {
	// Limit controls a limitation of the top retrieved ticker results.
	Limit int

	// Convert controls a convertation of price, 24h volume, and market
	// cap of the ticker results into another currency. Possible are
	// "AUD", "BRL", "CAD", "CHF", "CNY", "EUR", "GBP", "HKD", "IDR",
	// "INR", "JPY", "KRW", "MXN", "RUB".
	Convert string
}

//--------------------
// CONTEXT
//--------------------

// contextKey is used to type keys for context values.
type contextKey int

// configurationKey addresses a configuration inside a context.
const configurationKey contextKey = 1

// NewContext returns a new context that carries a configuration.
func NewContext(ctx context.Context, cfg Configuration) context.Context {
	return context.WithValue(ctx, configurationKey, cfg)
}

// FromContext returns the configuration stored in ctx, if any.
func FromContext(ctx context.Context) (Configuration, bool) {
	cfg, ok := ctx.Value(configurationKey).(Configuration)
	return cfg, ok
}

// EOS
