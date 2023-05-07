// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"context"
	"net/http"
)

func WithContextValue(key, value interface{}) Middleware {
	return MiddlewareFunc(func(next http.HandlerFunc) http.Handler {
		return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(
				wri,
				req.WithContext(context.WithValue(req.Context(), key, value)),
			)
		})
	})
}
