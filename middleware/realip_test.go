// Copyright (c) 2026, Roel Schut. All rights reserved.
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

func TestRealIP(t *testing.T) {
	tests := map[string]struct {
		header http.Header
		want   string
	}{
		"skip local, single header": {
			header: http.Header{
				headerForwardedFor: {"127.0.0.1, 25.26.27.28"},
			},
			want: "25.26.27.28",
		},
		"skip local, multiple headers": {
			header: http.Header{
				headerForwardedFor: {"127.0.0.1", "25.26.27.28"},
			},
			want: "25.26.27.28",
		},
		"realip header": {
			header: http.Header{
				headerForwardedFor: {""},
				headerRealIP:       {"25.26.27.28"},
			},
			want: "25.26.27.28",
		},
		"default": {
			want: "192.0.2.1:1234",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			for k, v := range tc.header {
				req.Header[k] = v
			}

			RealIP(response.NoopHandler()).ServeHTTP(nil, req)
			assert.Equal(t, tc.want, req.RemoteAddr)
		})
	}
}
