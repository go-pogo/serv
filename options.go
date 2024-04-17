// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"crypto/tls"
	"github.com/go-pogo/errors"
	"net"
	"net/http"
)

type Option interface {
	apply(s *Server) error
}

type optionFunc func(s *Server) error

func (fn optionFunc) apply(s *Server) error { return fn(s) }

// WithOptions wraps multiple options into a single [Option].
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

// WithHandler sets the [Server]'s [Server.Handler] to h.
func WithHandler(h http.Handler) Option {
	return optionFunc(func(s *Server) error {
		s.Handler = h
		return nil
	})
}

const ErrHandlerIsNoRouteHandler errors.Msg = "server handler is not a RouteHandler"

// WithRoutes adds the provided [RoutesRegisterer] to the [Server]'s
// [Server.Handler]. It will use [DefaultServeMux] as handler when
// [Server.Handler] is nil.
// It returns an [ErrHandlerIsNoRouteHandler] error when
// [Server.Handler] is not a [RouteHandler].
func WithRoutes(reg ...RoutesRegisterer) Option {
	return optionFunc(func(s *Server) error {
		if s.Handler == nil {
			mux := DefaultServeMux()
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

// WithName adds the [Server]'s name as value to the [http.Request]'s context
// by wrapping the [Server.Handler] with [AddServerName]. This is done when the
// [Server] starts.
func WithName(name string) Option {
	return optionFunc(func(s *Server) error {
		s.name = name
		return nil
	})
}

// BaseContext returns a function which returns the provided context.
func BaseContext(ctx context.Context) func(_ net.Listener) context.Context {
	return func(_ net.Listener) context.Context { return ctx }
}

// WithBaseContext sets the provided [context.Context] ctx to the [Server]'s
// internal [http.Server.BaseContext].
func WithBaseContext(ctx context.Context) Option {
	return optionFunc(func(s *Server) error {
		s.httpServer.BaseContext = BaseContext(ctx)
		return nil
	})
}

const panicNilTLSConfig = "serv.WithTLS: tls.Config should not be nil"

// WithTLS sets the provided [tls.Config] to the [Server]'s internal
// [http.Server.TLSConfig]. Any provided [TLSOption](s) will be applied to this
// [tls.Config].
func WithTLS(conf *tls.Config, opts ...TLSOption) Option {
	return optionFunc(func(s *Server) error {
		if conf == nil {
			s.TLSConfig = DefaultTLSConfig()
		} else {
			s.TLSConfig = conf
		}

		var err error
		for _, opt := range opts {
			err = errors.Append(err, opt.ApplyTo(conf))
		}
		return err
	})
}
