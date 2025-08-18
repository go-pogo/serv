// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-pogo/errors"
)

const (
	ErrApplyOptions     errors.Msg = "error while applying option(s)"
	ErrAlreadyStarted   errors.Msg = "server has already started"
	ErrUnableToStart    errors.Msg = "unable to start server"
	ErrUnableToShutdown errors.Msg = "unable to shutdown server"
	ErrUnableToClose    errors.Msg = "unable to close server"
	ErrServerShutdown   errors.Msg = "error during server shutdown"
	ErrServerClose      errors.Msg = "error during closing server"
)

type httpServer = http.Server

// Server is a wrapper for [http.Server]. The zero value is safe and ready to
// use, and will apply safe defaults on starting the [Server].
type Server struct {
	httpServer

	// Config to apply to the internal [http.Server], [DefaultConfig] if zero.
	// Changes to [Config] after starting the [Server] will not be applied
	// until after the [Server] is restarted.
	Config Config
	// Addr optionally specifies the TCP address for the server to listen on.
	// Changing Addr after starting the [Server] will not affect it until after
	// the [Server] is restarted.
	// See [net.Dial] for details of the address format.
	// See [http.Server] for additional information.
	Addr string
	// Handler to invoke, [DefaultServeMux] if nil. Changing Handler after the
	// [Server] has started will not have any effect until after the [Server]
	// is restarted.
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
		return nil, errors.Wrap(err, ErrApplyOptions)
	}
	return &srv, nil
}

// With applies additional [Option](s) to the server. It will return an
// [InvalidStateError] containing a [ErrAlreadyStarted] error when the
// server has already started.
func (srv *Server) With(opts ...Option) error {
	if state := srv.State(); state == StateStarted {
		return errors.WithStack(&InvalidStateError{
			Err:   ErrAlreadyStarted,
			State: state,
		})
	}

	srv.mut.Lock()
	defer srv.mut.Unlock()

	if err := srv.with(opts); err != nil {
		return errors.Wrap(err, ErrApplyOptions)
	}
	return nil
}

func (srv *Server) with(opts []Option) error {
	var err error
	for _, opt := range opts {
		if opt == nil {
			continue
		}
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

	if srv.state == StateStarted || srv.state == StateClosing {
		return errors.WithStack(&InvalidStateError{
			Err:   ErrUnableToStart,
			State: srv.state,
		})
	}
	if srv.state == StateClosed {
		srv.httpServer = http.Server{
			DisableGeneralOptionsHandler: srv.DisableGeneralOptionsHandler,
			TLSConfig:                    srv.TLSConfig,
			TLSNextProto:                 srv.TLSNextProto,
			ConnState:                    srv.ConnState,
			ErrorLog:                     srv.ErrorLog,
			BaseContext:                  srv.BaseContext,
			ConnContext:                  srv.ConnContext,
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
	return nil
}

// Serve is a wrapper for [http.Server.Serve].
func (srv *Server) Serve(l net.Listener) error {
	if err := srv.start(); err != nil {
		return err
	}

	srv.log.LogServerStart(srv.name, srv.Addr)
	err := srv.httpServer.Serve(l)
	if !srv.isClosed(err) {
		err = errors.WithStack(err)
	}
	return err
}

// ListenAndServe is a wrapper for [http.Server.ListenAndServe].
func (srv *Server) ListenAndServe() error {
	if err := srv.start(); err != nil {
		return err
	}

	srv.log.LogServerStart(srv.name, srv.Addr)
	err := srv.httpServer.ListenAndServe()
	if !srv.isClosed(err) {
		err = errors.WithStack(err)
	}
	return err
}

// ServeTLS is a wrapper for [http.Server.ServeTLS].
func (srv *Server) ServeTLS(l net.Listener, certFile, keyFile string) error {
	if err := srv.start(); err != nil {
		return err
	}

	srv.log.LogServerStartTLS(srv.name, srv.Addr, certFile, keyFile)
	err := srv.httpServer.ServeTLS(l, certFile, keyFile)
	if !srv.isClosed(err) {
		err = errors.WithStack(err)
	}
	return err
}

// ListenAndServeTLS is a wrapper for [http.Server.ListenAndServeTLS].
func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	if err := srv.start(); err != nil {
		return err
	}

	srv.log.LogServerStartTLS(srv.name, srv.Addr, certFile, keyFile)
	err := srv.httpServer.ListenAndServeTLS(certFile, keyFile)
	if !srv.isClosed(err) {
		err = errors.WithStack(err)
	}
	return err
}

// isClosed checks if the provided error is http.ErrServerClosed, which
// indicates the server has been successfully closed. In this case, state is
// set to StateClosed and isClosed returns true.
// If an error occurs while starting the internal server, state is set to
// StateErrored and isClosed returns false.
func (srv *Server) isClosed(err error) (ok bool) {
	state := StateErrored
	if ok = errors.Is(err, http.ErrServerClosed); ok {
		state = StateClosed
	}

	srv.mut.Lock()
	srv.state = state
	srv.mut.Unlock()
	return ok
}

// Run starts the server and calls either [Server.ListenAndServe] or
// [Server.ListenAndServeTLS], depending on the provided TLS config/option(s).
// Unlike [Server.Serve], [Server.ListenAndServe], [Server.ServeTLS], and
// [Server.ListenAndServeTLS], Run will not return a [http.ErrServerClosed]
// error when the server is closed.
func (srv *Server) Run() error {
	srv.mut.RLock()
	useTLS := srv.TLSConfig != nil &&
		(len(srv.TLSConfig.Certificates) != 0 ||
			srv.TLSConfig.GetCertificate != nil)
	srv.mut.RUnlock()

	var err error
	if useTLS {
		err = srv.ListenAndServeTLS("", "")
	} else {
		err = srv.ListenAndServe()
	}
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
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
	if state := srv.State(); state != StateStarted {
		return errors.WithStack(&InvalidStateError{
			Err:   ErrUnableToShutdown,
			State: state,
		})
	}

	srv.mut.Lock()
	srv.state = StateClosing
	srv.log.LogServerShutdown(srv.name)
	srv.SetKeepAlivesEnabled(false)
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

	defer srv.close()
	return errors.Wrap(srv.httpServer.Shutdown(ctx), ErrServerShutdown)
}

// Close immediately closes all active [net.Listener](s) and any connections in
// state [http.StateNew], [http.StateActive], or [http.StateIdle].
// An [InvalidStateError] containing a [ErrUnableToClose] error is returned
// when the server is not started.
// For a graceful shutdown, use [Server.Shutdown].
func (srv *Server) Close() error {
	if state := srv.State(); state != StateStarted {
		return errors.WithStack(&InvalidStateError{
			Err:   ErrUnableToClose,
			State: state,
		})
	}

	srv.mut.Lock()
	srv.state = StateClosing
	srv.log.LogServerClose(srv.name)
	srv.mut.Unlock()

	defer srv.close()
	return errors.Wrap(srv.httpServer.Close(), ErrServerClose)
}

func (srv *Server) close() {
	srv.mut.Lock()
	srv.state = StateClosed
	srv.mut.Unlock()
}
