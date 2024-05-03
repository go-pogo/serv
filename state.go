// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"fmt"
)

type State uint32

const (
	// StateUnstarted is the default state of a [Server] when it is created.
	StateUnstarted State = iota
	// StateStarted indicates the [Server] is (almost) ready to start listening
	// for incoming connections.
	StateStarted
	// StateErrored indicates the [Server] has encountered an error while
	// listening for incoming connections.
	StateErrored
	// StateClosing indicates the [Server] is in the process of closing and
	// does not accept any incoming connections.
	StateClosing
	// StateClosed indicates the [Server] has been completely closed and is no
	// longer listening for incoming connections.
	StateClosed
)

func (s State) String() string {
	switch s {
	case StateUnstarted:
		return "unstarted"
	case StateStarted:
		return "started"
	case StateErrored:
		return "errored"
	case StateClosing:
		return "closing"
	case StateClosed:
		return "closed"

	default:
		panic(fmt.Sprintf("serv: %d is not a valid State", s))
	}
}

// InvalidStateError is returned when an operation is attempted on a [Server]
// that is in an invalid state for that operation to succeed.
type InvalidStateError struct {
	Err   error
	State State
}

func (u InvalidStateError) Unwrap() error { return u.Err }

func (u InvalidStateError) Error() string {
	return "unexpected state " + u.State.String()
}
