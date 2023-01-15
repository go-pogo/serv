// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpheader

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestBasicAuth_String(t *testing.T) {
	ba := BasicAuth{
		Username: "roeldev",
		Password: "S0m3SeCr37Pa55",
	}

	req := http.Request{Header: make(http.Header)}
	req.SetBasicAuth(ba.Username, ba.Password)
	assert.Equal(t, req.Header.Get(Authorization), ba.String())
}
