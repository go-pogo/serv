// Copyright (c) 2026, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httperr

import (
	"net/http"

	"github.com/go-pogo/errors"
	"github.com/go-pogo/serv/middleware"
)

// Wrap wraps the [http.HandlerFunc] next with middleware that catches any
// possible panics and returns them as the error result of the returned
// [Handler].
func Wrap(next http.HandlerFunc) Handler {
	if next == nil {
		panic(panicNilNextHandler)
	}
	return HandlerFunc(func(wri http.ResponseWriter, req *http.Request) (err error) {
		defer errors.CatchPanic(&err)
		next.ServeHTTP(wri, req)
		return
	})
}

// WrapWith wraps the [Handler] next with an [http.Handler] which is passed to
// the wrap function, which further wraps the next handlers. Any error returned
// by next is kept and eventually returned by the resulting [Handler].
func WrapWith(next Handler, wrap middleware.Wrapper) Handler {
	if next == nil {
		panic(panicNilNextHandler)
	}
	if wrap == nil {
		return next
	}
	return &wrapper{next: next, wrap: wrap}
}

var (
	_ http.Handler = (*wrapper)(nil)
	_ Handler      = (*wrapper)(nil)
)

type wrapper struct {
	next Handler
	wrap middleware.Wrapper
}

func (b *wrapper) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	b.wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := b.next.ServeHTTPError(w, r); err != nil {
			panic(err) // http.Server catches panics and handles the error
		}
	})).ServeHTTP(wri, req)
}

func (b *wrapper) ServeHTTPError(wri http.ResponseWriter, req *http.Request) error {
	var err error
	b.wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err = b.next.ServeHTTPError(w, r)
	})).ServeHTTP(wri, req)
	return err
}
