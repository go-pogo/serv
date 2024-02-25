// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
)

// Wrapper wraps a http.HandlerFunc with additional logic and returns a new
// http.Handler.
type Wrapper interface {
	Wrap(next http.HandlerFunc) http.Handler
}

var _ Wrapper = (*WrapperFunc)(nil)

//goland:noinspection GoNameStartsWithPackageName
type WrapperFunc func(next http.HandlerFunc) http.Handler

func (fn WrapperFunc) Wrap(next http.HandlerFunc) http.Handler { return fn(next) }

// Wrap http.Handler h with the provided Wrapper.
func Wrap(handler http.Handler, wrap ...Wrapper) http.Handler {
	for i := len(wrap) - 1; i >= 0; i-- {
		handler = wrap[i].Wrap(handler.ServeHTTP)
	}
	return handler
}

var _ Wrapper = (*Middleware)(nil)

type Middleware []Wrapper

func (m Middleware) Wrap(next http.HandlerFunc) http.Handler {
	return Wrap(next, m...)
}

var _ Wrapper = (*HandlerFunc)(nil)

// HandlerFunc is a http.HandlerFunc which implements the Wrapper interface.
type HandlerFunc http.HandlerFunc

func (fn HandlerFunc) Wrap(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		fn(wri, req)
		next(wri, req)
	})
}
