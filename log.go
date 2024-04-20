// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"log"
)

// Logger logs a [Server]'s lifecycle events.
type Logger interface {
	ServerStart(name, addr string)
	ServerShutdown(name string)
	ServerClose(name string)
}

type ErrorLoggerProvider interface {
	ErrorLogger() *log.Logger
}

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

// DefaultLogger returns a [Logger] that uses a [log.Logger] to log the
// [Server]'s lifecycle events. It defaults to [log.Default] if the provided
// [log.Logger] l is nil.
func DefaultLogger(l *log.Logger) Logger {
	if l == nil {
		l = log.Default()
	}
	return &defaultLogger{l}
}

// NopLogger returns a [Logger] that does nothing.
func NopLogger() Logger { return new(nopLogger) }

var (
	_ Logger              = (*defaultLogger)(nil)
	_ ErrorLoggerProvider = (*defaultLogger)(nil)
)

type defaultLogger struct {
	*log.Logger
}

func (l *defaultLogger) ErrorLogger() *log.Logger { return l.Logger }

func (l *defaultLogger) name(name string) string {
	if name == "" {
		return "server"
	}
	return "server " + name
}

func (l *defaultLogger) ServerStart(name, addr string) {
	l.Logger.Println(l.name(name) + " starting on " + addr)
}

func (l *defaultLogger) ServerShutdown(name string) {
	l.Logger.Println(l.name(name) + " shutting down")
}

func (l *defaultLogger) ServerClose(name string) {
	l.Logger.Println(l.name(name) + " closing")
}

type nopLogger struct{}

func (*nopLogger) ServerStart(_, _ string) {}
func (*nopLogger) ServerShutdown(string)   {}
func (*nopLogger) ServerClose(string)      {}
