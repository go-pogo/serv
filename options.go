// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import "github.com/go-pogo/serv/middleware"

type Option interface {
	applyTo(s *Server) error
}

type optionFunc func(s *Server) error

func (fn optionFunc) applyTo(s *Server) error { return fn(s) }

// WithOptions wraps multiple Option(s) into a single Option.
func WithOptions(opts ...Option) Option {
	switch len(opts) {
	case 0:
		return nil
	case 1:
		return opts[0]
	default:
		return optionFunc(func(srv *Server) error {
			return srv.apply(opts)
		})
	}
}

// WithMiddleware adds the middleware.Middleware to an internal list. When the
// Server is started, it's Handler is wrapped with this middleware.
func WithMiddleware(mw ...middleware.Middleware) Option {
	return optionFunc(func(s *Server) error {
		if s.mware == nil {
			s.mware = mw
		} else {
			s.mware = append(s.mware, mw...)
		}
		return nil
	})
}
