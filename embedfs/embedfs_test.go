// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package embedfs

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed _test/*
var embedded embed.FS

func TestNew(t *testing.T) {
	handler, err := New(embedded)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest("GET", "/_test/some-file.txt", nil))
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "some-file.txt", strings.TrimSpace(rec.Body.String()))
	assert.Equal(t, "", rec.Header().Get("Last-Modified"))
}

func TestWithSubDir(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		handler, err := New(embedded, WithSubDir("_test"))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest("GET", "/some-file.txt", nil))
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "some-file.txt", strings.TrimSpace(rec.Body.String()))
	})
	t.Run("empty", func(t *testing.T) {
		handler, err := New(embedded, WithSubDir(""))
		require.NoError(t, err)

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest("GET", "/some-file.txt", nil))
		assert.Equal(t, 404, rec.Code)
	})
	t.Run("invalid", func(t *testing.T) {
		_, err := New(embedded, WithSubDir("../invalid"))
		assert.ErrorIs(t, err, ErrInvalidSubDir)
	})
}

func TestWithModTime(t *testing.T) {
	now := time.Now()
	handler, err := New(embedded, WithModTime(now))
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest("GET", "/_test/some-file.txt", nil))
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "some-file.txt", strings.TrimSpace(rec.Body.String()))
	assert.Equal(t, now.UTC().Format(http.TimeFormat), rec.Header().Get("Last-Modified"))
	assert.Equal(t, now.UTC(), handler.ModTime())
}
