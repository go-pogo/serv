// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNewNilLogger, func() { _ = NewLogger(nil) })
	})
}

func TestDefaultLogger(t *testing.T) {
	want := slog.Default()
	assert.Same(t, want, DefaultLogger().(*logger).Logger)
}
