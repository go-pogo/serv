// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.22

package serv

func (mux *ServeMux) handle(route Route) {
	mux.serveMux.Handle(route.Pattern, route.Handler)
}
