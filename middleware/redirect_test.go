// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-pogo/serv/response"
	"github.com/stretchr/testify/assert"
)

func TestRedirectHTTPS(t *testing.T) {
	tests := []struct {
		method       string
		target       string
		wantCode     int
		wantLocation string
	}{
		{
			method:   http.MethodGet,
			target:   "https://example.com",
			wantCode: http.StatusOK,
		},
		{
			method:       http.MethodGet,
			target:       "http://example.com",
			wantCode:     http.StatusMovedPermanently,
			wantLocation: "https://example.com",
		},
		{
			method:       http.MethodPost,
			target:       "http://example.com/post-me",
			wantCode:     http.StatusTemporaryRedirect,
			wantLocation: "https://example.com/post-me",
		},
	}

	handler := RedirectHTTPS(response.NoopHandler())
	for _, tc := range tests {
		t.Run(tc.method+" "+tc.target, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.target, nil))
			assert.Equal(t, tc.wantCode, rec.Code)
			assert.Equal(t, tc.wantLocation, rec.Header().Get("Location"))
		})
	}
}

func TestRemoveTrailingSlash(t *testing.T) {
	tests := []struct {
		method       string
		target       string
		wantCode     int
		wantLocation string
	}{
		{
			method:       http.MethodGet,
			target:       "/test/",
			wantCode:     http.StatusMovedPermanently,
			wantLocation: "/test",
		},
		{
			method:       http.MethodPost,
			target:       "/test/",
			wantCode:     http.StatusTemporaryRedirect,
			wantLocation: "/test",
		},
		{
			method:   http.MethodGet,
			target:   "/no-redirect",
			wantCode: http.StatusOK,
		},
		{
			method:   http.MethodGet,
			target:   "http://example.com",
			wantCode: http.StatusOK,
		},
		{
			method:   http.MethodGet,
			target:   "http://example.com/",
			wantCode: http.StatusOK,
		},
	}

	handler := RemoveTrailingSlash(response.NoopHandler())
	for _, tc := range tests {
		t.Run(tc.method+" "+tc.target, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.target, nil))
			assert.Equal(t, tc.wantCode, rec.Code)
			assert.Equal(t, tc.wantLocation, rec.Header().Get("Location"))
		})
	}
}
