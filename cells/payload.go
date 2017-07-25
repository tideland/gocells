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
	"encoding/json"
	"fmt"

	"github.com/tideland/golib/errors"
)

//--------------------
// PAYLOAD
//--------------------

// Payload is a write-once/read-multiple container for the
// transport of additional information with events.
type Payload interface {
	fmt.Stringer

	// Len returns the size of the payload.
	Len() int

	// Bytes returns the raw payload bytes.
	Bytes() []byte

	// Unmarshal parses the JSON-encoded payload bytes and
	// stores the result in the value pointed to by v.
	Unmarshal(v interface{}) error
}

// payload implements the Payload interface.
type payload struct {
	Data []byte
}

// NewPayload creates a new payload based on the passed value. In
// case of a byte slice this is taken directly, otherwise it is
// marshalled into JSON.
func NewPayload(v interface{}) (Payload, error) {
	var data []byte
	var err error
	switch tv := v.(type) {
	case []byte:
		data = make([]byte, len(tv))
		copy(data, tv)
	case string:
		data = []byte(tv)
	case Payload:
		return tv, nil
	default:
		if v != nil {
			data, err = json.Marshal(v)
			if err != nil {
				return nil, errors.Annotate(err, ErrMarshal, errorMessages)
			}
		}
	}
	return &payload{
		Data: data,
	}, nil
}

// newEmptyPayload returns a payload with empty data for
// lazy access to events with no payload.
func newEmptyPayload() Payload {
	return &payload{
		Data: []byte{},
	}
}

// Len implements Payload.
func (p *payload) Len() int {
	return len(p.Data)
}

// Bytes implements Payload.
func (p *payload) Bytes() []byte {
	data := make([]byte, len(p.Data))
	copy(data, p.Data)
	return data
}

// Unmarshal implements Payload.
func (p *payload) Unmarshal(v interface{}) error {
	err := json.Unmarshal(p.Data, v)
	if err != nil {
		return errors.Annotate(err, ErrUnmarshal, errorMessages)
	}
	return nil
}

// String implements fmt.Stringer.
func (p *payload) String() string {
	return string(p.Data)
}

// EOF
