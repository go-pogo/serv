// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collect

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandlerName(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		assert.Equal(t, "", HandlerName(context.Background()))
	})

	t.Run("value", func(t *testing.T) {
		want := "foobar"
		ctx := context.WithValue(context.Background(), handlerNameKey{}, want)
		assert.Equal(t, want, HandlerName(ctx))
	})
}

func TestHandlerNameOr(t *testing.T) {
	t.Run("empty value", func(t *testing.T) {
		assert.Equal(t, "foo", HandlerNameOr(context.Background(), "foo"))
	})

	t.Run("value", func(t *testing.T) {
		want := "foobar"
		ctx := context.WithValue(context.Background(), handlerNameKey{}, want)
		assert.Equal(t, want, HandlerNameOr(ctx, "xoo"))
	})
}
