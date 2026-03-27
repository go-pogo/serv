// Copyright (c) 2026, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httperr

import (
	"net/http"
	"testing"

	"github.com/go-pogo/errors"
	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	t.Run("panic on nil handler", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilNextHandler, func() {
			Wrap(nil)
		})
	})

	t.Run("error from panic", func(t *testing.T) {
		const wantErr errors.Msg = "some panic"
		haveErr := Wrap(func(_ http.ResponseWriter, _ *http.Request) {
			panic(wantErr)
		}).ServeHTTPError(nil, nil)

		// todo: in toekomst PanicError.Value testen
		assert.Contains(t, haveErr.Error(), wantErr.String())
	})
}
