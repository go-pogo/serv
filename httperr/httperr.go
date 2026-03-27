// Copyright (c) 2026, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httperr

import "net/http"

// A Handler responds to an HTTP request the same way an [http.Handler] does
// with the difference of it being able to return an error.
type Handler interface {
	ServeHTTPError(http.ResponseWriter, *http.Request) error
}

// The HandlerFunc type is an adapter to allow the use of ordinary functions as
// HTTP handlers, the same way as [http.HandlerFunc]. The difference is that
// HandlerFunc may return an error.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// ServeHTTPError calls f(w, r).
func (f HandlerFunc) ServeHTTPError(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}
