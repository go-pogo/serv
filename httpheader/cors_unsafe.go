// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build cors_unsafe
// +build cors_unsafe

package httpheader

import (
	"net/http"
)

// extractOrigin get the origin from the request's Origin of Referer headers.
func extractOrigin(h http.Header) string {
	origin := h.Get(Origin)
	if origin == "" {
		origin = h.Get(Referer)
	}
	return origin
}
