// Copyright (c) 2026, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httperr

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-pogo/serv"
	"github.com/go-pogo/serv/response"
)

type ErrHandler interface {
	HandleError(ctx context.Context, err error)
}

type ErrHandlerFunc func(ctx context.Context, err error)

// HandleError calls fn(ctx, err).
func (fn ErrHandlerFunc) HandleError(ctx context.Context, err error) { fn(ctx, err) }

func Log(l *slog.Logger) ErrHandlerFunc {
	return func(ctx context.Context, err error) {
		if l == nil {
			l = slog.Default()
		}
		if !l.Handler().Enabled(ctx, slog.LevelError) {
			return
		}

		attrs := make([]any, 0, 4)
		attrs = append(attrs, slog.Any("err", err))

		if v := serv.ServerName(ctx); v != "" {
			attrs = append(attrs, slog.String("server", v))
		}
		if v := serv.HandlerName(ctx); v != "" {
			attrs = append(attrs, slog.String("handler", v))
		}
		if v := serv.RequestID(ctx); v != "" {
			attrs = append(attrs, slog.String("request_id", v))
		}

		l.ErrorContext(ctx, "handler error", attrs...)
	}
}

const panicNilNextHandler = "httperr: next handler should not be nil"

// HandleError returns an [http.Handler] which wraps the [Handler] next and
// handles any returned errors by it using the provided [ErrHandlerFunc].
func HandleError(next Handler, handleErr ErrHandlerFunc) http.Handler {
	if next == nil {
		panic(panicNilNextHandler)
	}
	if handleErr == nil {
		handleErr = func(_ context.Context, err error) { panic(err) }
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := next.ServeHTTPError(w, r); err != nil {
			handleErr.HandleError(r.Context(), err)
		}
	})
}

func WriteJSONError(next Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		if err := next.ServeHTTPError(wri, req); err != nil {
			_ = response.WriteJSONError(wri, err)
		}
	})
}
