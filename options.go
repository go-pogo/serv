// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"github.com/go-pogo/easytls"
	"github.com/go-pogo/errors"
)

type Option interface {
	apply(srv *Server) error
}

type optionFunc func(srv *Server) error

func (fn optionFunc) apply(srv *Server) error { return fn(srv) }

const panicNilLogger = "serv.WithLogger: Logger should not be nil"

// WithLogger adds a [Logger] to the [Server]. It will also set the internal
// [http.Server.ErrorLog] if [Logger] l also implements [ErrorLoggerProvider].
func WithLogger(log Logger) Option {
	return optionFunc(func(srv *Server) error {
		if log == nil {
			panic(panicNilLogger)
		}

		srv.log = log
		if srv.httpServer.ErrorLog == nil {
			if el, ok := log.(ErrorLoggerProvider); ok {
				srv.httpServer.ErrorLog = el.ErrorLogger()
			}
		}
		return nil
	})
}

// WithDefaultLogger adds a default [Logger] to the [Server] using
// [DefaultLogger] and [WithLogger].
func WithDefaultLogger() Option { return WithLogger(DefaultLogger()) }

const panicNilErrorLogger = "serv.WithErrorLogger: log.Logger should not be nil"

func WithErrorLogger(l *log.Logger) Option {
	return optionFunc(func(srv *Server) error {
		if l == nil {
			panic(panicNilErrorLogger)
		}

		srv.httpServer.ErrorLog = l
		return nil
	})
}

// WithName adds the [Server]'s name as value to the [http.Request]'s context
// by wrapping the [Server.Handler] with [AddServerName]. This is done when the
// [Server] starts.
func WithName(name string) Option {
	return optionFunc(func(srv *Server) error {
		srv.name = name
		return nil
	})
}

// WithHandler sets the [Server]'s [Server.Handler] to h.
func WithHandler(h http.Handler) Option {
	return optionFunc(func(s *Server) error {
		s.Handler = h
		return nil
	})
}

const ErrHandlerIsNoRouteHandler errors.Msg = "server handler is not a RouteHandler"

// WithRoutesRegisterer uses the provided [RoutesRegisterer](s) to add [Route]s
// to the [Server]'s [Server.Handler]. It will use [DefaultServeMux] as handler
// when [Server.Handler] is nil.
// It returns an [ErrHandlerIsNoRouteHandler] error when
// [Server.Handler] is not a [RouteHandler].
func WithRoutesRegisterer(reg ...RoutesRegisterer) Option {
	return optionFunc(func(srv *Server) error {
		if srv.Handler == nil {
			mux := DefaultServeMux()
			for _, rr := range reg {
				rr.RegisterRoutes(mux)
			}
			srv.Handler = mux
			return nil
		}
		if r, ok := srv.Handler.(RouteHandler); ok {
			for _, rr := range reg {
				rr.RegisterRoutes(r)
			}
		}
		return errors.New(ErrHandlerIsNoRouteHandler)
	})
}

// BaseContext returns a function which returns the provided context.
func BaseContext(ctx context.Context) func(_ net.Listener) context.Context {
	return func(_ net.Listener) context.Context { return ctx }
}

// WithBaseContext sets the provided [context.Context] ctx to the [Server]'s
// internal [http.Server.BaseContext].
func WithBaseContext(ctx context.Context) Option {
	return optionFunc(func(srv *Server) error {
		srv.httpServer.BaseContext = BaseContext(ctx)
		return nil
	})
}

const panicNilTLSConfig = "serv.WithTLSConfig: tls.Config should not be nil"

// WithTLSConfig sets the provided [tls.Config] to the [Server]'s internal
// [http.Server.TLSConfig]. Any provided [easytls.Option](s) will be applied to
// this [tls.Config].
func WithTLSConfig(conf *tls.Config, opts ...easytls.Option) Option {
	return optionFunc(func(srv *Server) error {
		if conf == nil {
			panic(panicNilTLSConfig)
		}

		if err := easytls.Apply(conf, easytls.TargetServer, opts...); err != nil {
			return err
		}

		srv.httpServer.TLSConfig = conf
		return nil
	})
}

// WithDefaultTLSConfig sets the [Server]'s internal [http.Server.TLSConfig] to
// the value of [easytls.DefaultTLSConfig]. Any provided [easytls.Option](s)
// will be applied to this [tls.Config].
func WithDefaultTLSConfig(opts ...easytls.Option) Option {
	return WithTLSConfig(easytls.DefaultTLSConfig(), opts...)
}
