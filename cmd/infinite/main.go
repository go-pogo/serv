// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/serv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Infinite serves an infinite stream of dots.

func main() {
	var port serv.Port = 80

	cli := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cli.Var(&port, "p", "Port of webserver")
	_ = cli.Parse(os.Args[1:])

	mux := serv.NewServeMux()
	mux.HandleRoute(serv.Route{
		Name:    "infinite",
		Method:  http.MethodGet,
		Pattern: "/",
		Handler: &handler{verbose: true},
	})

	ctx, stopFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	srv, err := serv.NewDefault(
		port, mux,
		serv.WithName("infinite"),
		serv.WithBaseContext(ctx),
		serv.WithDefaultLogger(),
	)
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

type handler struct {
	verbose bool
}

func (h *handler) serveInfinite(ctx context.Context) {
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	defer fmt.Print("\n")

	for {
		select {
		case <-timer.C:
			if h.verbose {
				fmt.Print(".")
			}

			timer.Reset(time.Second)
		case <-ctx.Done():
			return
		}
	}
}

func (h *handler) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	if path == "" {
		h.serveInfinite(req.Context())
		return
	}

	result := make(map[string]any, 3)
	result["start"] = time.Now()

	timeout, err := time.ParseDuration(path)

	if err != nil {
		result["err"] = err.Error()
		result["ok"] = false
	} else {
		time.Sleep(timeout)
		result["end"] = time.Now()
		result["ok"] = true
	}

	_ = serv.WriteJSON(wri, result)
}
