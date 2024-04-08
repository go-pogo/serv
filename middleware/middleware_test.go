// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrap(t *testing.T) {
	handler := Wrap(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte("c"))
		}),
		WrapperFunc(func(next http.HandlerFunc) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write([]byte("a"))
				next(w, nil)
			})
		}),
		WrapperFunc(func(next http.HandlerFunc) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write([]byte("b"))
				next(w, nil)
			})
		}),
	)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, nil)
	assert.Equal(t, "abc", rec.Body.String())
}

func TestHandlerFunc_Wrap(t *testing.T) {
	handlerA := func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("a"))
	}
	handlerB := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("b"))
	})

	rec := httptest.NewRecorder()
	HandlerFunc(handlerA).Wrap(handlerB).ServeHTTP(rec, nil)
	assert.Equal(t, "ab", rec.Body.String())
}
