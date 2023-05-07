// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"github.com/go-pogo/serv/middleware"
	"net/http"
)

type serverNameKey struct{}

// WithName adds the server's name as value to the request's context.
func WithName(name string) Option {
	return optionFunc(func(s *Server) error {
		s.name = name
		WithMiddleware(
			middleware.WithContextValue(serverNameKey{}, name),
		).applyTo(s)
		return nil
	})
}

// ServerName gets the name from the context values which may be an empty string.
func ServerName(ctx context.Context) string {
	if v := ctx.Value(serverNameKey{}); v != nil {
		return v.(string)
	}
	return ""
}

type handlerNameKey struct{}

// WithHandlerName adds name as value to the request's context. It should be
// used on a per route/handler basis.
func WithHandlerName(name string, next http.Handler) http.Handler {
	return middleware.WithContextValue(handlerNameKey{}, name).Wrap(next.ServeHTTP)
}

// HandlerName gets the handler name from the context values which may be an
// empty string.
func HandlerName(ctx context.Context) string {
	if v := ctx.Value(handlerNameKey{}); v != nil {
		return v.(string)
	}
	return ""
}
