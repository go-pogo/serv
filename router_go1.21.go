// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.22

package serv

import "net/http"

func (mux *ServeMux) handle(_, pattern string, handler http.Handler) {
	mux.serveMux.Handle(pattern, handler)
}

func (mux *ServeMux) handleFunc(_, pattern string, handler http.HandlerFunc) {
	mux.serveMux.HandleFunc(pattern, handler)
}
