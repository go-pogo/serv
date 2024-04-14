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

func TestRoute_ServeHTTP(t *testing.T) {
	var haveName string
	route := Route{
		Name: "test",
		Handler: http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
			haveName = HandlerName(req.Context())
		}),
	}
	route.ServeHTTP(nil, httptest.NewRequest(http.MethodGet, "/", nil))
	assert.Equal(t, route.Name, haveName)
}
