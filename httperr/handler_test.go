// Copyright (c) 2026, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httperr

import (
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-pogo/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	const wantErr errors.Msg = "some error"
	handler := HandlerFunc(func(wri http.ResponseWriter, req *http.Request) error {
		return wantErr
	})

	t.Run("with nil handler", func(t *testing.T) {
		var out strings.Builder

		srv := httptest.NewServer(HandleError(handler, nil))
		srv.Config.ErrorLog = log.New(&out, "", 0)
		_, _ = srv.Client().Get(srv.URL)

		// close server before getting the output string to prevent a race condition
		srv.Close()

		have := out.String()
		assert.Contains(t, have, "http: panic serving")
		assert.Contains(t, have, wantErr.String())
	})

	t.Run("with slog handler", func(t *testing.T) {
		var out strings.Builder
		logger := slog.New(slog.NewTextHandler(&out, &slog.HandlerOptions{}))

		srv := httptest.NewServer(HandleError(handler, Log(logger)))
		_, _ = srv.Client().Get(srv.URL)
		srv.Close()

		have := out.String()
		assert.Contains(t, have, `level=ERROR`)
		assert.Contains(t, have, `msg="handler error"`)
		assert.Contains(t, have, `err="some error"`)
	})

	t.Run("with err handler", func(t *testing.T) {
		var have error
		h := func(err error) { have = err }

		srv := httptest.NewServer(HandleError(handler, h))
		_, _ = srv.Client().Get(srv.URL)
		srv.Close()

		assert.Equal(t, wantErr, have)
	})

	t.Run("panic on nil handler", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilNextHandler, func() {
			HandleError(nil, nil)
		})
	})
}
