// Copyright (c) 2022, Roel Schut. All rights reserved.
// applyOptions of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"net/http"
)

type Header interface {
	Add(key, value string)
	Set(key, value string)
}

type HeaderModifier interface {
	ModifyHeader(s Header)
}

func WithHeader(h HeaderModifier, next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		h.ModifyHeader(response.Header())
		next.ServeHTTP(response, request)
	})
}
