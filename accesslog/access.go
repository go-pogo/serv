// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
)

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

type Access struct {
	// Request received from the client.
	Request *http.Request
	// The Time the request was received.
	Time time.Time
	// The Duration of handling the request.
	Duration time.Duration
	// StatusCode which was sent back to the client.
	StatusCode int
	// Size of the response body returned to the client.
	Size int64
}

const panicNilHandler = "accesslog.Collect: http.Handler must not be nil"

func Collect(log AccessLogger, level Level, next http.Handler) http.Handler {
	if next == nil {
		panic(panicNilHandler)
	}
	if log == nil || level == ResponseStatusNone {
		return next
	}

	return &collector{
		next:  next,
		level: level,
		log:   log,
	}
}

type collector struct {
	next  http.Handler
	level Level
	log   AccessLogger
}

func (c *collector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	metrics := httpsnoop.CaptureMetrics(c.next, w, r)
	if !c.level.InRange(metrics.Code) {
		return
	}

	c.log.LogAccess(&Access{
		Request:    r,
		Time:       start,
		Duration:   metrics.Duration,
		StatusCode: metrics.Code,
		Size:       metrics.Written,
	})
}
