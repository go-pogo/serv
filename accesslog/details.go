// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"net/http"
	"time"

	"github.com/go-pogo/serv"
)

// Details are collected using [NewHandler] and contain additional details of a
// request and it's corresponding response.
type Details struct {
	serv.Info

	UserAgent string

	// StatusCode is the first http response code passed to the
	// [http.ResponseWriter.WriteHeader]. See [httpsnoop.Metrics] for additional
	// information.
	StatusCode int
	// StartTime is the time the request was received.
	StartTime time.Time
	// Duration is the time it took to execute the handler.
	Duration time.Duration
	// BytesWritten is the number of bytes successfully written by the
	// [http.ResponseWriter.Write] or [http.ResponseWriter.ReadFrom] functions.
	// See [httpsnoop.Metrics] for additional information.
	BytesWritten int64
	// RequestCount is the amount of open requests during the execution of the
	// handler.
	RequestCount int64
}

// RemoteAddr returns a sanitized remote address from the [http.Request].
// Add [middleware.RealIP] middleware to your [http.Handler] to handle (proxy)
// forwarded traffic.
//
// Deprecated: use serv.RemoteAddr instead.
func RemoteAddr(r *http.Request) string { return serv.RemoteAddr(r) }

// RequestURI
// https://www.rfc-editor.org/rfc/rfc7540#section-8.3
//
// Deprecated: use serv.RequestURI instead.
func RequestURI(r *http.Request) string { return serv.RequestURI(r) }

// Username returns a username when available in the request's url.
//
// Deprecated: will be removed in future releases.
func Username(r *http.Request) string {
	if r.URL != nil {
		return r.URL.User.Username()
	}
	return ""
}
