// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestServeMux(t *testing.T) {
	t.Run("apply to server", func(t *testing.T) {
		var srv Server
		var mux ServeMux
		assert.NoError(t, srv.With(&mux))
		assert.Same(t, &mux, srv.Handler)
	})
}

func TestServeMux_HandleRoute(t *testing.T) {
	mux := NewServeMux()
	mux.HandleRoute(Route{
		Pattern: "/",
		Handler: http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
			res.WriteHeader(http.StatusOK)
		}),
	})

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestServeMux_ServeHTTP(t *testing.T) {
	t.Run("default not found", func(t *testing.T) {
		rec := httptest.NewRecorder()
		NewServeMux().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("custom not found", func(t *testing.T) {
		const want = "my custom not found message"
		mux := NewServeMux().
			WithNotFoundHandler(http.HandlerFunc(func(wri http.ResponseWriter, _ *http.Request) {
				_, _ = wri.Write([]byte(want))
			}))

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
		assert.Equal(t, want, rec.Body.String())
	})
}
