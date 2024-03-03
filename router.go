// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"github.com/go-pogo/serv/accesslog"
	"net/http"
)

// RouteHandler handles routes.
type RouteHandler interface {
	HandleRoute(route Route)
}

// RoutesRegisterer registers routes to a RouteHandler.
type RoutesRegisterer interface {
	RegisterRoutes(r RouteHandler)
}

// RoutesRegistererFunc registers routes to a RouteHandler.
type RoutesRegistererFunc func(r RouteHandler)

func (fn RoutesRegistererFunc) RegisterRoutes(r RouteHandler) { fn(r) }

var _ http.Handler = (*Route)(nil)

type Route struct {
	// Name of the route.
	Name string
	// Method used to handle the route.
	Method string
	// Pattern to access the route.
	Pattern string
	// Handler is the http.Handler that handles the route.
	Handler http.Handler
}

func (r Route) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	if r.Name == "" {
		r.Handler.ServeHTTP(wri, req)
		return
	}

	accesslog.WithHandlerName(r.Name, r.Handler).ServeHTTP(wri, req)
}

// Router is a http.Handler that can handle routes.
type Router interface {
	RouteHandler
	http.Handler
}

var (
	_ Router = (*ServeMux)(nil)
	_ Option = (*ServeMux)(nil)
)

type serveMux = http.ServeMux

// ServeMux is a http.ServeMux wrapper which implements the Router interface.
// See http.ServeMux for more information.
type ServeMux struct{ serveMux }

// NewServeMux creates a new ServeMux and is ready to be used.
func NewServeMux() *ServeMux { return new(ServeMux) }

func (mux *ServeMux) HandleRoute(route Route) { mux.handle(route) }

func (mux *ServeMux) apply(s *Server) error {
	s.Handler = mux
	return nil
}
