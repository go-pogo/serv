// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package metrics

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/felixge/httpsnoop"
)

type snoop = httpsnoop.Metrics

// Metrics holds metrics collected with Collect.
type Metrics struct {
	snoop
	// The Time the request was received.
	Time    time.Time
	Traffic int64
}

// Recorder receives the Metrics data decides what to do with it. This can range
// from logging the data and/or making it available for scraping (prometheus).
type Recorder interface {
	Record(met Metrics, req *http.Request)
}

// Collect metrics from http.Handler h and pass it to a Recorder.
func Collect(h http.Handler, rec ...Recorder) *Collector {
	col := &Collector{
		next: h,
		rec:  make([]Recorder, 0, len(rec)),
	}
	for _, r := range rec {
		col.WithRecorder(r)
	}
	return col
}

var _ http.Handler = &Collector{}

type Collector struct {
	next    http.Handler
	rec     []Recorder
	traffic atomic.Int64
}

func (col *Collector) WithRecorder(rec Recorder) *Collector {
	if rec != nil {
		col.rec = append(col.rec, rec)
	}
	return col
}

func (col *Collector) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	var m Metrics
	m.Time = time.Now()
	m.Traffic = col.traffic.Add(1)
	defer col.traffic.Add(-1)

	m.snoop = httpsnoop.CaptureMetrics(col.next, wri, req)

	go func() {
		r := req.Clone(context.Background())
		for _, rec := range col.rec {
			rec.Record(m, r)
		}
	}()
}
