// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"log"
	"net/http"

	"github.com/go-pogo/serv/metrics"
)

type AccessLogger interface {
	LogAccess(m metrics.Metrics, r *http.Request)
}

func NewRecorder(log AccessLogger) metrics.Recorder {
	if log == nil {
		return nil
	}

	return &recorder{
		log: log,
	}
}

type recorder struct {
	log AccessLogger
}

func (c *recorder) Record(met metrics.Metrics, req *http.Request) {
	c.log.LogAccess(met, req)
}

type Logger struct{}

func (l Logger) LogAccess(m metrics.Metrics, r *http.Request) {
	log.Printf("%s \"%s %s %s\" %d %db\n",
		RemoteAddr(r),
		r.Method,
		RequestURI(r),
		r.Proto,
		m.Code,
		m.Written,
	)
}
