// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"net/http"
	"strings"
)

// Middleware wraps a http.Handler with additional logic and returns a new
// http.Handler that executes it.
type Middleware interface {
	Wrap(next http.Handler) http.Handler
}

type MiddlewareFunc func(next http.Handler) http.Handler

func (fn MiddlewareFunc) Wrap(next http.Handler) http.Handler { return fn(next) }

// Wrap http.Handler h with the provided Middleware.
func Wrap(h http.Handler, mw ...Middleware) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i].Wrap(h)
	}
	return h
}

func WithContextValue(key, value interface{}, next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(
			wri,
			req.WithContext(context.WithValue(req.Context(), key, value)),
		)
	})
}

func WithContextValueM(key, value interface{}) Middleware {
	return MiddlewareFunc(func(next http.Handler) http.Handler {
		return WithContextValue(key, value, next)
	})
}

// RemoveTrailingSlash adds middleware that redirects any request which has a
// trailing slash to the equivalent path without trailing slash.
// It uses the redirect code http.StatusMovedPermanently for any GET request
// and http.StatusTemporaryRedirect for any other method.
func RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" || strings.HasSuffix(req.URL.Path, "/") {
			http.Redirect(wri, req, strings.TrimRight(req.URL.Path, "/"), statusCode(req.Method))
		} else {
			next.ServeHTTP(wri, req)
		}
	})
}

func RemoveTrailingSlashM() Middleware { return MiddlewareFunc(RemoveTrailingSlash) }

// RedirectHttps adds middleware that redirects any http request to its https
// equivalent url.
// It uses the redirect code http.StatusMovedPermanently for any GET request
// and http.StatusTemporaryRedirect for any other method.
func RedirectHttps(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		if req.URL.Scheme == "http" {
			req.URL.Scheme += "s"
			http.Redirect(wri, req, req.URL.String(), statusCode(req.Method))
		} else {
			next.ServeHTTP(wri, req)
		}
	})
}

func RedirectHttpsM() Middleware { return MiddlewareFunc(RedirectHttps) }

func statusCode(method string) int {
	switch method {
	case http.MethodGet:
		return http.StatusMovedPermanently
	default:
		return http.StatusTemporaryRedirect
	}
}
