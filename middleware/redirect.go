// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
	"strings"
)

// RedirectHTTPS adds middleware that redirects any http request to its https
// equivalent url, using [http.Redirect].
// It uses the redirect code [http.StatusMovedPermanently] for any GET request
// and [http.StatusTemporaryRedirect] for any other method.
func RedirectHTTPS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		if req.URL.Scheme == "http" {
			req.URL.Scheme += "s"
			redirect(wri, req)
		} else {
			next.ServeHTTP(wri, req)
		}
	})
}

// RemoveTrailingSlash adds middleware that redirects any request which has a
// trailing slash to the equivalent path without trailing slash, using
// [http.Redirect].
// It uses the redirect code [http.StatusMovedPermanently] for any GET request
// and [http.StatusTemporaryRedirect] for any other method.
func RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" && strings.HasSuffix(req.URL.Path, "/") {
			req.URL.Path = strings.TrimRight(req.URL.Path, "/")
			redirect(wri, req)
		} else {
			next.ServeHTTP(wri, req)
		}
	})
}

func redirect(wri http.ResponseWriter, req *http.Request) {
	code := http.StatusTemporaryRedirect
	if req.Method == http.MethodGet {
		code = http.StatusMovedPermanently
	}

	http.Redirect(wri, req, req.URL.String(), code)
}
