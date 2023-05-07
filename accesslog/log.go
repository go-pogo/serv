// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

const TimeLayout string = "02/Jan/2006:15:04:05 -0700"

type Logger interface {
	Log(ctx context.Context, det Details, req *http.Request)
}

type loggerFunc func(ctx context.Context, det Details, req *http.Request)

func (fn loggerFunc) Log(ctx context.Context, det Details, req *http.Request) { fn(ctx, det, req) }

var _ Logger = new(DefaultLogger)

type DefaultLogger struct{ io.Writer }

// https://httpd.apache.org/docs/current/logs.html#common
type ApacheLogger struct{ io.Writer }

func (l *DefaultLogger) Log(_ context.Context, det Details, req *http.Request) {
	if l.Writer == nil {
		l.Writer = os.Stdout
	}

	_, _ = fmt.Fprintf(l, "%s %s \"%s %s %s\" %d %db\n",
		RemoteAddr(req),
		or(det.HandlerName, "-"),
		req.Method,
		RequestURI(req),
		req.Proto,
		det.StatusCode,
		det.BytesWritten,
	)
}

func (l *ApacheLogger) Log(_ context.Context, det Details, req *http.Request) {
	if l.Writer == nil {
		l.Writer = os.Stdout
	}

	_, _ = fmt.Fprintf(l, "%s - %s [%s] \"%s %s %s\" %d %d\n",
		RemoteAddr(req),
		or(Username(req), "-"),
		det.StartTime.Format(TimeLayout),
		req.Method,
		RequestURI(req),
		req.Proto,
		det.StatusCode,
		det.BytesWritten,
	)
}

// RemoteAddr returns a sanitize                                                                d remote address from
// the/http.Request.
func RemoteAddr(r *http.Request) string {
	addr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return addr
}

// RequestURI
// https://www.rfc-editor.org/rfc/rfc7540#section-8.3
func RequestURI(r *http.Request) string {
	var uri string
	if r.ProtoMajor == 2 && r.Method == "CONNECT" {
		uri = r.Host
	} else {
		uri = r.RequestURI
	}
	if uri == "" {
		uri = r.URL.RequestURI()
	}
	return uri
}

func Username(r *http.Request) string {
	if r.URL != nil && r.URL.User != nil {
		if user := r.URL.User.Username(); user != "" {
			return user
		}
	}
	return ""
}

func or(a, b string) string {
	if a == "" {
		return b
	}
	return a
}

type nopLogger struct{}

func NopLogger() Logger { return new(nopLogger) }

func (*nopLogger) Log(context.Context, Details, *http.Request) {}
