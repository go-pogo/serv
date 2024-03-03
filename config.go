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
	// See http.Server for additional information.
	ReadTimeout time.Duration `default:"5s"`

	// ReadHeaderTimeout is the amount of time allowed to read request headers.
	// See http.Server for additional information.
	ReadHeaderTimeout time.Duration `default:"2s"`

	// WriteTimeout is the maximum duration before timing out writes of the
	// response.
	// See http.Server for additional information.
	WriteTimeout time.Duration `default:"10s"`

	// IdleTimeout is the maximum amount of time to wait for the next request
	// when keep-alives are enabled.
	// See http.Server for additional information.
	IdleTimeout time.Duration `default:"120s"`

	// ShutdownTimeout is the maximum duration for shutting down the server and
	// waiting for all connections to be closed.
	ShutdownTimeout time.Duration `default:"60s"`

	// MaxHeaderBytes controls the maximum number of bytes the server will read
	// parsing the request header's keys and values, including the request line.
	// It does not limit the size of the request body.
	// See http.Server for additional information.
	MaxHeaderBytes uint64 `default:"10 KiB"` // data.Bytes
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() *Config {
	var c Config
	c.Default()
	return &c
}

// Default sets any zero values on Config to a default non-zero value.
func (cfg *Config) Default() {
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 5 * time.Second
	}
	if cfg.ReadHeaderTimeout == 0 {
		cfg.ReadHeaderTimeout = 2 * time.Second
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 10 * time.Second
	}
	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = 120 * time.Second
	}
	if cfg.MaxHeaderBytes == 0 {
		cfg.MaxHeaderBytes = 10240
		//cfg.MaxHeaderBytes = 10 * data.Kibibyte
	}
	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = 60 * time.Second
	}
}

// ApplyTo applies the Config fields values to *http.Server s.
func (cfg *Config) ApplyTo(s *http.Server) {
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

func (cfg *Config) apply(s *Server) error {
	s.Config = cfg
	return nil
}
