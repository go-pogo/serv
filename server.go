// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"context"
	"github.com/go-pogo/errors"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	ErrAlreadyStarted    errors.Msg = "server is already started"
	ErrUnstartedShutdown errors.Msg = "cannot shutdown server that is not started"
	ErrUnstartedClose    errors.Msg = "cannot close server that is not started"
)

type httpServer = http.Server

// Server is a wrapper for http.Server. The zero value is safe and ready to use,
// and will apply safe defaults on starting the server.
type Server struct {
	httpServer

	// Config to apply to the internal http.Server, DefaultConfig if zero.
	// Changes to Config after starting the server will not be applied.
	Config Config
	// Addr optionally specifies the TCP address for the server to listen on.
	// See net.Dial for details of the address format.
	// See http.Server for additional information.
	Addr string
	// Handler to invoke, http.DefaultServeMux if nil.
	Handler http.Handler

	log     Logger
	name    string
	started atomic.Bool
}

// New creates a new Server with a default Config.
func New(opts ...Option) (*Server, error) {
	srv := Server{Config: defaultConfig}
	if err := srv.With(opts...); err != nil {
		return nil, err
	}
	return &srv, nil
}

// With applies additional option to the server. It will return an
// ErrAlreadyStarted error when the server is already started.
func (srv *Server) With(opts ...Option) error {
	if srv.IsStarted() {
		return errors.New(ErrAlreadyStarted)
	}

	var err error
	for _, opt := range opts {
		err = errors.Append(err, opt.apply(srv))
	}
	return err
}

// Name returns an optional provided name of the server. Use WithName to set
// the server's name.
func (srv *Server) Name() string { return srv.name }

// IsStarted indicates if the server is started.
func (srv *Server) IsStarted() bool { return srv.started.Load() }

func (srv *Server) start() error {
	if srv.IsStarted() {
		return errors.New(ErrAlreadyStarted)
	}

	srv.started.Store(true)
	if srv.log == nil {
		srv.log = NopLogger()
	}
	if srv.Config.IsZero() {
		srv.Config = defaultConfig
	}

	handler := srv.Handler
	if srv.name != "" {
		if srv.Handler == nil {
			handler = http.DefaultServeMux
		}
		handler = WithServerName(srv.name, handler)
	}

	srv.Config.ApplyTo(&srv.httpServer)
	srv.httpServer.Addr = srv.Addr
	srv.httpServer.Handler = handler

	srv.log.ServerStart(srv.name, srv.Addr)
	return nil
}

// Serve starts the server by accepting incoming connections on the
// net.Listener l. See http.Server for additional information.
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

// ListenAndServe is a wrapper for http.Server.ListenAndServe.
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

// ServeTLS is a wrapper for http.Server.ServeTLS.
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

// ListenAndServeTLS is a wrapper for http.Server.ListenAndServeTLS.
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

// Run starts the server and calls either ListenAndServe or ListenAndServeTLS,
// depending on the provided TLS Option(s).
// Unlike Serve, ListenAndServe, ServeTLS, and ListenAndServeTLS, Run will not
// return a http.ErrServerClosed error when the server is closed.
func (srv *Server) Run() error {
	if srv.IsStarted() {
		return errors.New(ErrAlreadyStarted)
	}

	if srv.httpServer.TLSConfig != nil &&
		(len(srv.httpServer.TLSConfig.Certificates) != 0 ||
			srv.httpServer.TLSConfig.GetCertificate != nil) {
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
// connections. Just like the underlying http.Server, Shutdown works by first
// closing all open listeners, then closing all idle connections, and then
// waiting indefinitely for connections to return to idle and then shut down.
// If ShutdownTimeout is set and/or the provided context expires before the
// shutdown is complete, Shutdown returns the context's error. Otherwise, it
// returns any error returned from closing the Server's underlying
// net.Listener(s).
func (srv *Server) Shutdown(ctx context.Context) error {
	if !srv.IsStarted() {
		return errors.New(ErrUnstartedShutdown)
	}

	srv.log.ServerShutdown(srv.name)
	srv.httpServer.SetKeepAlivesEnabled(false)

	if srv.Config.ShutdownTimeout != 0 {
		if t, ok := ctx.Deadline(); !ok || srv.Config.ShutdownTimeout < time.Until(t) {
			// shutdown timeout is set to a lower value, update context
			var cancelFn context.CancelFunc
			ctx, cancelFn = context.WithTimeout(ctx, srv.Config.ShutdownTimeout)
			defer cancelFn()
		}
	}

	return errors.WithStack(srv.httpServer.Shutdown(ctx))
}

// Close immediately closes all active net.Listeners and any connections in
// state http.StateNew, http.StateActive, or http.StateIdle. For a graceful
// shutdown, use Shutdown.
// It is a wrapper for http.NewClient.Close.
func (srv *Server) Close() error {
	if !srv.IsStarted() {
		return errors.New(ErrUnstartedClose)
	}

	srv.log.ServerClose(srv.name)
	return srv.httpServer.Close()
}
