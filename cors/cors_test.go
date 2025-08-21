// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var nopHttpHandlerFunc = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
})

func TestMiddleware(t *testing.T) {
	tests := map[string]struct {
		allow      []string
		header     http.Header
		wantSafe   string
		wantUnsafe string
	}{
		"origin header": {
			allow: []string{"*"},
			header: map[string][]string{
				originKey:  {"https://foo.bar"},
				refererKey: {"https://qux.xoo"},
			},
			wantSafe:   "*",
			wantUnsafe: "https://foo.bar",
		},
		"referer header": {
			allow: []string{"*"},
			header: map[string][]string{
				refererKey: {"https://qux.xoo"},
			},
			wantSafe:   "*",
			wantUnsafe: "https://qux.xoo",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			Middleware(AccessControl{
				AllowOrigin: tc.allow,
			})(nopHttpHandlerFunc).
				ServeHTTP(rec, &http.Request{
					Header: tc.header,
				})

			have := rec.Header().Get(AccessControlAllowOriginKey)
			if corsUnsafe {
				assert.Equal(t, tc.wantUnsafe, have)
			} else {
				assert.Equal(t, tc.wantSafe, have)
			}
		})
	}
}
