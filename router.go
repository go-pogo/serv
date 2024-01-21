// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import "net/http"

type RouteHandler interface {
	Handle(method, pattern string, handler http.Handler)
	HandleFunc(method, pattern string, handler http.HandlerFunc)
}

type RoutesRegisterer interface {
	RegisterRoutes(router RouteHandler)
}

var _ RouteHandler = (*ServeMux)(nil)

type serveMux = http.ServeMux

type ServeMux struct{ serveMux }

func NewServeMux() *ServeMux { return new(ServeMux) }

func (r *ServeMux) Handle(method, pattern string, handler http.Handler) {
	r.handle(method, pattern, handler)
}

func (r *ServeMux) HandleFunc(method, pattern string, handler http.HandlerFunc) {
	r.handleFunc(method, pattern, handler)
}
