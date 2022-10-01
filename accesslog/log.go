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

type Level int

const (
	ResponseStatusSuccess     Level = 1 << iota // 200-299
	ResponseStatusRedirect                      // 300-399
	ResponseStatusClientError                   // 400-499
	ResponseStatusServerError                   // 500-599

	ResponseStatusNone Level = 0
	ResponseStatusAll        = ResponseStatusSuccess | ResponseStatusRedirect | ResponseStatusClientError | ResponseStatusServerError
)

func (l Level) InRange(statusCode int) bool {
	return (l&ResponseStatusSuccess != 0 && statusCode >= 200 && statusCode <= 299) ||
		(l&ResponseStatusRedirect != 0 && statusCode >= 300 && statusCode <= 399) ||
		(l&ResponseStatusClientError != 0 && statusCode >= 400 && statusCode <= 499) ||
		(l&ResponseStatusServerError != 0 && statusCode >= 500 && statusCode <= 599)
}

func NewRecorder(log AccessLogger, level Level) metrics.Recorder {
	if log == nil || level == ResponseStatusNone {
		return nil
	}

	return &recorder{
		level: level,
		log:   log,
	}
}

type recorder struct {
	level Level
	log   AccessLogger
}

func (c *recorder) Record(met metrics.Metrics, req *http.Request) {
	if !c.level.InRange(met.Code) {
		return
	}

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
