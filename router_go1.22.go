// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.22

package serv

import "net/http"

func (r *ServeMux) handle(method, pattern string, handler http.Handler) {
	r.serveMux.Handle(method+" "+pattern, handler)
}

func (r *ServeMux) handleFunc(method, pattern string, handler http.HandlerFunc) {
	r.serveMux.HandleFunc(method+" "+pattern, handler)
}
