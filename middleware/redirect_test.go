// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedirectHTTPS(t *testing.T) {
	const wantBody = "Hello, World!"
	handler := RedirectHTTPS(http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		_, _ = wri.Write([]byte(wantBody))
	}))

	t.Run("redirect", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest("GET", "http://example.com", nil))
		assert.Equal(t, http.StatusMovedPermanently, rec.Code)
		assert.Equal(t, "https://example.com", rec.Header().Get("Location"))
	})
	t.Run("no redirect", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest("GET", "https://example.com", nil))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "", rec.Header().Get("Location"))
		assert.Equal(t, wantBody, rec.Body.String())
	})
}

func TestRemoveTrailingSlash(t *testing.T) {
	const wantBody = "Hello, World!"
	handler := RemoveTrailingSlash(http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		_, _ = wri.Write([]byte(wantBody))
	}))

	t.Run("redirect", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest("GET", "/test/", nil))
		assert.Equal(t, http.StatusMovedPermanently, rec.Code)
		assert.Equal(t, "/test", rec.Header().Get("Location"))
	})
	t.Run("no redirect", func(t *testing.T) {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest("GET", "http://example.com", nil))
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "", rec.Header().Get("Location"))
		assert.Equal(t, wantBody, rec.Body.String())
	})
}
