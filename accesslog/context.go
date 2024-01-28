// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"context"
	"github.com/go-pogo/serv/middleware"
	"net/http"
)

type handlerNameKey struct{}

// WithHandlerName adds name as value to the request's context. It should be
// used on a per route/handler basis.
func WithHandlerName(name string, next http.Handler) http.Handler {
	return middleware.WithContextValue(handlerNameKey{}, &name).Wrap(next.ServeHTTP)
}

// SetHandlerName adds name as the value for handler name to the context.
func SetHandlerName(ctx context.Context, name string) context.Context {
	if v := ctx.Value(handlerNameKey{}); v != nil {
		*v.(*string) = name
		return ctx
	}
	return setHandlerName(ctx, &name)
}

func setHandlerName(ctx context.Context, ptr *string) context.Context {
	return context.WithValue(ctx, handlerNameKey{}, ptr)
}

// HandlerName gets the handler name from the context values.
// It may be an empty string.
func HandlerName(ctx context.Context) string {
	if v := ctx.Value(handlerNameKey{}); v != nil {
		return *v.(*string)
	}
	return ""
}
