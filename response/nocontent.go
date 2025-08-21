// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package response

import (
	"net/http"
)

// NoContent replies to the request with an HTTP 204 "no content" status code.
func NoContent(wri http.ResponseWriter) {
	wri.WriteHeader(http.StatusNoContent)
}

// NoContentHandler returns a simple request handler that replies to each
// request with an HTTP 204 “no content” reply.
func NoContentHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		NoContent(w)
	})
}
