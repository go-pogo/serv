// Copyright (c) 2026, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !go1.26

package serv

import "net/http"

func (srv *Server) resetServer() {
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
