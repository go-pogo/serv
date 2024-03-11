// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"crypto/tls"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/serv/internal"
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
			return srv.With(opts...)
		})
	}
}

func WithHandler(h http.Handler) Option {
	return optionFunc(func(s *Server) error {
		s.Handler = h
		return nil
	})
}

const ErrHandlerIsNoRouteHandler errors.Msg = "server handler is not a RouteHandler"

func WithRoutes(reg ...RoutesRegisterer) Option {
	return optionFunc(func(s *Server) error {
		if s.Handler == nil {
			mux := NewServeMux()
			for _, rr := range reg {
				rr.RegisterRoutes(mux)
			}
			s.Handler = mux
			return nil
		}
		if r, ok := s.Handler.(RouteHandler); ok {
			for _, rr := range reg {
				rr.RegisterRoutes(r)
			}
		}
		return errors.New(ErrHandlerIsNoRouteHandler)
	})
}

// WithMiddleware adds the middleware.Middleware to an internal list. When the
// Server is started, it's Handler is wrapped with this middleware.
func WithMiddleware(mw ...middleware.Wrapper) Option {
	return optionFunc(func(s *Server) error {
		if s.middleware == nil {
			s.middleware = mw
		} else {
			s.middleware = append(s.middleware, mw...)
		}
		return nil
	})
}

// WithName adds the server's name as value to the request's context.
func WithName(name string) Option {
	return optionFunc(func(s *Server) error {
		s.name = name
		return WithMiddleware(internal.ServerNameMiddleware(name)).apply(s)
	})
}

// ServerName gets the server's name from context values. Its return value may
// be an empty string.
func ServerName(ctx context.Context) string { return internal.ServerName(ctx) }

// BaseContext returns a function which returns the provided context.
func BaseContext(ctx context.Context) func(_ net.Listener) context.Context {
	return func(_ net.Listener) context.Context { return ctx }
}

func WithBaseContext(ctx context.Context) Option {
	return optionFunc(func(s *Server) error {
		s.httpServer.BaseContext = BaseContext(ctx)
		return nil
	})
}

type TLSOption interface {
	Apply(conf *tls.Config) error
}

func WithTLS(conf *tls.Config, opts ...TLSOption) Option {
	return optionFunc(func(s *Server) error {
		if conf == nil {
			s.TLSConfig = DefaultTLSConfig()
		} else {
			s.TLSConfig = conf
		}

		var err error
		for _, opt := range opts {
			err = errors.Append(err, opt.Apply(conf))
		}
		return err
	})
}
