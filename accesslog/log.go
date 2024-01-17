// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"context"
	"log"
	"net/http"
)

const Prefix string = "access request"

type Logger interface {
	Log(ctx context.Context, det Details, req *http.Request)
}

func DefaultLogger(l *log.Logger) Logger {
	if l == nil {
		l = log.Default()
	}
	return &defaultLogger{l}
}

var _ Logger = (*defaultLogger)(nil)

type defaultLogger struct {
	*log.Logger
}

func (l *defaultLogger) Log(_ context.Context, det Details, req *http.Request) {
	handlerName := det.HandlerName
	if handlerName == "" {
		handlerName = "-"
	}

	l.Printf("%s: %s %s \"%s %s %s\" %d %db\n",
		Prefix,
		RemoteAddr(req),
		handlerName,
		req.Method,
		RequestURI(req),
		req.Proto,
		det.StatusCode,
		det.BytesWritten,
	)
}

type nopLogger struct{}

func NopLogger() Logger { return new(nopLogger) }

func (*nopLogger) Log(context.Context, Details, *http.Request) {}

// loggerFunc acts as middleware wrapper for Logger
type loggerFunc func(ctx context.Context, det Details, req *http.Request)

func (fn loggerFunc) Log(ctx context.Context, det Details, req *http.Request) { fn(ctx, det, req) }
