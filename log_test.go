// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNewNilLogger, func() { _ = NewLogger(nil) })
	})

	t.Run("custom logger", func(t *testing.T) {
		want := log.New(io.Discard, "", 0)
		l := NewLogger(want)
		assert.Same(t, want, l.(*logger).Logger)
	})

	t.Run("ErrorLoggerProvider", func(t *testing.T) {
		want := log.Default()
		assert.Same(t, want, NewLogger(want).ErrorLogger())
	})
}

func TestDefaultLogger(t *testing.T) {
	want := log.Default()
	assert.Same(t, want, DefaultLogger().(*logger).Logger)
}
