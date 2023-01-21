// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/serv"
	"github.com/go-pogo/serv/accesslog"
	"github.com/go-pogo/serv/collect"
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

	router := http.NewServeMux()
	router.Handle("/", handler)

	server, err := serv.NewDefault(
		collect.Wrap(router,
			collect.LimitCodes(
				collect.ResponseStatusErrors,
				accesslog.Collector(),
			),
		),
		port,
		serv.WithLogger(new(serv.DefaultLogger)),
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

	if err = server.Shutdown(nil); err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			_, _ = fmt.Fprintf(os.Stderr, "\nShutdown error: %+v\n", err)
		} else if err = server.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "\nClose error: %+v\n", err)
		}
	}
}
