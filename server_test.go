// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDefault(t *testing.T) {
	var want Server
	_ = DefaultConfig().apply(&want)

	have, err := NewDefault()
	assert.NoError(t, err)
	assert.Equal(t, &want, have)
}
