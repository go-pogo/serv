// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Portions of this file are cloned from https://github.com/zenazn/goji.
// The original code is licensed under the MIT license and is owned by its
// author(s).

package response

import (
	"net/http"
	"time"
)

// Unix epoch time
var epoch = time.Unix(0, 0).Format(time.RFC1123)

// Taken from https://github.com/mytrile/nocache
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

// NoCache is a simple piece of middleware that sets a number of HTTP headers
// to prevent a router from being cached by an upstream proxy and/or client.
//
// As per http://wiki.nginx.org/HttpProxyModule - NoCache sets:
//
//	Expires: Thu, 01 Jan 1970 00:00:00 UTC
//	Cache-Control: no-cache, private, max-age=0
//	X-Accel-Expires: 0
//	Pragma: no-cache (for HTTP/1.0 proxies/clients)
//
// This function is based on the NoCache middleware from the Goji web framework.
func NoCache(wri http.ResponseWriter, req *http.Request) {
	// Delete any ETag headers that may have been set
	for _, v := range etagHeaders {
		if req.Header.Get(v) != "" {
			req.Header.Del(v)
		}
	}

	// Set our WithNoCache headers
	for k, v := range noCacheHeaders {
		wri.Header().Set(k, v)
	}
}
