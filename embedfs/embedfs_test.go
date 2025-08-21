// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package embedfs

import (
	"embed"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed _test-fixtures/*
var embedded embed.FS

func TestNew(t *testing.T) {
	handler, err := New(embedded)
	require.NoError(t, err)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest("GET", "/_test-fixtures/some-file.txt", nil))
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "some-file.txt", strings.TrimSpace(rec.Body.String()))
	assert.Equal(t, "", rec.Header().Get("Last-Modified"))
}
