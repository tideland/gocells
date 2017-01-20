// Tideland Go Cells - Behaviors - Types
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
)

//--------------------
// TYPES
//--------------------

// EventData represents the pure collected event data.
type EventData struct {
	Timestamp time.Time
	Topic     string
	Payload   cells.Payload
}

// newEventData returns the passed event as event data to collect.
func newEventData(event cells.Event) EventData {
	data := EventData{
		Timestamp: time.Now(),
		Topic:     event.Topic(),
		Payload:   event.Payload(),
	}
	return data
}

// EventDatas stores a number of event datas.
type EventDatas struct {
	max   int
	datas []*EventData
}

// NewEventDatas creates a store for event datas.
func NewEventDatas(max int) *EventDatas {
	return &EventDatas{
		max: max,
	}
}

// Add adds a new event data based on the passed event.
func (d *EventDatas) Add(event cells.Event) *EventData {
	data := &EventData{
		Timestamp: time.Now(),
		Topic:     event.Topic(),
		Payload:   event.Payload(),
	}
	d.datas = append(d.datas, data)
	if d.max > 0 && len(d.datas) > d.max {
		d.datas = d.datas[1:]
	}
	return data
}

// Len returns the number of stored event datas.
func (d *EventDatas) Len() int {
	return len(d.datas)
}

// First returns the first of the collected event datas.
func (d *EventDatas) First() (*EventData, bool) {
	if len(d.datas) < 1 {
		return nil, false
	}
	return d.datas[0], true
}

// Last returns the last of the collected event datas.
func (d *EventDatas) Last() (*EventData, bool) {
	if len(d.datas) < 1 {
		return nil, false
	}
	return d.datas[len(d.datas)-1], true
}

// TimestampAt returns the collected timestamp at a given index.
func (d *EventDatas) TimestampAt(index int) (time.Time, bool) {
	if index < 0 || index > len(d.datas)-1 {
		return time.Time{}, false
	}
	return d.datas[index].Timestamp, true
}

// TopicAt returns the collected topic at a given index.
func (d *EventDatas) TopicAt(index int) (string, bool) {
	if index < 0 || index > len(d.datas)-1 {
		return "", false
	}
	return d.datas[index].Topic, true
}

// PayloadAt returns the collected payload at a given index.
func (d *EventDatas) PayloadAt(index int) (cells.Payload, bool) {
	if index < 0 || index > len(d.datas)-1 {
		return nil, false
	}
	return d.datas[index].Payload, true
}

// Do iterates over all collected event datas.
func (d *EventDatas) Do(f func(index int, data *EventData) error) error {
	for index, data := range d.datas {
		if err := f(index, data); err != nil {
			return err
		}
	}
	return nil
}

// EOF
