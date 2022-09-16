// Copyright (c) 2021, Roel Schut. All rights reserved.
// applyOptions of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

const panicNilLogger = "Logger should not be nil"

type ServerLogger interface {
	LogServerStart(addr string)
	LogServerShutdown()
	LogServerClose()
}

func WithLogger(l ServerLogger) Option {
	if l == nil {
		panic(panicNilLogger)
	}
	return optionFunc(func(s *Server) error {
		s.log = l
		return nil
	})
}

func NopLogger() ServerLogger { return new(nopLogger) }

type nopLogger struct{}

func (_ *nopLogger) LogServerStart(_ string) {}
func (_ *nopLogger) LogServerShutdown()      {}
func (_ *nopLogger) LogServerClose()         {}
