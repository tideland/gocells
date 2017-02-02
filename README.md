# Tideland Go Cells

## Description

The *Tideland Go Cells* provide a package for the creation of event based
applications with networked concurrently working cells. The way how they
process the recevied events is defined by behaviors. During the processing
of an event a cell can emit multiple events to its subscribers.

I hope you like them. ;)

## Version

Version 6.0.0-beta.2017-02-02

## Packages

### Cells

Major package providing the infrastructure for event based applications. It
is organized as an environment of networked cells. Each cell is controlled
by its behavior. It receives events, processes them, and can emit a number
of new events during this time. Those events are then received by those
cells which subscribed to the individual cell.

[![GoDoc](https://godoc.org/github.com/tideland/gocells/cells?status.svg)](https://godoc.org/github.com/tideland/gocells/cells)

### Behaviors

The project already contains some standard behaviors, the number is
still growing.

- **Broadcaster** simply emits received events to all subscribers.
- **Callback** calls a number of passed functions for each received event.
- **Collector** collects events, theese can be retrieved and reset.
- **Combo** waits for a user-defined combination of events.
- **Configurator** reads a configuration file based on an event and emits it.
- **Counter** counts events, the counters can be retrieved.
- **Evaluator** evaluates events based on a user-defined function which
  returns a rating.
- **Filter** emits received events based on a user-defined filter.
- **Finite State Machine** allows to build finite state machines for events.
- **Logger** logs received events with level INFO.
- **Mapper** maps received events based on a user-defined function to new events.
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
- **Ticker** emits tick events in a defined interval.
- **Waiter** sets the payload of the first received event to a payload waiter.

[![GoDoc](https://godoc.org/github.com/tideland/gocells/behaviors?status.svg)](https://godoc.org/github.com/tideland/gocells/behaviors)

## Contributors

- Frank Mueller (https://github.com/TheMue / https://github.com/tideland)
- Jonathan Camp (https://github.com/kung-foo)

## License

*Tideland Go Cells* is distributed under the terms of the BSD 3-Clause license.
