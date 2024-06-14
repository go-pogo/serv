// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	have, err := New()
	assert.NoError(t, err)
	assert.Equal(t, *DefaultConfig(), have.Config)
}

func TestServer_With(t *testing.T) {
	t.Run("nil option", func(t *testing.T) {
		var srv Server
		assert.NoError(t, srv.With(nil))
	})

	t.Run("started", func(t *testing.T) {
		var srv Server
		require.NoError(t, srv.start())

		var wantErr *InvalidStateError
		assert.ErrorAs(t, srv.With(WithName("foobar")), &wantErr)
		assert.Equal(t, ErrAlreadyStarted, wantErr.Unwrap())
		assert.Equal(t, StateStarted, wantErr.State)
	})
}

func TestServer_State(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		var srv Server
		assert.Equal(t, StateUnstarted, srv.State())
	})
	t.Run("start", func(t *testing.T) {
		var srv Server
		require.NoError(t, srv.start())
		assert.Equal(t, StateStarted, srv.State())
		assert.ErrorIs(t, srv.start(), ErrUnableToStart)
	})
	t.Run("shutdown", func(t *testing.T) {
		var srv Server
		require.NoError(t, srv.start())
		assert.NoError(t, srv.Shutdown(context.Background()))
		assert.Equal(t, StateClosed, srv.State())
		assert.ErrorIs(t, srv.Shutdown(context.Background()), ErrUnableToShutdown)
	})
	t.Run("close", func(t *testing.T) {
		var srv Server
		require.NoError(t, srv.start())
		assert.NoError(t, srv.Close())
		assert.Equal(t, StateClosed, srv.State())
		assert.ErrorIs(t, srv.Close(), ErrUnableToClose)
	})
}
