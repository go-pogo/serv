// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"context"
	"log"
	"net/http"

	"github.com/go-pogo/serv/collect"
)

func Collector() collect.Collector {
	return collect.CollectorFunc(func(ctx context.Context, met collect.Metrics, req *http.Request) {
		log.Printf("%s %s \"%s %s %s\" %d %db\n",
			RemoteAddr(req),
			collect.HandlerNameOr(ctx, "-"),
			req.Method,
			RequestURI(req),
			req.Proto,
			met.Code,
			met.Written,
		)
	})
}

type Logger interface {
	Log(ctx context.Context, entry Entry)
}

type LoggerFunc func(ctx context.Context, entry Entry)

func (fn LoggerFunc) Log(ctx context.Context, entry Entry) { fn(ctx, entry) }

func Log(l Logger) collect.Collector {
	if l == nil {
		return nil
	}

	return collect.CollectorFunc(func(ctx context.Context, met collect.Metrics, req *http.Request) {
		l.Log(ctx, Entry{
			Request: req,
			Metrics: met,
		})
	})
}

var nop Logger = new(nopLogger)

type nopLogger struct{}

func NopLogger() Logger { return nop }

func (*nopLogger) Log(context.Context, Entry) {}
