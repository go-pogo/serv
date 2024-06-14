// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"context"
	"net/http"
)

// ShouldIgnore returns if the [http.Request]'s context contains a
// "should ignore" value set using [SetShouldIgnore].
func ShouldIgnore(ctx context.Context) bool {
	if v := ctx.Value(ctxSettingsKey{}); v != nil {
		return v.(*ctxSettings).shouldIgnore
	}
	return false
}

// SetShouldIgnore marks the context from a [http.Request] as "should ignore".
// It returns true if the context was successfully updated. If the context
// does not already contain a reference to
func SetShouldIgnore(ctx context.Context, ignore bool) bool {
	if _, settings, exists := withSettings(ctx); exists {
		settings.shouldIgnore = ignore
		return true
	}
	return false
}

// IgnoreHandler returns a [http.Handler] that marks the [http.Request] as
// "should ignore".
func IgnoreHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		ctx, settings, existing := withSettings(req.Context())
		settings.shouldIgnore = true
		if !existing {
			req = req.WithContext(ctx)
		}

		next.ServeHTTP(wri, req)
		wri.WriteHeader(http.StatusNoContent)
	})
}
