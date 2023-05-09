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

// Details are collected With Wrap and contain additional details of a request
// and it's corresponding response.
type Details struct {
	ServerName  string
	HandlerName string

	// StatusCode is the first http response code passed to the WriteHeader func
	// of the ResponseWriter. If no such call is made, a default code of 200 is
	// assumed instead.
	StatusCode int
	// StartTime is the time the request was received.
	StartTime time.Time
	// Duration is the time it took to execute the handler.
	Duration time.Duration
	// BytesWritten is the number of bytes successfully written by the Write or
	// ReadFrom function of the ResponseWriter. ResponseWriters may also write
	// data to their underlying connection directly (e.g. headers), but those
	// are not tracked. Therefor the number of BytesWritten bytes will usually
	// match the size of the response body.
	BytesWritten int64
	// RequestCount is the amount of open requests during the execution of the
	// handler.
	RequestCount int64
}

type handler struct {
	log     Logger
	next    http.Handler
	traffic int64
}

// Wrap wraps a http.Handler so it's request and response details are tracked
// and send to Logger log.
func Wrap(log Logger, next http.Handler) http.Handler {
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

	met := httpsnoop.CaptureMetrics(c.next, wri, req)
	det.StatusCode = met.Code
	det.Duration = met.Duration
	det.BytesWritten = met.Written

	ctx := req.Context()
	det.ServerName = serv.ServerName(ctx)
	det.HandlerName = HandlerName(ctx)

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
