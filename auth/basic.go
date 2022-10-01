// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"

	"github.com/go-pogo/serv/httpheader"
)

func Basic(user, pass string, next http.Handler) http.Handler {
	if pass == "" {
		return next
	}

	return NewBasic(user, pass).Wrap(next)
}

type BasicMiddleware struct {
	user, pass [sha256.Size]byte
}

func NewBasic(user, pass string) *BasicMiddleware {
	var h BasicMiddleware
	h.SetUser(user)
	h.SetPass(pass)
	return &h
}

func (h *BasicMiddleware) SetUser(v string) { h.user = sha256.Sum256([]byte(v)) }
func (h *BasicMiddleware) SetPass(v string) { h.pass = sha256.Sum256([]byte(v)) }

func (h *BasicMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		if u, p, submitted := req.BasicAuth(); submitted {
			uh := sha256.Sum256([]byte(u))
			ph := sha256.Sum256([]byte(p))

			if subtle.ConstantTimeCompare(h.user[:], uh[:]) == 1 &&
				subtle.ConstantTimeCompare(h.pass[:], ph[:]) == 1 {
				next.ServeHTTP(wri, req)
				return
			}
		}

		wri.Header().Set(httpheader.WWWAuthenticate, `Basic charset="UTF-8"`)
		http.Error(wri, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	})
}
