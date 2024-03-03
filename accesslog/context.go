// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"context"
	"net/http"
)

type ctxValuesKey struct{}

type HandlerSettings struct {
	HandlerName  string
	ShouldIgnore bool
}

// WithHandlerName adds name as value to the request's context. It should be
// used on a per route/handler basis.
func WithHandlerName(name string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(
			wri,
			req.WithContext(SetHandlerName(req.Context(), name)),
		)
	})
}

// SetHandlerName adds name as the value for handler name to the context.
func SetHandlerName(ctx context.Context, name string) context.Context {
	return setSettings(ctx, func(s *HandlerSettings) {
		s.HandlerName = name
	})
}

func SetShouldIgnore(ctx context.Context, ignore bool) context.Context {
	return setSettings(ctx, func(s *HandlerSettings) {
		s.ShouldIgnore = ignore
	})
}

func Settings(ctx context.Context) HandlerSettings {
	if v := ctx.Value(ctxValuesKey{}); v != nil {
		return *v.(*HandlerSettings)
	}
	return HandlerSettings{}
}

// HandlerName gets the handler name from the context values.
// It may be an empty string.
func HandlerName(ctx context.Context) string {
	return Settings(ctx).HandlerName
}

func ShouldIgnore(ctx context.Context) bool {
	return Settings(ctx).ShouldIgnore
}

func setSettings(ctx context.Context, fn func(s *HandlerSettings)) context.Context {
	ctx, val, _ := withSettings(ctx)
	fn(val)
	return ctx
}

func withSettings(ctx context.Context) (context.Context, *HandlerSettings, bool) {
	if v := ctx.Value(ctxValuesKey{}); v != nil {
		return ctx, v.(*HandlerSettings), true
	}

	v := new(HandlerSettings)
	return context.WithValue(ctx, ctxValuesKey{}, v), v, false
}
