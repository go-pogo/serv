// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"context"
	"log"
	"net/http"
)

const Message string = "access request"

type Logger interface {
	LogAccess(ctx context.Context, det Details, req *http.Request)
}

const panicNewNilLogger = "accesslog.NewLogger: log.Logger should not be nil"

func NewLogger(l *log.Logger) Logger {
	if l == nil {
		panic(panicNewNilLogger)
	}
	return &logger{l}
}

func DefaultLogger() Logger { return &logger{log.Default()} }

type logger struct{ *log.Logger }

func (l *logger) LogAccess(_ context.Context, det Details, req *http.Request) {
	handlerName := det.HandlerName
	if handlerName == "" {
		handlerName = "-"
	}

	l.Logger.Printf("%s: %s %s \"%s %s %s\" %d %db %s\n",
		Message,
		RemoteAddr(req),
		handlerName,
		req.Method,
		RequestURI(req),
		req.Proto,
		det.StatusCode,
		det.BytesWritten,
		det.Duration,
	)
}

func NopLogger() Logger { return new(nopLogger) }

type nopLogger struct{}

func (*nopLogger) LogAccess(context.Context, Details, *http.Request) {}

// loggerFunc acts as middleware wrapper for Logger
type loggerFunc func(ctx context.Context, det Details, req *http.Request)

func (fn loggerFunc) LogAccess(ctx context.Context, det Details, req *http.Request) {
	fn(ctx, det, req)
}
