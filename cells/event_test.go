// Tideland Go Cells - Unit Tests
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cells_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"github.com/tideland/golib/audit"
	"github.com/tideland/golib/errors"

	"github.com/tideland/gocells/cells"
)

//--------------------
// TESTS
//--------------------

// TestEvent tests the event construction.
func TestEvent(t *testing.T) {
	assert := audit.NewTestingAssertion(t, true)
	now := time.Now().UTC()

	event, err := cells.NewEvent("foo", "bar")
	assert.Nil(err)
	assert.True(event.Timestamp().After(now))
	assert.True(time.Now().UTC().After(event.Timestamp()))
	assert.Equal(event.Topic(), "foo")

	bar := event.Payload().String()
	assert.Equal(bar, "bar")

	_, err = cells.NewEvent("", nil)
	assert.True(errors.IsError(err, cells.ErrNoTopic))

	_, err = cells.NewEvent("yadda", nil)
	assert.Nil(err)
}

// TestPayload tests the payload creation and access.
func TestPayload(t *testing.T) {
	type loading struct {
		Bool     bool
		Int      int
		Float    float64
		String   string
		Time     time.Time
		Duration time.Duration
	}

	assert := audit.NewTestingAssertion(t, true)

	in := loading{
		Bool:     true,
		Int:      42,
		Float:    47.11,
		String:   "Hello, world!",
		Time:     time.Now(),
		Duration: 30 * time.Second,
	}
	payload, err := cells.NewPayload(in)
	assert.Nil(err)
	var out loading
	err = payload.Unmarshal(&out)
	assert.Nil(err)
	assert.Equal(in, out)

	payload, err = cells.NewPayload([]byte{1, 3, 3, 7})
	assert.Nil(err)
	bs := payload.Bytes()
	assert.Equal(bs, []byte{1, 3, 3, 7})

	same, err := cells.NewPayload(payload)
	assert.Nil(err)
	assert.Equal(same, payload)
}

// EOF
