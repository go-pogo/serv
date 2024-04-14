// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clflogger

import (
	"context"
	"fmt"
	"github.com/go-pogo/serv/accesslog"
	"io"
	"net/http"
	"os"
)

const TimeLayout string = "02/Jan/2006:15:04:05 -0700"

// https://httpd.apache.org/docs/current/logs.html#common
type CommonLogger struct{ io.Writer }

func (l *CommonLogger) Log(_ context.Context, det accesslog.Details, req *http.Request) {
	if l.Writer == nil {
		l.Writer = os.Stdout
	}

	username := accesslog.Username(req)
	if username == "" {
		username = "-"
	}

	_, _ = fmt.Fprintf(l, "%s - %s [%s] \"%s %s %s\" %d %d\n",
		accesslog.RemoteAddr(req),
		username,
		det.StartTime.Format(TimeLayout),
		req.Method,
		accesslog.RequestURI(req),
		req.Proto,
		det.StatusCode,
		det.BytesWritten,
	)
}
