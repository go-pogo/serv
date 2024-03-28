// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package embedfs

import (
	"embed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

//go:embed test/*
var embedded embed.FS

func TestNew(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		handler, err := New(embedded)
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest("GET", "/test/some-file.txt", nil))
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "some-file.txt\n", rec.Body.String())
		assert.Equal(t, "", rec.Header().Get("Last-Modified"))
	})

	t.Run("WithSubDir", func(t *testing.T) {
		handler, err := New(embedded, WithSubDir("test"))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest("GET", "/some-file.txt", nil))
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "some-file.txt\n", rec.Body.String())
	})

	t.Run("WithModTime", func(t *testing.T) {
		now := time.Now()
		handler, err := New(embedded, WithModTime(now))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest("GET", "/test/some-file.txt", nil))
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "some-file.txt\n", rec.Body.String())
		assert.Equal(t, now.UTC().Format(http.TimeFormat), rec.Header().Get("Last-Modified"))
	})
}
