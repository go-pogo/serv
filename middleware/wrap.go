// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
)

// Wrapper wraps the [http.Handler] next with additional logic.
type Wrapper func(next http.Handler) http.Handler

// Wrap [http.Handler] h with additional logic via the provided [Wrapper]s.
func Wrap(h http.Handler, wrap ...Wrapper) http.Handler {
	for i := len(wrap) - 1; i >= 0; i-- {
		if wrap[i] == nil {
			continue
		}

		h = wrap[i](h)
	}
	return h
}

// Middleware consists of multiple [Wrapper]s which can be wrapped around a
// [http.Handler] using [Wrap].
type Middleware []Wrapper

func (m Middleware) Wrap(next http.Handler) http.Handler {
	return Wrap(next, m...)
}
