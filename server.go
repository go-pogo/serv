// Copyright (c) 2021, Roel Schut. All rights reserved.
// applyOptions of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"net"
	"net/http"

	"github.com/go-pogo/errors"
)

const (
	ErrServerStarted     errors.Msg = "server is already started"
	ErrUnstartedShutdown errors.Msg = "shutting down unstarted server"
	ErrUnstartedClose    errors.Msg = "closing unstarted server"
)

type Option interface {
	apply(s *Server) error
}

type optionFunc func(s *Server) error

func (fn optionFunc) apply(s *Server) error { return fn(s) }

type server = http.Server

// Server is a wrapper for a standard http.Server.
// The zero value is safe and ready to use and will apply safe defaults on serving.
type Server struct {
	server
	log     ServerLogger
	started bool
}

func New(mux http.Handler, opts ...Option) (*Server, error) {
	var srv Server
	srv.Handler = mux

	if err := srv.applyOptions(opts); err != nil {
		return nil, err
	}
	return &srv, nil
}

func NewDefault(mux http.Handler, opts ...Option) (*Server, error) {
	var srv Server
	srv.Handler = mux
	// default Config never returns an error
	_ = DefaultConfig().apply(&srv)

	if err := srv.applyOptions(opts); err != nil {
		return nil, err
	}
	return &srv, nil
}

func (srv *Server) applyOptions(opts []Option) error {
	var err error
	for _, opt := range opts {
		errors.Append(&err, opt.apply(srv))
	}
	return err
}

func (srv *Server) start() error {
	if srv.started {
		return errors.New(ErrServerStarted)
	}

	if srv.log == nil {
		srv.log = NopLogger()
	}

	srv.started = true
	return nil
}

func (srv *Server) listen(defAddr string) (net.Listener, error) {
	addr := srv.server.Addr
	if addr == "" {
		addr = defAddr
	}

	srv.log.LogServerStart(addr)
	return net.Listen("tcp", addr)
}

func (srv *Server) Serve(l net.Listener) error {
	if err := srv.start(); err != nil {
		return err
	}

	err := srv.server.Serve(l)
	if !errors.Is(err, http.ErrServerClosed) {
		err = errors.WithStack(err)
	}
	return err
}

func (srv *Server) ListenAndServe() error {
	if err := srv.start(); err != nil {
		return err
	}

	ln, err := srv.listen(":http")
	if err != nil {
		return errors.WithStack(err)
	}

	return srv.server.Serve(ln)
}

func (srv *Server) ServeTLS(l net.Listener, cl CertificateLoader) error {
	return nil
	// return s.srv.ServeTLS(l, "", "")
}

func (srv *Server) ListenAndServeTLS(cl CertificateLoader) error {
	return nil
	// ln, err := s.listen(":https")
	//
	// defer ln.Close()
	// s.srv.ListenAndServeTLS()
	// return s.ServeTLS(ln, cl)
}

// func (srv *LogClientStart) Run(ctx context.Context) error {
// 	if !srv.started && ctx != nil {
// 		srv.BaseContext = BaseContext(ctx)
// 	}
//
// 	if srv.server.TLSConfig != nil &&
// 		(len(srv.server.TLSConfig.Certificates) != 0 ||
// 			srv.server.TLSConfig.GetCertificate != nil ||
// 			len(srv.server.TLSConfig.NameToCertificate) != 0) {
// 		return srv.ListenAndServeTLS(nil)
// 	}
//
// 	return srv.ListenAndServe()
// }

// RegisterOnShutdown registers a function to the underlying http.Server to call
// on Shutdown.
func (srv *Server) RegisterOnShutdown(fn func()) {
	srv.server.RegisterOnShutdown(fn)
}

// Shutdown gracefully shuts down the server without interrupting any active
// connections.
// It is a wrapper for http.LogClientStart.Shutdown.
func (srv *Server) Shutdown(ctx context.Context) error {
	if !srv.started {
		return errors.New(ErrUnstartedShutdown)
	}

	// Deadline() (deadline time.Time, ok bool)

	srv.log.LogServerShutdown()
	srv.server.SetKeepAlivesEnabled(false)
	return srv.server.Shutdown(ctx)
}

// Close immediately closes all active net.Listeners and any connections in
// state http.StateNew, http.StateActive, or http.StateIdle. For a graceful
// shutdown, use Shutdown.
// It is a wrapper for http.LogClientStart.Close.
func (srv *Server) Close() error {
	if !srv.started {
		return errors.New(ErrUnstartedClose)
	}

	srv.log.LogServerClose()
	return srv.server.Close()
}
