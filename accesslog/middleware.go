// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"context"
	"github.com/felixge/httpsnoop"
	"github.com/go-pogo/serv"
	"github.com/go-pogo/serv/middleware"
	"net/http"
	"sync/atomic"
	"time"
)

type handler struct {
	log     Logger
	next    http.Handler
	traffic int64
}

// Wrap wraps a http.Handler so it's request and response details are tracked
// and send to Logger log.
func Wrap(log Logger, next http.Handler) http.Handler {
	if log == nil {
		log = NopLogger()
	}
	return &handler{
		log:  log,
		next: next,
	}
}

// Middleware returns Wrap as middleware.Middleware.
func Middleware(log Logger) middleware.Middleware {
	return middleware.MiddlewareFunc(func(next http.HandlerFunc) http.Handler {
		return Wrap(log, next)
	})
}

func (c *handler) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	var det Details
	det.StartTime = time.Now()
	det.RequestCount = atomic.AddInt64(&c.traffic, 1)
	defer atomic.AddInt64(&c.traffic, -1)

	// create a string pointer for the handler name, so we can populate
	// this with the correct value in the handler itself
	ctx := setHandlerName(req.Context(), new(string))
	met := httpsnoop.CaptureMetrics(c.next, wri, req.WithContext(ctx))

	det.StatusCode = met.Code
	det.Duration = met.Duration
	det.BytesWritten = met.Written
	det.ServerName = serv.ServerName(ctx)
	det.HandlerName = HandlerName(ctx)
	det.UserAgent = req.UserAgent()

	c.log.Log(ctx, det, req.Clone(&noopCtx{ctx}))
}

var _ context.Context = new(noopCtx)

// An noopCtx is similar to context.Background as it is never canceled and has
// no deadline. However, it does return values from its parent context, when
// available.
type noopCtx struct {
	parent context.Context
}

func (*noopCtx) Deadline() (deadline time.Time, ok bool) { return }

func (*noopCtx) Done() <-chan struct{} { return nil }

func (*noopCtx) Err() error { return nil }

func (c *noopCtx) Value(key interface{}) interface{} { return c.parent.Value(key) }
