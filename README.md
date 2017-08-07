# Tideland Go Cells

## Description

The *Tideland Go Cells* provide a package for the creation of event based
applications with networked concurrently working cells. The way how they
process the recevied events is defined by behaviors. During the processing
of an event a cell can emit multiple events to its subscribers.

I hope you like them. ;)

[![GitHub release](https://img.shields.io/github/release/tideland/gocells.svg)](https://github.com/tideland/gocells)
[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/tideland/gocells/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/tideland/gocells?status.svg)](https://godoc.org/github.com/tideland/gocells)
[![Sourcegraph](https://sourcegraph.com/github.com/tideland/gocells/-/badge.svg)](https://sourcegraph.com/github.com/tideland/gocells?badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/tideland/gocells)](https://goreportcard.com/report/github.com/tideland/gocells)

## Version

Version 6.0.0-beta.2017-08-07

## Packages

### Cells

Major package providing the infrastructure for event based applications. It
is organized as an environment of networked cells. Each cell is controlled
by its behavior implementing the according interface. It receives events,
processes them, and can emit a number of new events during this time. Those
events are then received by those cells which subscribed to the individual
cell.

For the implementation of own behaviors the `EventSink`, the `EventSinkAccessor`,
and the `EventSinkAnalyzer` provide help for their tasks.

### Behaviors

The project already contains some standard behaviors, the number is
still growing.

- **Aggregator** aggregates events and emits each aggregated value.
- **Broadcaster** simply emits received events to all subscribers.
- **Callback** calls a number of passed functions for each received event.
- **Collector** collects events which can be processed on demand.
- **Combo** waits for a user-defined combination of events.
- **Condition** tests events for conditions using a tester function
  and calls a processor then.
- **Countdown** counts a number of events down to zero and executes an
  event returning function. The event will be emitted then.
- **Counter** counts events, the counters can be retrieved.
- **Evaluator** evaluates events based on a user-defined function which
  returns a rating.
- **Filter** re-emits received events based on a user-defined filter. It can
  be selective or excluding.
- **Finite State Machine** allows to build finite state machines for events.
- **Key/Value** collects and emits payloads grouped by topics.
- **Logger** logs received events with level INFO.
- **Mapper** maps received events based on a user-defined function to new events.
- **Once** calls the once function only for the first event it receives.
- **Pair** checks if the event stream contains two matching ones based on a
  user-based criterion in a given timespan.
- **Rate** measures times between a number of criterion fitting events and
  emits the result.
- **Rate Window** checks if a number of events in a given timespan matches
  a given criterion.
- **Round Robin** distributes events round robin to its subscribers.
- **Sequence** checks the event stream for a defined sequence of events
  discovered by a user-defined criterion.
- **Simple Processor** allows to not implement a behavior but only use
  one function for event processing.
- **Status** receives and processes status events by other behaviors.
  Those have to emit it when receiving the topic "status" with a status
  cell ID as payload.
- **Ticker** emits tick events in a defined interval.

## Contributors

- Frank Mueller (https://github.com/TheMue / https://github.com/tideland)
- Jonathan Camp (https://github.com/kung-foo)

## License

*Tideland Go Cells* is distributed under the terms of the BSD 3-Clause license.
