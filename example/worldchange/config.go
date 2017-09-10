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

// Configuration contains the sewttings for the example application.
type Configuration struct {
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
