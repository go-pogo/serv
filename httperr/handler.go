// Copyright (c) 2026, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httperr

import (
	"log/slog"
	"net/http"

	"github.com/go-pogo/serv/response"
)

type ErrHandler interface {
	HandleError(err error)
}

type ErrHandlerFunc func(err error)

func (fn ErrHandlerFunc) HandleError(err error) { fn(err) }

func Log(l *slog.Logger) ErrHandlerFunc {
	return func(err error) {
		if l == nil {
			l = slog.Default()
		}
		l.Error("handler error", slog.Any("err", err))
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
		handleErr = func(err error) { panic(err) }
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := next.ServeHTTPError(w, r); err != nil {
			handleErr.HandleError(err)
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
