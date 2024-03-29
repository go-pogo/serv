// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"context"
	"github.com/felixge/httpsnoop"
	"github.com/go-pogo/serv/internal"
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
func Middleware(log Logger) middleware.Wrapper {
	return middleware.WrapperFunc(func(next http.HandlerFunc) http.Handler {
		return Wrap(log, next)
	})
}

func (c *handler) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	var det Details
	det.StartTime = time.Now()
	det.RequestCount = atomic.AddInt64(&c.traffic, 1)
	defer atomic.AddInt64(&c.traffic, -1)

	ctx, settings, existing := withSettings(req.Context())
	if !existing {
		req = req.WithContext(ctx)
	}

	met := httpsnoop.CaptureMetrics(c.next, wri, req)
	if settings.ShouldIgnore {
		return
	}

	det.StatusCode = met.Code
	det.Duration = met.Duration
	det.BytesWritten = met.Written
	det.ServerName = internal.ServerName(ctx)
	det.HandlerName = settings.HandlerName
	det.UserAgent = req.UserAgent()

	c.log.Log(ctx, det, req.Clone(&noopCtx{ctx}))
}

var _ context.Context = (*noopCtx)(nil)

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
