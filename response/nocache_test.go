// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoCache(t *testing.T) {
	t.Run("headers present", func(t *testing.T) {
		rec := httptest.NewRecorder()
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			NoCache(w, r)
			w.WriteHeader(http.StatusOK)
		}).ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))

		for k, v := range noCacheHeaders {
			assert.Equalf(t, v, rec.Header().Get(k), "%s header not set by middleware.", k)
		}
	})
}
