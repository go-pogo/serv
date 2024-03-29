// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestHandlerName(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	require.Nil(t, err)

	var have string
	WithHandlerName("foobar", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		have = HandlerName(req.Context())
	})).ServeHTTP(nil, req)

	assert.Equal(t, "foobar", have)
}
