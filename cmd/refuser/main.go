// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"github.com/go-pogo/errors"
	"github.com/go-pogo/serv"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Refusers refuses any attempt to connect to its server.

func main() {
	var port serv.Port = 80

	cli := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cli.Var(&port, "port", "Port to listen to")
	_ = cli.Parse(os.Args[1:])

	tcp, err := net.Listen("tcp", port.Addr())
	errors.FatalOnErr(err)
	defer tcp.Close()
	log.Println("listening on", port.Addr())

	go func() {
		for {
			conn, err := tcp.Accept()
			if err != nil && !strings.Contains(err.Error(), "closed network connection") {
				log.Println("err:", err)
			}

			log.Println("refuse:", conn.RemoteAddr())
			_ = conn.Close()
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
