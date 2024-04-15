// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestAddServerName(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	t.Run("normal", func(t *testing.T) {
		var have string
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			have = ServerName(req.Context())
		})

		AddServerName("foobar", handler).ServeHTTP(nil, req)
		assert.Equal(t, "foobar", have)
	})
	t.Run("wrapped", func(t *testing.T) {
		var have string
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			have = ServerName(req.Context())
		})

		AddHandlerName("myhandler", AddServerName("foobar", handler)).ServeHTTP(nil, req)
		assert.Equal(t, "foobar", have)
	})
}

func TestAddHandlerName(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	require.NoError(t, err)

	t.Run("normal", func(t *testing.T) {
		var have string
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			have = HandlerName(req.Context())
		})

		AddHandlerName("foobar", handler).ServeHTTP(nil, req)
		assert.Equal(t, "foobar", have)
	})
	t.Run("wrapped", func(t *testing.T) {
		var have string
		handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			have = HandlerName(req.Context())
		})

		AddServerName("myserver", AddHandlerName("foobar", handler)).ServeHTTP(nil, req)
		assert.Equal(t, "foobar", have)
	})
}
