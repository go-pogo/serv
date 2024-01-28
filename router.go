// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import "net/http"

// RouteHandler handles routes.
type RouteHandler interface {
	Handle(method, pattern string, handler http.Handler)
	HandleFunc(method, pattern string, handler http.HandlerFunc)
}

// RoutesRegisterer registers routes to a RouteHandler.
type RoutesRegisterer interface {
	RegisterRoutes(r RouteHandler)
}

// RoutesRegistererFunc registers routes to a RouteHandler.
type RoutesRegistererFunc func(r RouteHandler)

func (fn RoutesRegistererFunc) RegisterRoutes(r RouteHandler) { fn(r) }

var _ RouteHandler = (*ServeMux)(nil)

type serveMux = http.ServeMux

// ServeMux is a http.ServeMux wrapper which implements the RouteHandler
// interface. See http.ServeMux for more information.
type ServeMux struct{ serveMux }

// NewServeMux creates a new ServeMux and is ready to be used.
func NewServeMux() *ServeMux { return new(ServeMux) }

func (mux *ServeMux) Handle(method, pattern string, handler http.Handler) {
	mux.handle(method, pattern, handler)
}

func (mux *ServeMux) HandleFunc(method, pattern string, handler http.HandlerFunc) {
	mux.handleFunc(method, pattern, handler)
}
