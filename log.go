// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import "log"

const panicNilLogger = "serv.WithLogger: Logger should not be nil"

type Logger interface {
	ServerStart(addr string)
	ServerShutdown()
	ServerClose()
}

type RouterLogger interface {
	RegisterRoute(name, method, path string)
}

func WithLogger(l Logger) Option {
	if l == nil {
		panic(panicNilLogger)
	}
	return optionFunc(func(s *Server) error {
		s.log = l
		return nil
	})
}

type DefaultLogger struct {
	Logger *log.Logger
}

func (l *DefaultLogger) log(v ...interface{}) {
	if l.Logger == nil {
		log.Println(v...)
	} else {
		l.Logger.Println(v...)
	}
}

func (l *DefaultLogger) ServerStart(addr string) { l.log("starting server on", addr) }
func (l *DefaultLogger) ServerShutdown()         { l.log("server shutdown") }
func (l *DefaultLogger) ServerClose()            { l.log("server close") }

func NopLogger() Logger { return new(nopLogger) }

type nopLogger struct{}

func (*nopLogger) ServerStart(string) {}
func (*nopLogger) ServerShutdown()    {}
func (*nopLogger) ServerClose()       {}
