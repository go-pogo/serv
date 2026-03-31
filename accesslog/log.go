// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-pogo/serv"
)

const Message string = "access request"

type Logger interface {
	LogAccess(ctx context.Context, det Details, req *http.Request)
}

const panicNewNilLogger = "accesslog.NewLogger: slog.Logger should not be nil"

func NewLogger(l *slog.Logger) Logger {
	if l == nil {
		panic(panicNewNilLogger)
	}
	return &logger{l}
}

func DefaultLogger() Logger { return &logger{slog.Default()} }

type logger struct{ *slog.Logger }

func (l *logger) LogAccess(ctx context.Context, det Details, req *http.Request) {
	attrs := make([]any, 0, 5)

	// keep matching attributes in sync with serv.Logger.LogError!
	if det.ServerName != "" {
		attrs = append(attrs, slog.String("server", det.ServerName))
	}
	if det.HandlerName != "" {
		attrs = append(attrs, slog.String("handler", det.HandlerName))
	}
	if det.RequestID != "" {
		attrs = append(attrs, slog.String("request_id", det.RequestID))
	}
	attrs = append(attrs,
		slog.GroupAttrs("request",
			slog.String("method", req.Method),
			slog.String("proto", req.Proto),
			slog.String("uri", serv.RequestURI(req)),
			slog.String("remote_addr", serv.RemoteAddr(req)),
		),
		slog.Int("status_code", det.StatusCode),
		slog.Int64("bytes_written", det.BytesWritten),
		slog.Duration("duration", det.Duration),
	)

	l.InfoContext(ctx, Message, attrs...)
}

func NopLogger() Logger { return new(nopLogger) }

type nopLogger struct{}

func (*nopLogger) LogAccess(context.Context, Details, *http.Request) {}

// loggerFunc acts as middleware wrapper for Logger
type loggerFunc func(ctx context.Context, det Details, req *http.Request)

func (fn loggerFunc) LogAccess(ctx context.Context, det Details, req *http.Request) {
	fn(ctx, det, req)
}
