// Tideland Go Cells - Unit Tests - Behaviors
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
	"time"

	"github.com/tideland/gocells/cells"
)

//--------------------
// TOPICS
//--------------------

const (
	// iterateTopic lets the test behavior iterate over its subscribers.
	iterateTopic = "iterate!"

	// panicTopic lets the test behavior panic to check recovering.
	panicTopic = "panic!"

	// subscribersTopic returns the current subscribers.
	subscribersTopic = "subscribers?"

	// emitTopic tells the cell to emit a test event.
	emitTopic = "emit!"

	// sleepTopic lets the cell sleep for a longer time so the queue gets full.
	sleepTopic = "sleep!"
)

//--------------------
// TEST BEHAVIORS
//--------------------

// nullBehavior does nothing.
type nullBehavior struct{}

var _ cells.Behavior = (*nullBehavior)(nil)

func (b *nullBehavior) Init(c cells.Cell) error { return nil }

func (b *nullBehavior) Terminate() error { return nil }

func (b *nullBehavior) ProcessEvent(event cells.Event) error { return nil }

func (b *nullBehavior) Recover(r interface{}) error { return nil }

// collectBehavior collects and re-emits all events, returns them
// on the topic "processed" and delets all collected on the
// topic "reset".
type collectBehavior struct {
	cell        cells.Cell
	sink        cells.EventSink
	recoverings int
}

var _ cells.Behavior = (*collectBehavior)(nil)

func newCollectBehavior(sink cells.EventSink) *collectBehavior {
	return &collectBehavior{nil, sink, 0}
}

func (b *collectBehavior) Init(cell cells.Cell) error {
	b.cell = cell
	return nil
}

func (b *collectBehavior) Terminate() error {
	return nil
}

func (b *collectBehavior) ProcessEvent(event cells.Event) error {
	switch event.Topic() {
	case cells.ProcessedTopic:
		payload, ok := cells.HasWaiterPayload(event)
		if !ok {
			panic("illegal payload, need waiter")
		}
		payload.GetWaiter().Set(cells.NewPayload(b.sink))
	case cells.ResetTopic:
		b.sink.Clear()
	case iterateTopic:
		err := b.cell.SubscribersDo(func(s cells.Subscriber) error {
			return s.ProcessNewEvent(event.Context(), "love", b.cell.ID()+" loves "+s.ID())
		})
		if err != nil {
			return err
		}
	case panicTopic:
		panic("Ouch!")
	case subscribersTopic:
		var ids []string
		b.cell.SubscribersDo(func(s cells.Subscriber) error {
			ids = append(ids, s.ID())
			return nil
		})
		payload, ok := cells.HasWaiterPayload(event)
		if !ok {
			panic("illegal payload, need waiter")
		}
		payload.GetWaiter().Set(cells.NewPayload(ids))
	default:
		b.sink.Push(event)
		return b.cell.Emit(event)
	}
	return nil
}

func (b *collectBehavior) Recover(r interface{}) error {
	b.recoverings++
	if b.recoverings > 5 {
		return cells.NewCannotRecoverError(b.cell.ID(), r)
	}
	return nil
}

// eventBufferBehavior allows testing the setting
// of the event buffer size.
type testEventBufferBehavior struct {
	*collectBehavior

	size int
}

var _ cells.BehaviorEventBufferSize = (*testEventBufferBehavior)(nil)

func newEventBufferBehavior(size int, sink cells.EventSink) cells.Behavior {
	return &testEventBufferBehavior{
		collectBehavior: newCollectBehavior(sink),
		size:            size,
	}
}

func (b *testEventBufferBehavior) EventBufferSize() int {
	return b.size
}

// recoveringFrequencyBehavior allows testing the setting
// of the recovering frequency.
type recoveringFrequencyBehavior struct {
	*collectBehavior

	number   int
	duration time.Duration
}

var _ cells.BehaviorRecoveringFrequency = (*recoveringFrequencyBehavior)(nil)

func newRecoveringFrequencyBehavior(number int, duration time.Duration, sink cells.EventSink) cells.Behavior {
	return &recoveringFrequencyBehavior{
		collectBehavior: newCollectBehavior(sink),
		number:          number,
		duration:        duration,
	}
}

func (b *recoveringFrequencyBehavior) RecoveringFrequency() (int, time.Duration) {
	return b.number, b.duration
}

// emitTimeoutBehavior allows testing the setting
// of the emit timeout time.
type emitTimeoutBehavior struct {
	*collectBehavior

	timeout time.Duration
}

var _ cells.BehaviorEmitTimeout = (*emitTimeoutBehavior)(nil)

func newEmitTimeoutBehavior(timeout time.Duration, sink cells.EventSink) cells.Behavior {
	return &emitTimeoutBehavior{
		collectBehavior: newCollectBehavior(sink),
		timeout:         timeout,
	}
}

func (b *emitTimeoutBehavior) EmitTimeout() time.Duration {
	return b.timeout
}

// emitBehavior simply emits the sleep topic to its subscribers.
type emitBehavior struct {
	c cells.Cell
}

var _ cells.Behavior = (*emitBehavior)(nil)

func newEmitBehavior() *emitBehavior {
	return &emitBehavior{}
}

func (b *emitBehavior) Init(c cells.Cell) error {
	b.c = c
	return nil
}

func (b *emitBehavior) Terminate() error {
	return nil
}

func (b *emitBehavior) ProcessEvent(event cells.Event) error {
	return b.c.EmitNew(event.Context(), sleepTopic, event.Payload())
}

func (b *emitBehavior) Recover(r interface{}) error {
	return nil
}

// sleepBehavior simply emits the sleep topic to its subscribers.
type sleepBehavior struct {
	cell cells.Cell
}

var _ cells.Behavior = (*sleepBehavior)(nil)

func newSleepBehavior() *sleepBehavior {
	return &sleepBehavior{}
}

func (b *sleepBehavior) Init(c cells.Cell) error {
	b.cell = c
	return nil
}

func (b *sleepBehavior) Terminate() error {
	return nil
}

func (b *sleepBehavior) ProcessEvent(event cells.Event) error {
	time.Sleep(4 * time.Second)
	return nil
}

func (b *sleepBehavior) Recover(r interface{}) error {
	return nil
}

func (b *sleepBehavior) EmitTimeout() time.Duration {
	return 2 * time.Second
}

// EOF
