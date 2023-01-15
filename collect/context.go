// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collect

import (
	"context"
	"net/http"

	"github.com/go-pogo/serv/middleware"
)

type handlerNameKey struct{}

// WithHandlerName adds name as value to the request's context.
func WithHandlerName(name string, next http.Handler) http.Handler {
	return middleware.WithContextValue(handlerNameKey{}, name, next)
}

// WithHandlerNameM adds name as value to the request's context.
func WithHandlerNameM(name string) middleware.Middleware {
	return middleware.WithContextValueM(handlerNameKey{}, name)
}

// HandlerName gets the handler name from the context values which may be an
// empty string.
func HandlerName(ctx context.Context) string {
	return ctx.Value(handlerNameKey{}).(string)
}

// HandlerNameOr uses HandlerName to get the handler name from the context
// values. It returns def when the result would otherwise be an empty string.
func HandlerNameOr(ctx context.Context, def string) string {
	if n := HandlerName(ctx); n != "" {
		return n
	}
	return def
}
