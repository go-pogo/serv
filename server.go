// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/serv/middleware"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	ErrServerStarted     errors.Msg = "server is already started"
	ErrUnstartedShutdown errors.Msg = "shutting down unstarted server"
	ErrUnstartedClose    errors.Msg = "closing unstarted server"
)

type server = http.Server

// Server is a wrapper for a standard http.Server.
// The zero value is safe and ready to use, and will apply safe defaults on
// starting the server.
type Server struct {
	server
	log     Logger
	mware   []middleware.Middleware
	name    string
	started uint32

	ShutdownTimeout time.Duration
}

// New creates a new Server.
func New(mux http.Handler, opts ...Option) (*Server, error) {
	var srv Server
	srv.Handler = mux

	if err := srv.apply(opts); err != nil {
		return nil, err
	}
	return &srv, nil
}

// NewDefault creates a new Server with DefaultConfig applied to it.
func NewDefault(mux http.Handler, opts ...Option) (*Server, error) {
	return New(mux, DefaultConfig(), WithOptions(opts...))
}

func (srv *Server) apply(opts []Option) error {
	var err error
	for _, opt := range opts {
		errors.Append(&err, opt.applyTo(srv))
	}
	return err
}

func (srv *Server) isStarted() bool {
	return atomic.LoadUint32(&srv.started) == 1
}

func (srv *Server) start() error {
	if srv.isStarted() {
		return errors.New(ErrServerStarted)
	}

	if srv.log == nil {
		srv.log = NopLogger()
	}
	if len(srv.mware) != 0 {
		srv.Handler = middleware.Wrap(srv.Handler, srv.mware...)
		srv.mware = nil
	}

	srv.log.ServerStart(srv.name, srv.Addr)
	atomic.StoreUint32(&srv.started, 1)
	return nil
}

func (srv *Server) Name() string { return srv.name }

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

	err := srv.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		err = errors.WithStack(err)
	}
	return err
}

func (srv *Server) ServeTLS(l net.Listener, certFile, keyFile string) error {
	if err := srv.start(); err != nil {
		return err
	}

	err := srv.server.ServeTLS(l, certFile, keyFile)
	if !errors.Is(err, http.ErrServerClosed) {
		err = errors.WithStack(err)
	}
	return err
}

func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	if err := srv.start(); err != nil {
		return err
	}

	err := srv.server.ListenAndServeTLS(certFile, keyFile)
	if !errors.Is(err, http.ErrServerClosed) {
		err = errors.WithStack(err)
	}
	return err
}

func (srv *Server) Run(ctx context.Context) error {
	if srv.isStarted() {
		return errors.New(ErrServerStarted)
	}
	if ctx != nil {
		srv.BaseContext = BaseContext(ctx)
	}

	if srv.server.TLSConfig != nil &&
		(len(srv.server.TLSConfig.Certificates) != 0 ||
			srv.server.TLSConfig.GetCertificate != nil) {
		return srv.ListenAndServeTLS("", "")
	}

	return srv.ListenAndServe()
}

// RegisterOnShutdown registers a function to the underlying http.Server, which
// is called on Shutdown.
func (srv *Server) RegisterOnShutdown(fn func()) {
	srv.server.RegisterOnShutdown(fn)
}

// Shutdown gracefully shuts down the server without interrupting any active
// connections. Just like the underlying http.Server, Shutdown works by first
// closing all open listeners, then closing all idle connections, and then
// waiting indefinitely for connections to return to idle and then shut down.
// If ShutdownTimeout is set and/or the provided context expires before the
// shutdown is complete, Shutdown returns the context's error. Otherwise it
// returns any error returned from closing the Server's underlying
// net.Listener(s).
func (srv *Server) Shutdown(ctx context.Context) error {
	if !srv.isStarted() {
		return errors.New(ErrUnstartedShutdown)
	}

	srv.log.ServerShutdown(srv.name)
	srv.server.SetKeepAlivesEnabled(false)

	if srv.ShutdownTimeout != 0 {
		if t, ok := ctx.Deadline(); !ok || srv.ShutdownTimeout < time.Until(t) {
			// shutdown timeout is set to a lower value, update context
			var cancelFn context.CancelFunc
			ctx, cancelFn = context.WithTimeout(ctx, srv.ShutdownTimeout)
			defer cancelFn()
		}
	}

	return errors.WithStack(srv.server.Shutdown(ctx))
}

// Close immediately closes all active net.Listeners and any connections in
// state http.StateNew, http.StateActive, or http.StateIdle. For a graceful
// shutdown, use Shutdown.
// It is a wrapper for http.NewClient.Close.
func (srv *Server) Close() error {
	if !srv.isStarted() {
		return errors.New(ErrUnstartedClose)
	}

	srv.log.ServerClose(srv.name)
	return srv.server.Close()
}
