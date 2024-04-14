// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"github.com/go-pogo/serv"
	"net/http"
)

func IgnoreFaviconHandler() http.HandlerFunc {
	return func(wri http.ResponseWriter, req *http.Request) {
		SetShouldIgnore(req.Context(), true)
		wri.WriteHeader(http.StatusNoContent)
	}
}

func IgnoreFaviconRoute() serv.Route {
	return serv.Route{
		Name:    "favicon",
		Method:  http.MethodGet,
		Pattern: "/favicon.ico",
		Handler: IgnoreFaviconHandler(),
	}
}
