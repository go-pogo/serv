// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"github.com/go-pogo/rawconv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestConfig_Default(t *testing.T) {
	prepare := func() (want Config, err error) {
		val := reflect.ValueOf(&want).Elem()
		typ := val.Type()

		var u rawconv.Unmarshaler
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			v := field.Tag.Get("default")
			if v == "" {
				continue
			}

			if err = u.Unmarshal(rawconv.Value(v), val.Field(i)); err != nil {
				break
			}
		}
		return
	}

	want, err := prepare()
	require.NoError(t, err, "failed to prepare expected value")
	assert.Equal(t, want, *DefaultConfig(), "default tag should match DefaultConfig() values")
}

func TestConfig_ApplyTo(t *testing.T) {
	var have http.Server
	DefaultConfig().ApplyTo(&have)

	assert.Equal(t, http.Server{
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    10240,
	}, have)
}
