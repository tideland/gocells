// Tideland Go Cells - Behaviors - Ticker
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors

//--------------------
// IMPORTS
//--------------------

import (
	"time"

	"github.com/tideland/gocells/cells"
	"github.com/tideland/golib/loop"
)

//--------------------
// CONSTANTS
//--------------------

const (
	// TopicTicker signals a tick event.
	TopicTicker = "tick!"

	// PayloadTickerID contains the ID of the ticker to differentiante
	// multiple ones.
	PayloadTickerID = "ticker:id"

	// PayloadTickerTime contains the time of the tick event.
	PayloadTickerTime = "ticker:time"
)

//--------------------
// TICKER BEHAVIOR
//--------------------

// tickerBehavior emits events in chronological order.
type tickerBehavior struct {
	cell     cells.Cell
	duration time.Duration
	loop     loop.Loop
}

// NewTickerBehavior creates a ticker behavior.
func NewTickerBehavior(duration time.Duration) cells.Behavior {
	return &tickerBehavior{
		duration: duration,
	}
}

// Init the behavior.
func (b *tickerBehavior) Init(c cells.Cell) error {
	b.cell = c
	b.loop = loop.Go(b.tickerLoop)
	return nil
}

// Terminate the behavior.
func (b *tickerBehavior) Terminate() error {
	return b.loop.Stop()
}

// PrecessEvent emits a ticker event each time the
// defined duration elapsed.
func (b *tickerBehavior) ProcessEvent(event cells.Event) error {
	if event.Topic() == TopicTicker {
		b.cell.EmitNew(TopicTicker, cells.Values{
			PayloadTickerID:   b.cell.ID(),
			PayloadTickerTime: time.Now(),
		}.Payload())
	}
	return nil
}

// Recover from an error. Counter will be set back to the initial counter.
func (b *tickerBehavior) Recover(err interface{}) error {
	return nil
}

// tickerLoop sends ticker events to its own process method.
func (b *tickerBehavior) tickerLoop(l loop.Loop) error {
	for {
		select {
		case <-l.ShallStop():
			return nil
		case now := <-time.After(b.duration):
			// Notify myself, action there to avoid
			// race when subscribers are updated.
			b.cell.Environment().EmitNew(b.cell.ID(), TopicTicker, cells.Values{
				cells.PayloadTickerTime: now,
			}.Payload())
		}
	}
}

// EOF
