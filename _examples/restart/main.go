// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/serv"
	"github.com/go-pogo/serv/accesslog"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
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
			_, _ = w.Write([]byte(strconv.FormatUint(restartCount.Load(), 10)))
		}),
	})
	mux.HandleRoute(accesslog.IgnoreFaviconRoute())
	mux.HandleRoute(serv.Route{
		Name:    "restart",
		Pattern: "/restart",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("Restarting server..."))

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
		serv.WithHandler(accesslog.Middleware(accesslog.DefaultLogger(nil), mux)),
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
