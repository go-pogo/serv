// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"context"
)

type ctxSettingsKey struct{}

type ctxSettings struct {
	shouldIgnore bool
}

func ShouldIgnore(ctx context.Context) bool {
	if v := ctx.Value(ctxSettingsKey{}); v != nil {
		return v.(*ctxSettings).shouldIgnore
	}
	return false
}

func SetShouldIgnore(ctx context.Context, ignore bool) context.Context {
	ctx, settings, _ := withSettings(ctx)
	settings.shouldIgnore = ignore
	return ctx
}

func withSettings(ctx context.Context) (context.Context, *ctxSettings, bool) {
	if v := ctx.Value(ctxSettingsKey{}); v != nil {
		return ctx, v.(*ctxSettings), true
	}

	v := new(ctxSettings)
	return context.WithValue(ctx, ctxSettingsKey{}, v), v, false
}
