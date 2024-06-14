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

func TestDefaultLogger(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		l := DefaultLogger(nil)
		assert.Same(t, log.Default(), l.(*defaultLogger).Logger)
	})

	t.Run("custom logger", func(t *testing.T) {
		want := log.New(io.Discard, "", 0)
		l := DefaultLogger(want)
		assert.Same(t, want, l.(*defaultLogger).Logger)
	})

	t.Run("ErrorLoggerProvider", func(t *testing.T) {
		want := log.Default()
		assert.Same(t, want, DefaultLogger(nil).ErrorLogger())
	})
}
