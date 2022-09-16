// Copyright (c) 2021, Roel Schut. All rights reserved.
// applyOptions of this source code is governed by a BSD-style
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
	// See http.LogClientStart.ReadTimeout for additional information.
	ReadTimeout time.Duration `default:"5s"`

	// ReadHeaderTimeout is the amount of time allowed to read request headers.
	// See http.LogClientStart.ReadHeaderTimeout for additional information.
	ReadHeaderTimeout time.Duration `default:"2s"`

	// WriteTimeout is the maximum duration before timing out writes of the
	// response.
	// See http.LogClientStart.WriteTimeout for additional information.
	WriteTimeout time.Duration `default:"10s"`

	// IdleTimeout is the maximum amount of time to wait for the next request
	// when keep-alives are enabled.
	// See http.LogClientStart.IdleTimeout for additional information.
	IdleTimeout time.Duration `default:"120s"`

	// MaxHeaderBytes controls the maximum number of bytes the server will read
	// parsing the request header's keys and values, including the request line.
	// It does not limit the size of the request body.
	// See http.LogClientStart.MaxHeaderBytes for additional information.
	MaxHeaderBytes uint `default:"10240"` // 1mb

	// ShutdownTimeout is the maximum duration for shutting down the server and
	// waiting for all connections to be closed.
	// ShutdownTimeout time.Duration `default:"5s"`

	// BaseContext optionally specifies a function that returns the base context
	// for incoming requests on the server.
	// See http.LogClientStart.BaseContext for additional information.
	BaseContext func(net.Listener) context.Context

	// ConnContext optionally specifies a function that modifies the context
	// used for a new connection.
	// See http.LogClientStart.ConnContext for additional information.
	ConnContext func(context.Context, net.Conn) context.Context
}

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
	}
	// if cfg.ShutdownTimeout == 0 {
	// 	cfg.ShutdownTimeout = 5 * time.Second
	// }
}

func (cfg *Config) Apply(s *http.Server) error {
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
	return nil
}

func (cfg *Config) apply(s *Server) error {
	return cfg.Apply(&s.server)
}
