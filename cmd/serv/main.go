// Copyright (c) 2021, Roel Schut. All rights reserved.
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
	"syscall"
	"time"
)

func main() {
	cli := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	port := serv.Port(*cli.Uint("p", 8080, "server port"))

	err := cli.Parse(os.Args[1:])
	errors.FatalOnErr(err)

	dir := cli.Arg(0)
	if dir == "" {
		dir = "./"
	}

	handler := http.FileServer(http.Dir(dir))
	handler = http.TimeoutHandler(handler, time.Second, "")

	mux := http.NewServeMux()
	mux.Handle("/", accesslog.WithHandlerName("files", handler))

	srv, err := serv.NewDefault(
		port,
		serv.WithHandler(mux),
		serv.WithName("serv"),
		serv.WithDefaultLogger(),
		serv.WithMiddleware(
			accesslog.Middleware(accesslog.DefaultLogger(nil)),
		),
	)
	errors.FatalOnErr(err)

	ctx, quitFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer quitFn()

	errCh := make(chan error)
	defer close(errCh)

	go func() {
		if err := srv.Run(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-errCh:
		log.Println("Server error:", err.Error())
	}

	if err = srv.Shutdown(context.Background()); err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Shutdown error: %+v\n", err)
		} else if err = srv.Close(); err != nil {
			log.Printf("Close error: %+v\n", err)
		}
	}
}
