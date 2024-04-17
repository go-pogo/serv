// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"github.com/go-pogo/errors"
	"net"
	"net/http"
	"sync"
	"time"
)

type State uint32

const (
	StateUnstarted State = iota
	StateClosed
	StateClosing
	StateStarted

	ErrAlreadyStarted    errors.Msg = "server has already started"
	ErrAlreadyClosing    errors.Msg = "server is already closing"
	ErrUnstartedShutdown errors.Msg = "cannot shutdown server that is not started"
	ErrUnstartedClose    errors.Msg = "cannot close server that is not started"
)

type httpServer = http.Server

// Server is a wrapper for [http.Server]. The zero value is safe and ready to
// use, and will apply safe defaults on starting the server.
type Server struct {
	httpServer

	// Config to apply to the internal [http.Server], [DefaultConfig] if zero.
	// Changes to Config after starting the server will not be applied.
	Config Config
	// Addr optionally specifies the TCP address for the server to listen on.
	// Changing Addr after starting the server will not affect the server.
	// See [net.Dial] for details of the address format.
	// See [http.Server] for additional information.
	Addr string
	// Handler to invoke, [DefaultServeMux] if nil. Changing Handler after the
	// server has started will not have any effect.
	Handler http.Handler

	mut   sync.RWMutex
	log   Logger
	name  string
	state State
}

// New creates a new [Server] with a default [Config].
func New(opts ...Option) (*Server, error) {
	srv := Server{Config: defaultConfig}
	if err := srv.with(opts); err != nil {
		return nil, err
	}
	return &srv, nil
}

// With applies additional option(s) to the server. It will return an
// [InvalidStateError] containing a [ErrAlreadyStarted] error when the
// server has already started.
func (srv *Server) With(opts ...Option) error {
	if srv.State() == StateStarted {
		return errors.New(ErrAlreadyStarted)
	}

	srv.mut.Lock()
	defer srv.mut.Unlock()
	return srv.with(opts)
}

func (srv *Server) with(opts []Option) error {
	var err error
	for _, opt := range opts {
		err = errors.Append(err, opt.apply(srv))
	}
	return err
}

// Name returns an optional provided name of the server. Use [WithName] to set
// the server's name.
func (srv *Server) Name() string {
	srv.mut.RLock()
	defer srv.mut.RUnlock()
	return srv.name
}

// State returns the current [State] of the [Server].
func (srv *Server) State() State {
	srv.mut.RLock()
	defer srv.mut.RUnlock()
	return srv.state
}

func (srv *Server) start() error {
	srv.mut.Lock()
	defer srv.mut.Unlock()

	if srv.state == StateStarted {
		return errors.New(ErrAlreadyStarted)
	}
	if srv.state == StateClosing {
		return errors.New(ErrAlreadyClosing)
	}
	if srv.state == StateClosed {
		srv.httpServer = http.Server{
			DisableGeneralOptionsHandler: srv.httpServer.DisableGeneralOptionsHandler,
			TLSConfig:                    srv.httpServer.TLSConfig,
			TLSNextProto:                 srv.httpServer.TLSNextProto,
			ConnState:                    srv.httpServer.ConnState,
			ErrorLog:                     srv.httpServer.ErrorLog,
			BaseContext:                  srv.httpServer.BaseContext,
			ConnContext:                  srv.httpServer.ConnContext,
		}
	}
	if srv.log == nil {
		srv.log = NopLogger()
	}
	if srv.Config.IsZero() {
		srv.Config = defaultConfig
	}

	handler := srv.Handler
	if srv.Handler == nil {
		handler = DefaultServeMux()
	}
	if srv.name != "" {
		handler = AddServerName(srv.name, handler)
	}

	srv.Config.ApplyTo(&srv.httpServer)
	srv.httpServer.Addr = srv.Addr
	srv.httpServer.Handler = handler

	srv.state = StateStarted
	srv.log.ServerStart(srv.name, srv.Addr)
	return nil
}

// Serve is a wrapper for [http.Server.Serve].
func (srv *Server) Serve(l net.Listener) error {
	if err := srv.start(); err != nil {
		return err
	}

	err := srv.httpServer.Serve(l)
	if !errors.Is(err, http.ErrServerClosed) {
		err = errors.WithStack(err)
	}
	return err
}

