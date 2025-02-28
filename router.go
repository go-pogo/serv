// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"net/http"
)

// RouteHandler handles routes.
type RouteHandler interface {
	HandleRoute(route Route)
}

// RoutesRegisterer registers routes to a [RouteHandler].
type RoutesRegisterer interface {
	RegisterRoutes(rh RouteHandler)
}

// RoutesRegistererFunc registers routes to a [RouteHandler].
type RoutesRegistererFunc func(rh RouteHandler)

func (fn RoutesRegistererFunc) RegisterRoutes(rh RouteHandler) { fn(rh) }

// RegisterRoutes registers [Route]s to a [RouteHandler].
func RegisterRoutes(rh RouteHandler, routes ...Route) {
	for _, r := range routes {
		rh.HandleRoute(r)
	}
}

var _ http.Handler = (*Route)(nil)

// Route is a [http.Handler] which represents a route that can be registered to
// a [RouteHandler].
type Route struct {
	// Name of the route.
	Name string
	// Method used to handle the route.
	Method string
	// Pattern to access the route.
	Pattern string
	// Handler is the [http.Handler] that handles the route.
	Handler http.Handler
}

func (r Route) handler() http.Handler {
	if r.Name == "" {
		return r.Handler
	}

	return AddHandlerName(r.Name, r.Handler)
}

func (r Route) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	r.handler().ServeHTTP(wri, req)
}

// Router is a [http.Handler] that can handle routes.
type Router interface {
	RouteHandler
	http.Handler
}

var (
	_ Router = (*ServeMux)(nil)
	_ Option = (*ServeMux)(nil)
)

type serveMux = http.ServeMux

// ServeMux uses an internal embedded [http.ServeMux] to handle routes. It
// implements the [Router] interface on top of that.
// See [http.ServeMux] for additional information about pattern syntax,
// compatibility etc.
type ServeMux struct {
	*serveMux
	notFound http.Handler
}

// NewServeMux creates a new [ServeMux] and is ready to be used.
func NewServeMux() *ServeMux {
	return &ServeMux{serveMux: http.NewServeMux()}
}

var defaultServeMux = ServeMux{serveMux: http.DefaultServeMux}

// DefaultServeMux returns a [ServeMux] containing [http.DefaultServeMux].
func DefaultServeMux() *ServeMux { return &defaultServeMux }

// HandleRoute registers a route to the [ServeMux] using its internal
// [http.ServeMux.Handle].
func (mux *ServeMux) HandleRoute(route Route) {
	mux.serveMux.Handle(route.Method+" "+route.Pattern, route.handler())
}

// WithNotFoundHandler sets a [http.Handler] which is called when there is no
// matching pattern. If not set, [ServeMux] will use the internal
// [http.ServeMux]'s default not found handler, which is [http.NotFound].
func (mux *ServeMux) WithNotFoundHandler(h http.Handler) *ServeMux {
	mux.notFound = h
	return mux
}

func (mux *ServeMux) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	if mux.notFound == nil {
		mux.serveMux.ServeHTTP(wri, req)
		return
	}

	// below code is taken from the http.ServeMux.ServeHTTP method
	if req.RequestURI == "*" {
		if req.ProtoAtLeast(1, 1) {
			wri.Header().Set("Connection", "close")
		}
		wri.WriteHeader(http.StatusBadRequest)
		return
	}
	if h, pattern := mux.serveMux.Handler(req); pattern != "" {
		// a known limitation is that the pattern and matches are not set in the
		// http.Request, thus http.Request.PathValue won't work as expected
		h.ServeHTTP(wri, req)
	} else {
		mux.notFound.ServeHTTP(wri, req)
	}
}

func (mux *ServeMux) apply(srv *Server) error {
	srv.Handler = mux
	return nil
}
