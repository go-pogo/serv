// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWithHandler(t *testing.T) {
	var srv Server
	want := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	assert.NoError(t, WithHandler(want).apply(&srv))
	assert.Equal(t, fmt.Sprintf("%v", want), fmt.Sprintf("%v", srv.Handler))
}

func TestWithRoutes(t *testing.T) {
	t.Run("nil handler", func(t *testing.T) {
		var srv Server
		assert.NoError(t, WithRoutes().apply(&srv))
		assert.NotNil(t, srv.Handler)
	})
}

func TestWithName(t *testing.T) {
	var srv Server
	assert.NoError(t, WithName("foobar").apply(&srv))
	assert.Equal(t, "foobar", srv.Name())
}
