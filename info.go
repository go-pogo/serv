// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"net/http"
)

type ctxInfoKey struct{}

type Info struct {
	ServerName  string
	HandlerName string
	RequestID   string
}

// ContextWithInfo adds an Info value to the context. It returns a derived
// context that points to the parent context.Context when Info is not already
// added. Otherwise, it will update the previously added Info with info and
// return the context as is.
func ContextWithInfo(ctx context.Context, info Info) context.Context {
	if v := InfoFromContext(ctx); v != nil {
		*v = info
		return ctx
	}
	return context.WithValue(ctx, ctxInfoKey{}, &info)
}

// InfoFromContext returns the Info value from the context values, or nil.
func InfoFromContext(ctx context.Context) *Info {
	if v := ctx.Value(ctxInfoKey{}); v != nil {
		return v.(*Info)
	}
	return nil
}

// ServerName gets the server's name from the context values. Its returned
// value may be an empty string.
func ServerName(ctx context.Context) string {
	if info := InfoFromContext(ctx); info != nil {
		return info.ServerName
	}
	return ""
}

// HandlerName gets the handler's name from the context values. Its returned
// value may be an empty string.
func HandlerName(ctx context.Context) string {
	if info := InfoFromContext(ctx); info != nil {
		return info.HandlerName
	}
	return ""
}

// RequestID gets the request id from the context values. Its returned value
// may be an empty string.
func RequestID(ctx context.Context) string {
	if info := InfoFromContext(ctx); info != nil {
		return info.RequestID
	}
	return ""
}

// AddServerName sets the request context value [Info.ServerName] field with
// the value name. It should be used on a per-server basis and is done
// automatically when a [Server]'s name is set using [WithName].
func AddServerName(name string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		req, info := requestWithInfo(req)
		info.ServerName = name
		next.ServeHTTP(wri, req)
	})
}

// AddHandlerName sets the request context value [Info.HandlerName] field with
// the value name. It should be used on a per route/handler basis.
func AddHandlerName(name string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		req, info := requestWithInfo(req)
		info.HandlerName = name
		next.ServeHTTP(wri, req)
	})
}

// AddRequestID sets the request context value [Info.RequestID] field with the
// value id. It should be used on a per-request basis.
func AddRequestID(id string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		req, info := requestWithInfo(req)
		info.RequestID = id
		next.ServeHTTP(wri, req)
	})
}

func requestWithInfo(req *http.Request) (*http.Request, *Info) {
	if v := InfoFromContext(req.Context()); v != nil {
		return req, v
	}

	var info Info
	return req.WithContext(
		context.WithValue(req.Context(), ctxInfoKey{}, &info),
	), &info
}
