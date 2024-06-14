// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/go-pogo/rawconv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var cfg Config
		assert.True(t, cfg.IsZero())
	})
	t.Run("non-zero", func(t *testing.T) {
		cfg := Config{ReadTimeout: 2 * time.Second}
		assert.False(t, cfg.IsZero())
	})
}

func TestDefaultConfig(t *testing.T) {
	t.Run("copies defaultConfig", func(t *testing.T) {
		cfg := DefaultConfig()
		assert.Equal(t, defaultConfig, *cfg)
		cfg.IdleTimeout *= 2
		assert.NotEqual(t, defaultConfig, *cfg)
	})
	t.Run("matches tags", func(t *testing.T) {
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
	})
}

func TestConfig_Default(t *testing.T) {
	var cfg Config
	cfg.Default()
	assert.Equal(t, defaultConfig, cfg)
}

func TestConfig_ApplyTo(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var server *http.Server
		DefaultConfig().ApplyTo(server)
		assert.Nil(t, server)
	})
	t.Run("non-nil", func(t *testing.T) {
		var have http.Server
		DefaultConfig().ApplyTo(&have)

		assert.Equal(t, &http.Server{
			ReadTimeout:       defaultConfig.ReadTimeout,
			ReadHeaderTimeout: defaultConfig.ReadHeaderTimeout,
			WriteTimeout:      defaultConfig.WriteTimeout,
			IdleTimeout:       defaultConfig.IdleTimeout,
			MaxHeaderBytes:    int(defaultConfig.MaxHeaderBytes),
		}, &have)
	})
}
