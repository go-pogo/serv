// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"crypto/tls"
	"github.com/go-pogo/errors"
	"log"
	"net"
	"net/http"
)

type Option interface {
	apply(srv *Server) error
}

type optionFunc func(srv *Server) error

func (fn optionFunc) apply(srv *Server) error { return fn(srv) }

const panicNilLogger = "serv.WithLogger: Logger should not be nil"

// WithLogger adds a [Logger] to the [Server]. It will also set the internal
// [http.Server.ErrorLog] if [Logger] l also implements [ErrorLoggerProvider].
func WithLogger(l Logger) Option {
	return optionFunc(func(srv *Server) error {
		if l == nil {
			panic(panicNilLogger)
		}

		srv.log = l
		if srv.httpServer.ErrorLog == nil {
			if el, ok := l.(ErrorLoggerProvider); ok {
				srv.httpServer.ErrorLog = el.ErrorLogger()
			}
		}
		return nil
	})
}

// WithDefaultLogger adds a [DefaultLogger] to the [Server].
func WithDefaultLogger() Option { return WithLogger(DefaultLogger(nil)) }

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

const panicNilTLSConfig = "serv.WithTLS: tls.Config should not be nil"

// WithTLS sets the provided [tls.Config] to the [Server]'s internal
// [http.Server.TLSConfig]. Any provided [TLSOption](s) will be applied to this
// [tls.Config].
func WithTLS(conf *tls.Config, opts ...TLSOption) Option {
	return optionFunc(func(srv *Server) error {
		if conf == nil {
			panic(panicNilTLSConfig)
		}

		var err error
		for _, opt := range opts {
			err = errors.Append(err, opt.ApplyTo(conf))
		}
		if err != nil {
			return err
		}

		srv.httpServer.TLSConfig = conf
		return nil
	})
}

// WithDefaultTLSConfig sets the [Server]'s internal [http.Server.TLSConfig] to
// the value of [DefaultTLSConfig]. Any provided [TLSOption](s) will be applied
// to this [tls.Config].
func WithDefaultTLSConfig(opts ...TLSOption) Option {
	return WithTLS(DefaultTLSConfig(), opts...)
}
