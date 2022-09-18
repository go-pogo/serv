// Copyright (c) 2022, Roel Schut. All rights reserved.
// applyOptions of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"log"
)

type AccessLogger interface {
	LogAccess(a *Access)
}

type Logger struct{}

func (l Logger) LogAccess(a *Access) {
	log.Printf("%s \"%s %s %s\" %d %d\n",
		RemoteAddr(a.Request),
		a.Request.Method,
		RequestURI(a.Request),
		a.Request.Proto,
		a.StatusCode,
		a.Size,
	)
}
