// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultServeMux(t *testing.T) {
	DefaultServeMux().HandleRoute(Route{
		Pattern: "/test",
		Handler: http.HandlerFunc(func(res http.ResponseWriter, _q *http.Request) {
			res.WriteHeader(http.StatusOK)
		}),
	})

	rec := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/test", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRoute_ServeHTTP(t *testing.T) {
	tests := map[string]string{
		"without name": "",
		"with name":    "test",
	}
	for name, want := range tests {
		t.Run(name, func(t *testing.T) {
			var haveName string
			route := Route{
				Name: want,
				Handler: http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
					haveName = HandlerName(req.Context())
				}),
			}
			route.ServeHTTP(nil, httptest.NewRequest(http.MethodGet, "/", nil))
			assert.Equal(t, want, haveName)
		})
	}
}

func TestRouter_HandleRoute(t *testing.T) {
	router := NewServeMux()
	router.HandleRoute(Route{
		Pattern: "/",
		Handler: http.HandlerFunc(func(res http.ResponseWriter, _q *http.Request) {
			res.WriteHeader(http.StatusOK)
		}),
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
}
