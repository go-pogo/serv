// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-pogo/serv/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextWithInfo(t *testing.T) {
	t.Run("add", func(t *testing.T) {
		want := Info{
			ServerName:  "foo",
			HandlerName: "bar",
		}

		ctx := context.Background()
		assert.Nil(t, InfoFromContext(ctx))

		have := InfoFromContext(ContextWithInfo(ctx, want))
		assert.NotNil(t, have)
		assert.Equal(t, want, *have)
	})
	t.Run("replace", func(t *testing.T) {
		want := Info{
			ServerName:  "foo",
			HandlerName: "bar",
		}

		ctx := ContextWithInfo(context.Background(), Info{ServerName: "initial"})
		assert.Equal(t, "initial", ServerName(ctx))

		_ = ContextWithInfo(ctx, want)
		assert.Equal(t, want, *InfoFromContext(ctx))
	})
}

func TestServerName(t *testing.T) {
	t.Run("empty context", func(t *testing.T) {
		assert.Equal(t, "", ServerName(context.Background()))
	})
	t.Run("from context", func(t *testing.T) {
		const want = "foobar"
		assert.Equal(t, want, ServerName(ContextWithInfo(
			context.Background(),
			Info{ServerName: want},
		)))
	})
	t.Run("server handler", func(t *testing.T) {
		const want = "quxoo"
		srv, err := New(WithName(want))
		require.NoError(t, err)

		srv.Handler = http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
			assert.Equal(t, want, ServerName(req.Context()))
		})
		require.NoError(t, srv.start())
		srv.httpServer.Handler.ServeHTTP(nil, httptest.NewRequest(http.MethodGet, "/", nil))
	})
}

func TestHandlerName(t *testing.T) {
	t.Run("empty context", func(t *testing.T) {
		assert.Equal(t, "", HandlerName(context.Background()))
	})
	t.Run("from context", func(t *testing.T) {
		const want = "foobar"
		assert.Equal(t, want, HandlerName(ContextWithInfo(
			context.Background(),
			Info{HandlerName: want},
		)))
	})
	t.Run("route handler", func(t *testing.T) {
		const want = "quxoo"
		route := Route{
			Name:    want,
			Method:  http.MethodGet,
			Pattern: "/",
			Handler: http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
				assert.Equal(t, want, HandlerName(req.Context()))
			}),
		}
		route.ServeHTTP(nil, httptest.NewRequest(route.Method, route.Pattern, nil))
	})
}

func TestRequestID(t *testing.T) {
	t.Run("empty context", func(t *testing.T) {
		assert.Equal(t, "", RequestID(context.Background()))
	})
	t.Run("from context", func(t *testing.T) {
		const want = "foobar"
		assert.Equal(t, want, RequestID(ContextWithInfo(
			context.Background(),
			Info{RequestID: want},
		)))
	})
}

func testInfoHandlers(t *testing.T,
	addFn func(string, http.Handler) http.Handler,
	getFn func(context.Context) string,
) {
	t.Run("inside handler", func(t *testing.T) {
		var have string
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			have = getFn(req.Context())
		})

		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
		addFn("foobar", handler).ServeHTTP(nil, req)
		assert.Equal(t, "foobar", have)
	})

	t.Run("outside handler", func(t *testing.T) {
		ctx := ContextWithInfo(context.Background(), Info{})
		assert.Empty(t, getFn(ctx))

		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		addFn("foobar", response.NoopHandler()).ServeHTTP(nil, req)
		assert.Equal(t, "foobar", getFn(ctx))
	})
}

func TestAddServerName(t *testing.T) {
	testInfoHandlers(t, AddServerName, ServerName)
}

func TestAddHandlerName(t *testing.T) {
	testInfoHandlers(t, AddHandlerName, HandlerName)
}

func TestAddRequestID(t *testing.T) {
	testInfoHandlers(t, AddRequestID, RequestID)
}
