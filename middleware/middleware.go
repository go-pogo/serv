// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
)

// Middleware wraps a http.HandlerFunc with additional logic and returns a new
// http.Handler.
type Middleware interface {
	Wrap(next http.HandlerFunc) http.Handler
}

var (
	_ Middleware = new(MiddlewareFunc)
	_ Middleware = new(HandlerFunc)
)

//goland:noinspection GoNameStartsWithPackageName
type MiddlewareFunc func(next http.HandlerFunc) http.Handler

func (fn MiddlewareFunc) Wrap(next http.HandlerFunc) http.Handler { return fn(next) }

type HandlerFunc http.HandlerFunc

func (fn HandlerFunc) Wrap(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		fn(wri, req)
		next(wri, req)
	})
}

// Wrap http.Handler h with the provided Middleware.
func Wrap(handler http.Handler, wrap ...Middleware) http.Handler {
	for i := len(wrap) - 1; i >= 0; i-- {
		handler = wrap[i].Wrap(handler.ServeHTTP)
	}
	return handler
}
