// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"log"
)

// Logger logs a [Server]'s lifecycle events.
type Logger interface {
	LogServerStart(name, addr string)
	LogServerStartTLS(name, addr, certFile, keyFile string)
	LogServerShutdown(name string)
	LogServerClose(name string)
}

type ErrorLoggerProvider interface {
	ErrorLogger() *log.Logger
}

type ErrorLogger interface {
	Logger
	ErrorLoggerProvider
}

const panicNewNilLogger = "serv.NewLogger: log.Logger should not be nil"

// NewLogger returns an [ErrorLogger] that uses a [log.Logger] to log the
// [Server]'s lifecycle events.
func NewLogger(l *log.Logger) ErrorLogger {
	if l == nil {
		panic(panicNewNilLogger)
	}
	return &logger{l}
}

// DefaultLogger returns an [ErrorLogger] that uses [log.Default] to log the
// [Server]'s lifecycle events.
func DefaultLogger() ErrorLogger { return &logger{log.Default()} }

type logger struct{ *log.Logger }

func (l *logger) ErrorLogger() *log.Logger { return l.Logger }

func (l *logger) name(name string) string {
	if name == "" {
		return "server"
	}
	return "server " + name
}

func (l *logger) LogServerStart(name, addr string) {
	l.Logger.Println(l.name(name) + " starting on " + addr)
}

func (l *logger) LogServerStartTLS(name, addr, certFile, keyFile string) {
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

func (l *logger) LogServerShutdown(name string) {
	l.Logger.Println(l.name(name) + " shutting down")
}

func (l *logger) LogServerClose(name string) {
	l.Logger.Println(l.name(name) + " closing")
}

// NopLogger returns a [Logger] that does nothing.
func NopLogger() Logger { return new(nopLogger) }

type nopLogger struct{}

func (*nopLogger) LogServerStart(_, _ string)          {}
func (*nopLogger) LogServerStartTLS(_, _, _, _ string) {}
func (*nopLogger) LogServerShutdown(string)            {}
func (*nopLogger) LogServerClose(string)               {}
