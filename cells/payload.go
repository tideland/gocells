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
	"time"
)

//--------------------
// CONSTANTS
//--------------------

const (
	DefaultKey = "default"

	sliceSep        = " || "
	mapSep          = ":"
	timestampFormat = "2006-1-2-15-4-5.999999999 (MST)"
)

//--------------------
// PAYLOAD
//--------------------

// PayloadValues is intended to easily set the information
// of a payload.
type PayloadValues map[string]interface{}

// Payload is a write-once/read-multiple container for the
// transport of additional information with events.
type Payload interface {
	fmt.Stringer

	// Len returns the number of values.
	Len() int

	// Get returns one of the payload values.
	Get(key string) (string, bool)

	// GetStringSlice splits one of the payload value into
	// substrings and returns them as slice.
	GetStringSlice(key string) ([]string, bool)

	// GetStringMap splits one of the payload values into
	// parts and these again into key and value. Those are
	// returned as map.
	GetStringMap(key string) (map[string]string, bool)

	// GetBool returns one of the payload values as bool.
	GetBool(key string) (bool, bool)

	// GetInt returns one of the payload values as int.
	GetInt(key string) (int, bool)

	// GetUint returns one of the payload values as uint.
	GetUint(key string) (uint, bool)

	// GetFloat64 returns one of the payload values as float64.
	GetFloat64(key string) (float64, bool)

	// GetTime returns one of the payload values as time.
	GetTime(key string) (time.Time, bool)

	// GetDuration returns one of the payload values as duration.
	GetDuration(key string) (time.Duration, bool)

	// Do iterates a function over all keys and values.
	Do(f func(key, value string) error) error

	// Apply creates a new payload based on this one
	// and the passed values. If a key already exists its
	// value will be overwritten.
	Apply(values PayloadValues) Payload
}

// payload implements the Payload interface.
type payload struct {
	values map[string]string
}

// NewPayload creates a new payload containing the passed
// values.
func NewPayload(values PayloadValues) Payload {
	p := &payload{}
	return p.Apply(values)
}

// NewDefaultPayload creates a payload containing the key
// "default" with the passed values.
func NewDefaultPayload(value string) Payload {
	return NewPayload(PayloadValues{DefaultKey: value})
}

// Len implements Payload.
func (p *payload) Len() int {
	return len(p.values)
}

// Get implements Payload.
func (p *payload) Get(key string) (string, bool) {
	value, ok := p.values[key]
	return value, ok
}

// GetStringSlice implements Payload.
func (p *payload) GetStringSlice(key string) ([]string, bool) {
	raw, ok := p.values[key]
	if !ok {
		return nil, false
	}
	values := strings.Split(raw, sliceSep)
	return values, true
}

// GetStringMap implements Payload.
func (p *payload) GetStringMap(key string) (map[string]string, bool) {
	raw, ok := p.values[key]
	if !ok {
		return nil, false
	}
	parts := strings.Split(raw, sliceSep)
	values := make(map[string]string, len(parts))
	for _, part := range parts {
		kv := strings.SplitN(part, mapSep, 2)
		switch len(kv) {
		case 2:
			values[kv[0]] = kv[1]
		default:
			values[kv[0]] = kv[0]
		}
	}
	return values, true
}

// GetBool implements Payload.
func (p *payload) GetBool(key string) (bool, bool) {
	raw, ok := p.values[key]
	if !ok {
		return false, false
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false, false
	}
	return value, true
}

// GetInt implements Payload.
func (p *payload) GetInt(key string) (int, bool) {
	raw, ok := p.values[key]
	if !ok {
		return 0, false
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, false
	}
	return value, true
}

// GetUint implements Payload.
func (p *payload) GetUint(key string) (uint, bool) {
	raw, ok := p.values[key]
	if !ok {
		return 0, false
	}
	value, err := strconv.ParseUint(raw, 10, 0)
	if err != nil {
		return 0, false
	}
	return uint(value), true
}

// GetFloat64 implements Payload.
func (p *payload) GetFloat64(key string) (float64, bool) {
	raw, ok := p.values[key]
	if !ok {
		return 0.0, false
	}
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0.0, false
	}
	return value, true
}

// GetTime implements Payload.
func (p *payload) GetTime(key string) (time.Time, bool) {
	raw, ok := p.values[key]
	if !ok {
		return time.Time{}, false
	}
	value, err := time.Parse(timestampFormat, raw)
	if err != nil {
		return time.Time{}, false
	}
	return value, true
}

// GetDuration implements Payload.
func (p *payload) GetDuration(key string) (time.Duration, bool) {
	raw, ok := p.values[key]
	if !ok {
		return 0, false
	}
	value, err := time.ParseDuration(raw)
	if err != nil {
		return 0, false
	}
	return value, true
}

// Do implements Payload.
func (p *payload) Do(f func(key, value string) error) error {
	for k, v := range p.values {
		if err := f(k, v); err != nil {
			return err
		}
	}
	return nil
}

// Apply implements Payload.
func (p *payload) Apply(values PayloadValues) Payload {
	np := &payload{
		values: make(map[string]string),
	}
	for k, v := range p.values {
		np.values[k] = v
	}
	for k, rv := range values {
		switch v := rv.(type) {
		case string:
			np.values[k] = v
		case []string:
			np.values[k] = strings.Join(v, sliceSep)
		case map[string]string:
			vs := []string{}
			for mk, mv := range v {
				vs = append(vs, mk+mapSep+mv)
			}
			np.values[k] = strings.Join(vs, sliceSep)
		case time.Time:
			np.values[k] = v.Format(timestampFormat)
		default:
			np.values[k] = fmt.Sprintf("%v", v)
		}
	}
	return np
}

// String implements fmt.Stringer.
func (p *payload) String() string {
	ps := []string{}
	for k, v := range p.values {
		ps = append(ps, fmt.Sprintf("{'%s':'%s'}", k, v))
	}
	return "[" + strings.Join(ps, " ") + "]"
}

// EOF
