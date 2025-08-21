// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package response

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteJSON(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		rec := httptest.NewRecorder()
		assert.NoError(t, WriteJSON(rec, nil))
		assert.NotEqual(t, rec.Header().Get("Content-Type"), contentTypeJSON)
		assert.Equal(t, rec.Body.String(), "")

	})
	t.Run("good", func(t *testing.T) {
		rec := httptest.NewRecorder()
		assert.NoError(t, WriteJSON(rec, struct{ Foo string }{"bar"}))
		assert.Equal(t, rec.Header().Get("Content-Type"), contentTypeJSON)
		assert.Equal(t, rec.Body.String(), `{"Foo":"bar"}`+"\n")
	})
	t.Run("bad", func(t *testing.T) {
		rec := httptest.NewRecorder()
		assert.Error(t, WriteJSON(rec, make(chan int)))
	})
}
