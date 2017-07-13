// Tideland Go Cells - Behaviors
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package behaviors provides several generic and always useful
// standard behaviors for the Tideland Go Cells. They are simply
// created by calling NewXyzBehavior(). Their configuration
// is done by constructor arguments. Additionally some of them take
// functions or implementations of interfaces to control their
// processing. These behaviors are:
//
// Aggregator aggregates events and emits each aggregated value.
//
// Broadcaster simply emits received events to all subscribers.
//
// Callback calls a number of passed functions for each received event.
//
// Collector collects events which can be processed on demand.
//
// Combo waits for a user-defined combination of events.
//
// Condition tests events for conditions using a tester function
// and calls a processor then.
//
// Counter counts events, the counters can be retrieved.
//
// Evaluator evaluates events based on a user-defined function which
// returns a rating.
//
// Filter emits received events based on a user-defined filter.
//
// Finite State Machine allows to build finite state machines for events.
//
// Key/Value collects and emits payloads grouped by topics.
//
// Logger logs received events with level INFO.
//
// Mapper maps received events based on a user-defined function to
// new events.
//
// Once calls the once function only for the first event it receives.
//
// Pair checks if the event stream contains two matching ones based on a
// user-based criterion in a given timespan.
//
// Rate measures times between a number of criterion fitting events and
// emits the result.
//
// Rate Window checks if a number of events in a given timespan matches
// a given criterion.
//
// Round Robin distributes events round robin to its subscribers.
//
// Sequence checks the event stream for a defined sequence of events
// discovered by a user-defined criterion.
//
// Simple Processor allows to not implement a behavior but only use
// one function for event processing.
//
// Status receives and processes status events by other behaviors.
// Those have to emit it when receiving the topic "status" with a status
// cell ID as payload.
//
// Ticker emits tick events in a defined interval.
package behaviors

// EOF
