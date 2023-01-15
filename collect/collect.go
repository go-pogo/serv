// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collect

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/go-pogo/serv/middleware"
)

type snoopMetrics = httpsnoop.Metrics

// Metrics holds metrics collected with Intercept.
type Metrics struct {
	snoopMetrics
	// The Time the request was received.
	Time time.Time
	// Traffic is the amount of open requests during the execution of the
	// handler.
	Traffic int64
}

// SnoopMetrics returns the original underlying httpsnoop.Metrics.
func (m Metrics) SnoopMetrics() httpsnoop.Metrics { return m.snoopMetrics }

// Collector receives the Metrics data and decides what to do with it. This may
// range from simple logging to making it available for scraping (prometheus).
type Collector interface {
	Collect(ctx context.Context, met Metrics, req *http.Request)
}

type CollectorFunc func(ctx context.Context, met Metrics, req *http.Request)

func (fn CollectorFunc) Collect(ctx context.Context, met Metrics, req *http.Request) {
	fn(ctx, met, req)
}

type collector struct {
	next    http.Handler
	col     []Collector
	traffic int64
}

// Wrap metrics from http.Handler h and pass it to a Collector.
func Wrap(h http.Handler, col ...Collector) http.Handler {
	if len(col) == 0 {
		return h
	}

	return &collector{
		next: h,
		col:  col,
	}
}

func Middleware(col ...Collector) middleware.Middleware {
	return middleware.MiddlewareFunc(func(next http.HandlerFunc) http.Handler {
		return Wrap(next, col...)
	})
}

func (c *collector) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	var m Metrics
	m.Time = time.Now()
	m.Traffic = atomic.AddInt64(&c.traffic, 1)
	defer atomic.AddInt64(&c.traffic, -1)

	m.snoopMetrics = httpsnoop.CaptureMetrics(c.next, wri, req)

	go func() {
		ctx := &noopCtx{req.Context()}
		req = req.Clone(ctx)

		for _, col := range c.col {
			if col != nil {
				col.Collect(ctx, m, req)
			}
		}
	}()
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
