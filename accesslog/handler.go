// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/go-pogo/serv"
)

func Middleware(log Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return NewHandler(next, log)
	}
}

var _ http.Handler = (*handler)(nil)

type handler struct {
	log     Logger
	next    http.Handler
	traffic int64
}

// NewHandler wraps a [http.Handler] so it's request and response details are
// tracked and send to [Logger] log.
func NewHandler(next http.Handler, log Logger) http.Handler {
	if log == nil {
		log = NopLogger()
	}

	return &handler{
		log:  log,
		next: next,
	}
}

func (h *handler) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	ctx, settings, existing := withSettings(req.Context())
	if !existing {
		req = req.WithContext(ctx)
	} else if settings.shouldIgnore {
		// when set outside handler h
		h.next.ServeHTTP(wri, req)
		return
	}

	var det Details
	det.StartTime = time.Now()
	det.RequestCount = atomic.AddInt64(&h.traffic, 1)
	defer atomic.AddInt64(&h.traffic, -1)

	met := httpsnoop.CaptureMetrics(h.next, wri, req)
	if settings.shouldIgnore {
		// when set inside h.next
		return
	}

	if info := serv.InfoFromContext(ctx); info != nil {
		det.Info = *info
	}

	det.StatusCode = met.Code
	det.Duration = met.Duration
	det.BytesWritten = met.Written
	det.UserAgent = req.UserAgent()

	h.log.LogAccess(ctx, det, req.Clone(&noopCtx{ctx}))
}

type ctxSettingsKey struct{}

type ctxSettings struct {
	shouldIgnore bool
}

func withSettings(ctx context.Context) (context.Context, *ctxSettings, bool) {
	if v := ctx.Value(ctxSettingsKey{}); v != nil {
		return ctx, v.(*ctxSettings), true
	}

	v := new(ctxSettings)
	return context.WithValue(ctx, ctxSettingsKey{}, v), v, false
}

var _ context.Context = (*noopCtx)(nil)

// A noopCtx is similar to context.Background as it is never canceled and has
// no deadline. However, it does return values from its parent context, when
// available.
type noopCtx struct {
	parent context.Context
}

func (*noopCtx) Deadline() (deadline time.Time, ok bool) { return }

func (*noopCtx) Done() <-chan struct{} { return nil }

func (*noopCtx) Err() error { return nil }

func (c *noopCtx) Value(key interface{}) interface{} { return c.parent.Value(key) }