// ListenAndServe is a wrapper for [http.Server.ListenAndServe].
func (srv *Server) ListenAndServe() error {
	if err := srv.start(); err != nil {
		return err
	}

	err := srv.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		err = errors.WithStack(err)
	}
	return err
}

// ServeTLS is a wrapper for [http.Server.ServeTLS].
func (srv *Server) ServeTLS(l net.Listener, certFile, keyFile string) error {
	if err := srv.start(); err != nil {
		return err
	}

	err := srv.httpServer.ServeTLS(l, certFile, keyFile)
	if !errors.Is(err, http.ErrServerClosed) {
		err = errors.WithStack(err)
	}
	return err
}

// ListenAndServeTLS is a wrapper for [http.Server.ListenAndServeTLS].
func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	if err := srv.start(); err != nil {
		return err
	}

	err := srv.httpServer.ListenAndServeTLS(certFile, keyFile)
	if !errors.Is(err, http.ErrServerClosed) {
		err = errors.WithStack(err)
	}
	return err
}

// Run starts the server and calls either [Server.ListenAndServe] or
// [Server.ListenAndServeTLS], depending on the provided TLS config/option(s).
// Unlike [Server.Serve], [Server.ListenAndServe], [Server.ServeTLS], and
// [Server.ListenAndServeTLS], Run will not return a [http.ErrServerClosed]
// error when the server is closed.
func (srv *Server) Run() error {
	srv.mut.RLock()
	useTLS := srv.httpServer.TLSConfig != nil &&
		(len(srv.httpServer.TLSConfig.Certificates) != 0 ||
			srv.httpServer.TLSConfig.GetCertificate != nil)
	srv.mut.RUnlock()

	if useTLS {
		return dismissErrServerClosed(srv.ListenAndServeTLS("", ""))
	}
	return dismissErrServerClosed(srv.ListenAndServe())
}

func dismissErrServerClosed(err error) error {
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// Shutdown gracefully shuts down the server without interrupting any active
// connections. Just like the underlying [http.Server], Shutdown works by first
// closing all open listeners, then closing all idle connections, and then
// waiting indefinitely for connections to return to idle and then shut down.
// If [Config.ShutdownTimeout] is set and/or the provided context expires before
// the shutdown is complete, Shutdown returns the context's error. Otherwise, it
// returns any error returned from closing the [Server]'s underlying
// [net.Listener](s).
// An [InvalidStateError] containing a [ErrUnableToShutdown] error is returned
// when the server is not started.
func (srv *Server) Shutdown(ctx context.Context) error {
	if state := srv.State(); state == StateClosing {
		return errors.New(ErrAlreadyClosing)
	} else if state != StateStarted {
		return errors.New(ErrUnstartedShutdown)
	}

	srv.mut.Lock()
	srv.state = StateClosing
	srv.log.ServerShutdown(srv.name)
	srv.httpServer.SetKeepAlivesEnabled(false)
	shutdownTimeout := srv.Config.ShutdownTimeout
	srv.mut.Unlock()

	if shutdownTimeout != 0 {
		if t, ok := ctx.Deadline(); !ok || shutdownTimeout < time.Until(t) {
			// shutdown timeout is set to a lower value, update context
			var cancelFn context.CancelFunc
			ctx, cancelFn = context.WithTimeout(ctx, shutdownTimeout)
			defer cancelFn()
		}
	}

	defer srv.closed()
	return errors.WithStack(srv.httpServer.Shutdown(ctx))
}

// Close immediately closes all active [net.Listener](s) and any connections in
// state [http.StateNew], [http.StateActive], or [http.StateIdle].
// An [InvalidStateError] containing a [ErrUnableToClose] error is returned
// when the server is not started.
// For a graceful shutdown, use [Server.Shutdown].
func (srv *Server) Close() error {
	if state := srv.State(); state == StateClosing {
		return errors.New(ErrAlreadyClosing)
	} else if state != StateStarted {
		return errors.New(ErrUnstartedClose)
	}

	srv.mut.Lock()
	srv.state = StateClosing
	srv.log.ServerClose(srv.name)
	srv.mut.Unlock()

	defer srv.closed()
	return errors.WithStack(srv.httpServer.Close())
}

func (srv *Server) closed() {
	srv.mut.Lock()
	srv.state = StateClosed
	srv.mut.Unlock()
}
