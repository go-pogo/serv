// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServerName(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		assert.Equal(t, "", ServerName(context.Background()))
	})

	t.Run("value", func(t *testing.T) {
		want := "foobar"
		ctx := context.WithValue(context.Background(), serverNameKey{}, want)
		assert.Equal(t, want, ServerName(ctx))
	})
}
