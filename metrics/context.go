// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metrics

import (
	"context"
	"net/http"

	"github.com/go-pogo/serv"
)

type handlerNameKey struct{}

// WithHandlerName adds name as value to the request's context.
func WithHandlerName(name string, next http.Handler) http.Handler {
	return serv.WithContextValue(handlerNameKey{}, name, next)
}

// WithHandlerNameM adds name as value to the request's context.
func WithHandlerNameM(name string) serv.Middleware {
	return serv.WithContextValueM(handlerNameKey{}, name)
}

// HandlerName gets the handler name from the context values. It returns an
// empty string and false when the handler name is not set.
func HandlerName(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(handlerNameKey{}).(string)
	return v, ok
}
