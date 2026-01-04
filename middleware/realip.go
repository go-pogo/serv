// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Portions of this file are cloned from https://github.com/zenazn/goji.
// The original code is licensed under the MIT license and is owned by its
// author(s).

package middleware

import (
	"net"
	"net/http"
	"strings"
)

const (
	headerForwardedFor = "X-Forwarded-For"
	headerRealIP       = "X-Real-Ip"
)

// RealIP is a middleware that sets a [http.Request]'s RemoteAddr to the results
// of parsing either "X-Forwarded-For" header(s) or the "X-Real-Ip" header (in
// that order).
//
// This middleware should be inserted fairly early in the middleware stack to
// ensure that subsequent layers (e.g., request loggers) which examine the
// RemoteAddr will see the intended value.
//
// You should only use this middleware if you can trust the headers passed to
// you (in particular, the two headers this middleware uses), for example
// because you have placed a reverse proxy like HAProxy or nginx in front of the
// server. If your reverse proxies are configured to pass along arbitrary header
// values from the client, or if you use this middleware without a reverse
// proxy, malicious clients will be able to make you very sad (or, depending on
// how you're using RemoteAddr, vulnerable to an attack of some sort).
//
// This function is based on the RealIP middleware from the Goji web framework.
func RealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		if rip := realIP(req); rip != "" {
			req.RemoteAddr = rip
		}
		next.ServeHTTP(wri, req)
	})
}

func realIP(r *http.Request) string {
	for _, xff := range r.Header.Values(headerForwardedFor) {
		for _, fwd := range strings.Split(xff, ",") {
			if fwd == "" {
				continue
			}
			if ip := net.ParseIP(strings.TrimSpace(fwd)); ip.IsGlobalUnicast() {
				return ip.String()
			}
		}
	}
	if rip := r.Header.Get(headerRealIP); rip != "" {
		return net.ParseIP(rip).String()
	}
	return ""
}
