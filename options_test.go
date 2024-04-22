// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestWithHandler(t *testing.T) {
	var srv Server
	want := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	assert.NoError(t, WithHandler(want).apply(&srv))
	assert.Equal(t, fmt.Sprintf("%v", want), fmt.Sprintf("%v", srv.Handler))
}

func TestWithRegisterRoutes(t *testing.T) {
	t.Run("nil handler", func(t *testing.T) {
		var srv Server
		assert.NoError(t, WithRoutesRegisterer().apply(&srv))
		assert.NotNil(t, srv.Handler)
	})
	t.Run("no routes handler", func(t *testing.T) {
		srv := Server{Handler: http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})}
		assert.ErrorIs(t, WithRoutesRegisterer().apply(&srv), ErrHandlerIsNoRouteHandler)
	})
}

func TestWithName(t *testing.T) {
	var srv Server
	assert.NoError(t, WithName("foobar").apply(&srv))
	assert.Equal(t, "foobar", srv.Name())
}

func TestWithTLS(t *testing.T) {
	t.Run("panic on nil", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilTLSConfig, func() {
			require.NoError(t, WithTLS(nil).apply(&Server{}))
		})
	})
}

func TestWithDefaultTLSConfig(t *testing.T) {
	var srv Server
	assert.NoError(t, WithDefaultTLSConfig().apply(&srv))
	assert.Equal(t, DefaultTLSConfig(), srv.TLSConfig)
}
