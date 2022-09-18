// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-pogo/errors"
	"github.com/go-pogo/serv"
	"github.com/go-pogo/serv/accesslog"
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

	router := http.NewServeMux()
	router.Handle("/", handler)

	server, err := serv.NewDefault(
		accesslog.Collect(
			new(accesslog.Logger),
			accesslog.ResponseStatusAll,
			router,
		),
		port,
	)
	errors.FatalOnErr(err)

	ctx, quitFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer quitFn()

	go func() {
		server.BaseContext = serv.BaseContext(ctx)
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			_, _ = fmt.Fprintf(os.Stderr, "\nServer error: %+v\n", err)
		}
	}()

	<-ctx.Done()

	ctx, closeFn := context.WithTimeout(context.Background(), time.Second*3)
	defer closeFn()

	done := make(chan error)
	go func() {
		done <- server.Shutdown(ctx)
	}()

	select {
	case err := <-done:
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "\nShutdown error: %+v\n", err)
		}
	case <-ctx.Done():
		if err := server.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "\nClose error: %+v\n", err)
		}
	}
}
