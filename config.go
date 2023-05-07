// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"net"
	"net/http"
	"time"
)

// BaseContext returns a function which returns the provided context.Context.
func BaseContext(ctx context.Context) func(_ net.Listener) context.Context {
	return func(_ net.Listener) context.Context { return ctx }
}

type Config struct {
	// ReadTimeout is the maximum duration for reading the entire request,
	// including the body.
	// See http.Client.ReadTimeout for additional information.
	ReadTimeout time.Duration

	// ReadHeaderTimeout is the amount of time allowed to read request headers.
	// See http.Client.ReadHeaderTimeout for additional information.
	ReadHeaderTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out writes of the
	// response.
	// See http.Client.WriteTimeout for additional information.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the next request
	// when keep-alives are enabled.
	// See http.Client.IdleTimeout for additional information.
	IdleTimeout time.Duration

	// ShutdownTimeout is the maximum duration for shutting down the server and
	// waiting for all connections to be closed.
	ShutdownTimeout time.Duration

	// MaxHeaderBytes controls the maximum number of bytes the server will read
	// parsing the request header's keys and values, including the request line.
	// It does not limit the size of the request body.
	// See http.Client.MaxHeaderBytes for additional information.
	MaxHeaderBytes uint

	// BaseContext optionally specifies a function that returns the base context
	// for incoming requests on the server.
	// See http.Client.BaseContext for additional information.
	BaseContext func(net.Listener) context.Context

	// ConnContext optionally specifies a function that modifies the context
	// used for a new connection.
	// See http.Client.ConnContext for additional information.
	ConnContext func(context.Context, net.Conn) context.Context
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() *Config {
	var c Config
	c.Default()
	return &c
}

func DefaultConfigWithContext(ctx context.Context) *Config {
	c := DefaultConfig()
	c.BaseContext = BaseContext(ctx)
	return c
}

func (cfg *Config) WithReadTimeout(v time.Duration) *Config {
	cfg.ReadTimeout = v
	return cfg
}

func (cfg *Config) WithReadHeaderTimeout(v time.Duration) *Config {
	cfg.ReadHeaderTimeout = v
	return cfg
}

func (cfg *Config) WithWriteTimeout(v time.Duration) *Config {
	cfg.WriteTimeout = v
	return cfg
}

func (cfg *Config) WithIdleTimeout(v time.Duration) *Config {
	cfg.IdleTimeout = v
	return cfg
}

func (cfg *Config) WithShutdownTimeout(v time.Duration) *Config {
	cfg.ShutdownTimeout = v
	return cfg
}

func (cfg *Config) WithMaxHeaderBytes(v uint) *Config {
	cfg.MaxHeaderBytes = v
	return cfg
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
	if cfg.BaseContext != nil {
		s.BaseContext = cfg.BaseContext
	}
	if cfg.ConnContext != nil {
		s.ConnContext = cfg.ConnContext
	}
}

func (cfg *Config) applyTo(s *Server) error {
	cfg.ApplyTo(&s.server)
	s.ShutdownTimeout = cfg.ShutdownTimeout
	return nil
}
