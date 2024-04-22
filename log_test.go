// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWithLogger(t *testing.T) {
	t.Run("panic on nil", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilLogger, func() {
			require.NoError(t, WithLogger(nil).apply(&Server{}))
		})
	})

	t.Run("set logger", func(t *testing.T) {
		var srv Server
		want := NopLogger()
		assert.NoError(t, WithLogger(want).apply(&srv))
		assert.Same(t, want, srv.log)
	})
}
