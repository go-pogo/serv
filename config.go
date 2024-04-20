// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"net/http"
	"time"
)

var _ Option = (*Config)(nil)

type Config struct {
	// ReadTimeout is the maximum duration for reading the entire request,
	// including the body.
	// See [http.Server.ReadTimeout] for additional information.
	ReadTimeout time.Duration `default:"5s"`
	// ReadHeaderTimeout is the amount of time allowed to read request headers.
	// See [http.Server.ReadHeaderTimeout] for additional information.
	ReadHeaderTimeout time.Duration `default:"2s"`
	// WriteTimeout is the maximum duration before timing out writes of the
	// response.
	// See [http.Server.WriteTimeout] for additional information.
	WriteTimeout time.Duration `default:"10s"`
	// IdleTimeout is the maximum amount of time to wait for the next request
	// when keep-alives are enabled.
	// See [http.Server.IdleTimeout] for additional information.
	IdleTimeout time.Duration `default:"120s"`
	// ShutdownTimeout is the default maximum duration for shutting down the
	// [Server] and waiting for all connections to be closed.
	ShutdownTimeout time.Duration `default:"60s"`
	// MaxHeaderBytes controls the maximum number of bytes the server will read
	// parsing the [http.Request] header's keys and values, including the
	// request line. It does not limit the size of the request body.
	// See [http.Server.MaxHeaderBytes] for additional information.
	MaxHeaderBytes uint64 `default:"10240"` // data.Bytes => 10 KiB
}

var defaultConfig = Config{
	ReadTimeout:       5 * time.Second,
	ReadHeaderTimeout: 2 * time.Second,
	WriteTimeout:      10 * time.Second,
	IdleTimeout:       120 * time.Second,
	ShutdownTimeout:   60 * time.Second,
	MaxHeaderBytes:    10240, // 10 KiB => 10 * data.Kibibyte
}

// DefaultConfig returns a [Config] with safe default values.
func DefaultConfig() *Config {
	c := defaultConfig
	return &c
}

// IsZero indicates [Config] equals its zero value.
func (cfg *Config) IsZero() bool { return *cfg == Config{} }

// Default sets any zero values on [Config] to a default non-zero value similar
// to [DefaultConfig].
func (cfg *Config) Default() {
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = defaultConfig.ReadTimeout
	}
	if cfg.ReadHeaderTimeout == 0 {
		cfg.ReadHeaderTimeout = defaultConfig.ReadHeaderTimeout
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = defaultConfig.WriteTimeout
	}
	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = defaultConfig.IdleTimeout
	}
	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = defaultConfig.ShutdownTimeout
	}
	if cfg.MaxHeaderBytes == 0 {
		cfg.MaxHeaderBytes = defaultConfig.MaxHeaderBytes
	}
}

// ApplyTo applies the [Config] fields' values to [http.Server] s.
func (cfg *Config) ApplyTo(s *http.Server) {
	if s == nil {
		return
	}

	if cfg.ReadTimeout != 0 {
		s.ReadTimeout = cfg.ReadTimeout
	}
	if cfg.ReadHeaderTimeout != 0 {
		s.ReadHeaderTimeout = cfg.ReadHeaderTimeout
	}
	if cfg.WriteTimeout != 0 {
		s.WriteTimeout = cfg.WriteTimeout
	}
	if cfg.IdleTimeout != 0 {
		s.IdleTimeout = cfg.IdleTimeout
	}
	if cfg.MaxHeaderBytes != 0 {
		s.MaxHeaderBytes = int(cfg.MaxHeaderBytes)
	}
}

func (cfg *Config) apply(srv *Server) error {
	srv.Config = *cfg
	return nil
}
