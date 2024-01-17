// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"net"
	"net/http"
	"time"
)

// Details are collected With Wrap and contain additional details of a request
// and it's corresponding response.
type Details struct {
	ServerName  string
	HandlerName string

	// StatusCode is the first http response code passed to the WriteHeader func
	// of the ResponseWriter. If no such call is made, a default code of 200 is
	// assumed instead.
	StatusCode int
	// StartTime is the time the request was received.
	StartTime time.Time
	// Duration is the time it took to execute the handler.
	Duration time.Duration
	// BytesWritten is the number of bytes successfully written by the Write or
	// ReadFrom function of the ResponseWriter. ResponseWriters may also write
	// data to their underlying connection directly (e.g. headers), but those
	// are not tracked. Therefor the number of BytesWritten bytes will usually
	// match the size of the response body.
	BytesWritten int64
	// RequestCount is the amount of open requests during the execution of the
	// handler.
	RequestCount int64
}

// RemoteAddr returns a sanitized remote address from the http.Request.
func RemoteAddr(r *http.Request) string {
	addr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return addr
}

// RequestURI
// https://www.rfc-editor.org/rfc/rfc7540#section-8.3
func RequestURI(r *http.Request) string {
	var uri string
	if r.ProtoMajor == 2 && r.Method == "CONNECT" {
		uri = r.Host
	} else {
		uri = r.RequestURI
	}
	if uri == "" {
		uri = r.URL.RequestURI()
	}
	return uri
}

func Username(r *http.Request) string {
	if r.URL != nil && r.URL.User != nil {
		if user := r.URL.User.Username(); user != "" {
			return user
		}
	}
	return ""
}
