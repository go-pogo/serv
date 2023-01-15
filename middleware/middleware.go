// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"context"
	"net/http"
	"strings"
)

// Middleware wraps a http.HandlerFunc with additional logic and returns a new
// http.Handler.
type Middleware interface {
	Wrap(next http.HandlerFunc) http.Handler
}

var (
	_ Middleware = new(MiddlewareFunc)
	_ Middleware = new(HandlerFunc)
)

type MiddlewareFunc func(next http.HandlerFunc) http.Handler

func (fn MiddlewareFunc) Wrap(next http.HandlerFunc) http.Handler { return fn(next) }

type HandlerFunc http.HandlerFunc

func (fn HandlerFunc) Wrap(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		fn(wri, req)
		next(wri, req)
	})
}

// Wrap http.Handler h with the provided Middleware.
func Wrap(handler http.Handler, wrap ...Middleware) http.Handler {
	for i := len(wrap) - 1; i >= 0; i-- {
		handler = wrap[i].Wrap(handler.ServeHTTP)
	}
	return handler
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
	return MiddlewareFunc(func(next http.HandlerFunc) http.Handler {
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

func RemoveTrailingSlashM() Middleware {
	return MiddlewareFunc(func(next http.HandlerFunc) http.Handler {
		return RemoveTrailingSlash(next)
	})
}

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

func RedirectHttpsM() Middleware {
	return MiddlewareFunc(func(next http.HandlerFunc) http.Handler {
		return RedirectHttps(next)
	})
}

func statusCode(method string) int {
	switch method {
	case http.MethodGet:
		return http.StatusMovedPermanently
	default:
		return http.StatusTemporaryRedirect
	}
}
