// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"log"
	"strings"
)

const panicNilLogger = "serv.WithLogger: Logger should not be nil"

type Logger interface {
	ServerStart(name, addr string)
	ServerShutdown(name string)
	ServerClose(name string)
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

func (l *DefaultLogger) log(v ...string) {
	if l.Logger == nil {
		log.Println(strings.Join(v, " "))
	} else {
		l.Logger.Println(strings.Join(v, " "))
	}
}

func (l *DefaultLogger) ServerStart(name, addr string) { l.log("server", name, "starting on", addr) }
func (l *DefaultLogger) ServerShutdown(name string)    { l.log("server", name, "shutting down") }
func (l *DefaultLogger) ServerClose(name string)       { l.log("server", name, "closing") }

func NopLogger() Logger { return new(nopLogger) }

type nopLogger struct{}

func (*nopLogger) ServerStart(_, _ string) {}
func (*nopLogger) ServerShutdown(string)   {}
func (*nopLogger) ServerClose(string)      {}
