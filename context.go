// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"net/http"
)

type ctxValuesKey struct{}

type ctxValues struct {
	serverName  string
	handlerName string
}

// ServerName gets the server's name from context values. Its return value may
// be an empty string.
func ServerName(ctx context.Context) string {
	if v := ctx.Value(ctxValuesKey{}); v != nil {
		return v.(*ctxValues).serverName
	}
	return ""
}

// AddServerName adds the server's name to the request's context. This is done
// automatically when a name is set using [WithName].
// The server's name can be retrieved using [ServerName].
func AddServerName(name string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		ctx, settings, exists := withCtxValues(req.Context())
		settings.serverName = name

		if !exists {
			// add new context to request
			req = req.WithContext(ctx)
		}
		next.ServeHTTP(wri, req)
	})
}

// HandlerName gets the handler's name from the context values. Its returned
// value may be an empty string.
func HandlerName(ctx context.Context) string {
	if v := ctx.Value(ctxValuesKey{}); v != nil {
		return v.(*ctxValues).handlerName
	}
	return ""
}

// AddHandlerName adds name as value to the request's context. It should
// be used on a per route/handler basis.
// The handler's name can be retrieved using [HandlerName].
func AddHandlerName(name string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		ctx, settings, exists := withCtxValues(req.Context())
		settings.handlerName = name

		if !exists {
			// add new context to request
			req = req.WithContext(ctx)
		}
		next.ServeHTTP(wri, req)
	})
}

func withCtxValues(ctx context.Context) (context.Context, *ctxValues, bool) {
	if v := ctx.Value(ctxValuesKey{}); v != nil {
		return ctx, v.(*ctxValues), true
	}

	v := new(ctxValues)
	return context.WithValue(ctx, ctxValuesKey{}, v), v, false
}
