// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"log"
	"log/slog"
	"net/http"
)

// ServerLogger logs a [Server]'s lifecycle events.
type ServerLogger interface {
	LogServerStart(name, addr string)
	LogServerStartTLS(name, addr, certFile, keyFile string)
	LogServerShutdown(name string)
	LogServerClose(name string)
}

// NopLogger returns a [ServerLogger] that does nothing.
func NopLogger() ServerLogger { return new(nopLogger) }

type nopLogger struct{}

func (*nopLogger) LogServerStart(_, _ string)          {}
func (*nopLogger) LogServerStartTLS(_, _, _, _ string) {}
func (*nopLogger) LogServerShutdown(string)            {}
func (*nopLogger) LogServerClose(string)               {}

type ErrorLoggerProvider interface {
	ErrorLogger() *log.Logger
}

var (
	_ ServerLogger        = (*Logger)(nil)
	_ ErrorLoggerProvider = (*Logger)(nil)
)

type Logger struct {
	*slog.Logger
	errLogger *log.Logger
}

const panicNewNilLogger = "serv.NewLogger: slog.Logger should not be nil"

// NewLogger returns an [ErrorLogger] that uses a [slog.Logger] to log the
// [Server]'s lifecycle events.
func NewLogger(l *slog.Logger) *Logger {
	if l == nil {
		panic(panicNewNilLogger)
	}
	return newLogger(l)
}

// DefaultLogger returns an [ErrorLogger] that uses [slog.Default] to log the
// [Server]'s lifecycle events.
func DefaultLogger() *Logger { return newLogger(slog.Default()) }

func newLogger(l *slog.Logger) *Logger {
	return &Logger{
		Logger:    l,
		errLogger: slogErrorLogger(l),
	}
}

func (l *Logger) ErrorLogger() *log.Logger { return l.errLogger }

func (l *Logger) LogServerStart(name, addr string) {
	l.Info("server starting",
		slog.String("name", name),
		slog.String("addr", addr),
	)
}

func (l *Logger) LogServerStartTLS(name, addr, certFile, keyFile string) {
	l.Info("server starting using TLS",
		slog.String("name", name),
		slog.String("addr", addr),
		slog.String("cert_file", certFile),
		slog.String("key_file", keyFile),
	)
}

func (l *Logger) LogServerShutdown(name string) {
	l.Info("server shutting down", slog.String("name", name))
}

func (l *Logger) LogServerClose(name string) {
	l.Info("server closing", slog.String("name", name))
}

// LogError logs any non-nil error together with the provided [http.Request]'s
// Method, Proto, RequestURI and RemoteAddr fields as attributes.
// Additional attributes that are logged are the server's name, handler's name
// and id of the request, when any of these are available from the provided
// context.
func (l *Logger) LogError(req *http.Request, err error) {
	if err == nil {
		return
	}

	ctx := req.Context()
	if !l.Handler().Enabled(ctx, slog.LevelError) {
		return
	}

	attrs := make([]any, 0, 5)
	attrs = append(attrs, slog.Any("err", err))

	// keep matching attributes in sync with accesslog.logger.LogAccess!
	if info := InfoFromContext(ctx); info != nil {
		if info.ServerName != "" {
			attrs = append(attrs, slog.String("server", info.ServerName))
		}
		if info.HandlerName != "" {
			attrs = append(attrs, slog.String("handler", info.HandlerName))
		}
		if info.RequestID != "" {
			attrs = append(attrs, slog.String("request_id", info.RequestID))
		}
	}
	attrs = append(attrs, slog.GroupAttrs("request",
		slog.String("method", req.Method),
		slog.String("proto", req.Proto),
		slog.String("uri", RequestURI(req)),
		slog.String("remote_addr", RemoteAddr(req)),
	))

	l.ErrorContext(ctx, "request handler error", attrs...)
}
