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
	ServerStartTLS(name, addr, certFile, keyFile string)
	ServerShutdown(name string)
	ServerClose(name string)
}

type ErrorLoggerProvider interface {
	ErrorLogger() *log.Logger
}

type ErrorLogger interface {
	Logger
	ErrorLoggerProvider
}

// DefaultLogger returns a [Logger] that uses a [log.Logger] to log the
// [Server]'s lifecycle events. It defaults to [log.Default] if the provided
// [log.Logger] l is nil.
func DefaultLogger(l *log.Logger) ErrorLogger {
	if l == nil {
		l = log.Default()
	}
	return &defaultLogger{l}
}

// NopLogger returns a [Logger] that does nothing.
func NopLogger() Logger { return new(nopLogger) }

type defaultLogger struct{ *log.Logger }

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

func (l *defaultLogger) ServerStartTLS(name, addr, certFile, keyFile string) {
	if certFile == "" || keyFile == "" {
		l.Logger.Println(l.name(name) + " starting on " + addr + " using TLS")
	} else {
		l.Logger.Printf(
			"%s starting on %s using TLS with cert file %s and key file %s\n",
			l.name(name),
			addr,
			certFile,
			keyFile,
		)
	}
}

func (l *defaultLogger) ServerShutdown(name string) {
	l.Logger.Println(l.name(name) + " shutting down")
}

func (l *defaultLogger) ServerClose(name string) {
	l.Logger.Println(l.name(name) + " closing")
}

type nopLogger struct{}

func (*nopLogger) ServerStart(_, _ string)          {}
func (*nopLogger) ServerStartTLS(_, _, _, _ string) {}
func (*nopLogger) ServerShutdown(string)            {}
func (*nopLogger) ServerClose(string)               {}
