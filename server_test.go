// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	have, err := New()
	assert.NoError(t, err)
	assert.Equal(t, *DefaultConfig(), have.Config)
}

func TestServer_State(t *testing.T) {
	t.Run("not started", func(t *testing.T) {
		var srv Server
		assert.Equal(t, StateUnstarted, srv.State())
	})
	t.Run("started", func(t *testing.T) {
		var srv Server
		require.NoError(t, srv.start())
		assert.Equal(t, StateStarted, srv.State())
		assert.ErrorIs(t, srv.start(), ErrAlreadyStarted)
	})
}
