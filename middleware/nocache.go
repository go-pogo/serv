// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"

	"github.com/go-pogo/serv/response"
)

// NoCache adds middleware that sets headers to prevent caching of the
// response, using [response.NoCache].
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.NoCache(w, r)
		next.ServeHTTP(w, r)
	})
}
