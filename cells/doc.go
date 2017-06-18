// Tideland Go Cells
//
// Copyright (C) 2010-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package cells is the main package of the Tideland Cells. They provide
// a framework for event and behavior based applications.
//
// Cell behaviors are defined based on an interface and can be added
// to an environment. Here they are running as concurrent cells that
// can be networked and communicate via events. Several useful behaviors
// are provided with the behaviors package.
//
// New environments are created with
//
//     env := cells.NewEnvironment(identifier)
//
// and cells are added with
//
//    env.StartCell("foo", NewFooBehavior())
//
// Cells then can be subscribed with
//
//    env.Subscribe("foo", "bar")
//
// so that events emitted by the "foo" cell during the processing of
// events will be received by the "bar" cell. Each cell can have
// multiple cells subscibed.
//
// Events from the outside are emitted using
//
//     env.Emit("foo", myEvent)
//
// or
//
//    env.EmitNew("foo", "myTopic", cells.PayloadValues{
//        "KeyA": 12345,
//        "KeyB": true,
//    })
//
// Behaviors have to implement the cells.Behavior interface. Here
// the Init() method is called with a cells.Context. This can be
// used inside the ProcessEvent() method to emit events to subscribers
// or directly to other cells of the environment.
//
// Sometimes it's needed to directly communicate with a cell to retrieve
// information. In this case the method
//
//     response, err := env.RequestContext("foo", "myRequest?", myPayload, myTimeout)
//
// is to be used. Inside the ProcessEvent() of the addressed cell the
// event can be used to send the response with
//
//    switch event.Topic() {
//    case "myRequest?":
//        event.Respond(someIncredibleData)
//    case ...:
//        ...
//    }
//
// Instructions without a response are simply done by emitting an event.
package cells

// EOF
