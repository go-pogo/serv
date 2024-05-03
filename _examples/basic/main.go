// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/serv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// This program starts and runs a basic server which will shutdown after
// receiving a SIGINT or SIGTERM signal.

func main() {
	var port serv.Port = 8080

	cli := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cli.Var(&port, "port", "Server port")
	_ = cli.Parse(os.Args[1:])

	mux := serv.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, world!"))
	})

	ctx, stopFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	srv, err := serv.New(port, mux, serv.WithBaseContext(ctx), serv.WithDefaultLogger())
	errors.FatalOnErr(err)

	go func() {
		defer stopFn()
		if err := srv.Run(); err != nil {
			log.Println("Server error:", err.Error())
		}
	}()
	<-ctx.Done()

	if err = srv.Shutdown(context.Background()); err != nil {
		if !errors.Is(err, context.DeadlineExceeded) {
			log.Printf("Shutdown error: %+v\n", err)
		} else if err = srv.Close(); err != nil {
			log.Printf("Close error: %+v\n", err)
		}
	}
}
