# Tideland Go Cells

## Description

The *Tideland Go Cells* provide a package for the creation of event based
applications with networked concurrently working cells. The way how they
process the recevied events is defined by behaviors. During the processing
of an event a cell can emit multiple events to its subscribers.

I hope you like them. ;)

## Version

Version 6.0.0-beta.2017-01-22

## Packages

### Cells

Major package providing the infrastructure for event based applications.

[![GoDoc](https://godoc.org/github.com/tideland/gocells/cells?status.svg)](https://godoc.org/github.com/tideland/gocells/cells)

### Behaviors

The project already contains some standard behaviors, the number is
still growing.

### Broadcaster

The broadcaster behavior simply emits all received events to all
of its subscribers. It is intended to be used as a top level behavior
to directly rigger multiple handlers instead of emitting an event
manually to those handlers.

### Callback

The callback behavior allows you to provide a number of functions
which will be called when an event is received. Those functions
have the topic and the payload of the event as argument.

### Collector

The collector behavior collects all received events. They can be
retrieved and resetted. It also emits all received events to its
subscribers.

### Configurator

After receiving a ReadConfigurationTopic with a filename as
payload the configuration behavior reads this configuration
and emits it. If it is started with a validator the configuration
is validated after the reading.

### Counter

The counter behavior is created with a counter function as argument.
This function is called for each event and returns the IDs of counters
which are incremented then. The counters are emitted each time and
also can be resetted.

### Filter

The filter behavior is created with a filtering function which is
called for each event. If this function call returns true the event
emitted, otherwise it is dropped.

### Finite State Machine

The FSM behavior implements a finite state machine. State functions
process the events and return the following state function.

### Logger

The logger behavior logs every event. The used level is INFO.

### Mapper

The mapper behavior is created with a mapping. It is called with each
received event and returns a new mapped one.

### Pair

The rate behavior checks if two events match a criterion defined by
a user function with the signature `func(event cells.Event, hitData interface{}) (interface{}, bool)` 
happen during a defined duration. If the second event is in time a positive
pair event is emitted, otherwise a timeout pair event.

The returned data in case of a positive first hit is passed to each future
call of the match function allowing it to make its decision depending on
this data. The payload of a pair event contains both timestamps and both
returned data.

### Rate

The rate behavior measures the time between two events fitting a criterion
defined by a user function with the signature `func(event cells.Event) bool`.
The time, the duration, and the lowest, highest, and average duration during
a defined number of fitting events is emitted each time.

### Rate Window

The rate window behavior checks if an event matches a given criterion defined
by a user function with the signature `func(event cells.Event) bool`. If a
defined number of events match during a defined timespan an event is emitted.
It contains the first and the last time as well as the number of matches.
This way bursts of events can be detected.

### Round Robin

The round robin behavior distributes each received event round robin
to its subscribers. It can be used for load balancing.

### Sequence

The sequence behavior checks the event stream for a sequence defined
by the criterion function with the signature
`func(event cells.Event, collected *cells.EventDatas) (bool, bool)`.
Here `event` is the current event, `collected` the so far collected
matching events. If the current event is matching the first returned
bool has to be true, if all are collected the second one has to be
true too. In this case an event containing the collected events will
be emitted.

### Simple Processor

The simple behavior is created with a simple event processing function.
Useful if no state and no complex recovery is needed.

### Ticker

The ticker behavior emits a tick event in a defined interval to its
subscribers. So they can process chronological tasks beside other
events.

[![GoDoc](https://godoc.org/github.com/tideland/gocells/behaviors?status.svg)](https://godoc.org/github.com/tideland/gocells/behaviors)

## Contributors

- Frank Mueller (https://github.com/TheMue / https://github.com/tideland)
- Jonathan Camp (https://github.com/kung-foo)

## License

*Tideland Go Cells* is distributed under the terms of the BSD 3-Clause license.
