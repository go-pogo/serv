// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	"github.com/go-pogo/errors"
	"github.com/go-pogo/serv"
	"github.com/go-pogo/serv/accesslog"
	"github.com/go-pogo/serv/middleware"
	"github.com/go-pogo/serv/response"
)

// This program restarts the server after the "/restart" url is visited.

func main() {
	var port serv.Port = 8080

	cli := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cli.Var(&port, "port", "Server port")
	_ = cli.Parse(os.Args[1:])

	var restartCount atomic.Uint64
	chRestart := make(chan struct{}, 1)
	chRestart <- struct{}{}

	mux := serv.NewServeMux()
	mux.HandleRoute(serv.Route{
		Name:    "count",
		Pattern: "/",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			_, _ = fmt.Fprintf(w,
				`<b>%d</b><br /><a href="/restart">restart</a>`,
				restartCount.Load(),
			)
		}),
	})
	mux.HandleRoute(serv.Route{
		Name:    "favicon",
		Method:  http.MethodGet,
		Pattern: "/favicon.ico",
		Handler: accesslog.IgnoreHandler(response.NoContentHandler()),
	})
	mux.HandleRoute(serv.Route{
		Name:    "restart",
		Pattern: "/restart",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(`<i>Restarting server...</i><br /><a href="/">show count</a>`))

			go func() {
				chRestart <- struct{}{}
			}()
		}),
	})

	ctx, stopFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopFn()

	srv, err := serv.New(port,
		serv.WithBaseContext(ctx),
		serv.WithDefaultLogger(),
		serv.WithHandler(middleware.Wrap(mux,
			accesslog.Middleware(accesslog.DefaultLogger()),
		)),
	)
	errors.FatalOnErr(err)

loop:
	for {
		select {
		case <-chRestart:
			if srv.State() == serv.StateStarted {
				log.Println("restarting server...")
				if err = srv.Shutdown(context.Background()); err != nil {
					if !errors.Is(err, context.DeadlineExceeded) {
						log.Printf("shutdown error while restarting: %+v\n", err)
					}
				}
			}
			if state := srv.State(); state != serv.StateStarted {
				if state == serv.StateClosed {
					restartCount.Add(1)
				}
				go func() {
					if err := srv.Run(); err != nil {
						log.Println("server error:", err.Error())
					}
				}()
			}

		case <-ctx.Done():
			break loop
		}
	}

	if srv.State() == serv.StateStarted {
		if err = srv.Shutdown(context.Background()); err != nil {
			log.Printf("%+v\n", err)
		}
	}
}
