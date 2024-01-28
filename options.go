// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"github.com/go-pogo/serv/middleware"
	"net"
	"net/http"
)

type Option interface {
	apply(s *Server) error
}

type optionFunc func(s *Server) error

func (fn optionFunc) apply(s *Server) error { return fn(s) }

// WithOptions wraps multiple Option(s) into a single Option.
func WithOptions(opts ...Option) Option {
	switch len(opts) {
	case 0:
		return nil
	case 1:
		return opts[0]
	default:
		return optionFunc(func(srv *Server) error {
			return srv.Apply(opts...)
		})
	}
}

func WithHandler(h http.Handler) Option {
	return optionFunc(func(s *Server) error {
		s.Handler = h
		return nil
	})
}

// WithMiddleware adds the middleware.Middleware to an internal list. When the
// Server is started, it's Handler is wrapped with this middleware.
func WithMiddleware(mw ...middleware.Middleware) Option {
	return optionFunc(func(s *Server) error {
		if s.middlewares == nil {
			s.middlewares = mw
		} else {
			s.middlewares = append(s.middlewares, mw...)
		}
		return nil
	})
}

type serverNameKey struct{}

// WithName adds the server's name as value to the request's context.
func WithName(name string) Option {
	return optionFunc(func(s *Server) error {
		s.name = name
		return WithMiddleware(
			middleware.WithContextValue(serverNameKey{}, name),
		).apply(s)
	})
}

// ServerName gets the name from the context values which may be an empty string.
func ServerName(ctx context.Context) string {
	if v := ctx.Value(serverNameKey{}); v != nil {
		return v.(string)
	}
	return ""
}

// BaseContext returns a function which returns the provided context.Context.
func BaseContext(ctx context.Context) func(_ net.Listener) context.Context {
	return func(_ net.Listener) context.Context { return ctx }
}

func WithBaseContext(ctx context.Context) Option {
	return optionFunc(func(s *Server) error {
		s.httpServer.BaseContext = BaseContext(ctx)
		return nil
	})
}
