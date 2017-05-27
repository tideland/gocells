// Tideland Go Cells - Payload
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cells

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"strconv"
	"strings"
)

//--------------------
// CONSTANTS
//--------------------

const (
	DefaultKey = "default"
)

//--------------------
// PAYLOAD
//--------------------

// PayloadValues is intended to set and get the information
// of a payload as bulk.
type PayloadValues map[string]string

// Payload creates a new payload out of the values.
func (pvs PayloadValues) Payload() Payload {
	return NewPayload(pvs)
}

// Payload is a write-once/read-multiple container for the
// transport of additional information with events.
type Payload interface {
	fmt.Stringer

	// Len returns the number of values.
	Len() int

	// Get returns one of the payload values.
	Get(key string) (string, bool)

	// Do iterates a function over all keys and values.
	Do(f func(key, value string) error) error

	// Apply creates a new payload based on this one
	// and the passed values. If a key already exists its
	// value will be overwritten.
	Apply(values PayloadValues) Payload
}

// payload implements the Payload interface.
type payload struct {
	values PayloadValues
}

// NewPayload creates a new payload containing the passed
// values.
func NewPayload(values PayloadValues) Payload {
	p := &payload{
		values: PayloadValues{},
	}
	for k, v := range values {
		p.values[k] = v
	}
	return p
}

// NewDefaultPayload creates a payload containing the key
// "default" with the passed values.
func NewDefaultPayload(value string) Payload {
	return NewPayload(PayloadValues{"default": value})
}

// Len implementes the Payload interface.
func (p *payload) Len() int {
	return len(p.values)
}

// Get implementes the Payload interface.
func (p *payload) Get(key string) (string, bool) {
	value, ok := p.values[key]
	return value, ok
}

// Do implementes the Payload interface.
func (p *payload) Do(f func(key, value string) error) error {
	for k, v := range p.values {
		if err := f(k, v); err != nil {
			return err
		}
	}
	return nil
}

// Apply implementes the Payload interface.
func (p *payload) Apply(values PayloadValues) Payload {
	np := &payload{
		values: PayloadValues{},
	}
	for k, v := range p.values {
		np.values[k] = v
	}
	for k, v := range values {
		np.values[k] = v
	}
	return np
}

// String implements the fmt.Stringer interface.
func (p *payload) String() string {
	ps := []string{}
	for k, v := range p.values {
		ps = append(ps, fmt.Sprintf("{'%s':'%s'}", k, v))
	}
	return "[" + strings.Join(ps, " ") + "]"
}

// EOF
