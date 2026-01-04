// Copyright (c) 2026, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-pogo/serv"
	"github.com/go-pogo/serv/response"
	"github.com/stretchr/testify/assert"
)

func TestAddRequestID(t *testing.T) {
	t.Run("add default", func(t *testing.T) {
		ctx := serv.ContextWithInfo(context.Background(), serv.Info{})
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

		AddRequestID(response.NoopHandler()).ServeHTTP(nil, req)
		assert.NotEmpty(t, serv.RequestID(ctx))
	})
	t.Run("from header", func(t *testing.T) {
		const want = "foobar"

		ctx := serv.ContextWithInfo(context.Background(), serv.Info{})
		req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		req.Header.Set(headerRequestID, want)

		AddRequestID(response.NoopHandler()).ServeHTTP(nil, req)
		assert.Equal(t, want, serv.RequestID(ctx))
	})
}
