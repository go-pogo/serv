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

func (v *ctxValues) set(set ctxValues) {
	if v.serverName == "" {
		v.serverName = set.serverName
	}
	if v.handlerName == "" {
		v.handlerName = set.handlerName
	}
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
		next.ServeHTTP(wri, withCtxValues(req, ctxValues{
			serverName: name,
		}))
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
		next.ServeHTTP(wri, withCtxValues(req, ctxValues{
			handlerName: name,
		}))
	})
}

func withCtxValues(req *http.Request, set ctxValues) *http.Request {
	ctx := req.Context()
	if v := ctx.Value(ctxValuesKey{}); v != nil {
		v.(*ctxValues).set(set)
		return req
	}

	return req.WithContext(context.WithValue(ctx, ctxValuesKey{}, &set))
}
