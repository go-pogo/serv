// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prometheus

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-pogo/serv/collect"
	"github.com/prometheus/client_golang/prometheus"
)

type Config struct {
	Namespace       string
	Subsystem       string
	HandlerLabel    string
	PathLabel       string
	MethodLabel     string
	CodeLabel       string
	DurationBuckets []float64
	SizeBuckets     []float64
}

// LatencyBuckets is like prometheus.DefBuckets but uses time as milliseconds.
func LatencyBuckets() []float64 {
	return []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}
}

// SizeBuckets returns exponential buckets from 10B to 1GB.
func SizeBuckets() []float64 {
	return prometheus.ExponentialBuckets(10, 10, 9)
}

func (c *Config) defaults() {
	if c.Subsystem == "" {
		c.Subsystem = "http_server"
	}
	if c.HandlerLabel == "" {
		c.HandlerLabel = "handler"
	}
	if c.PathLabel == "" {
		c.PathLabel = "path"
	}
	if c.MethodLabel == "" {
		c.MethodLabel = "method"
	}
	if c.CodeLabel == "" {
		c.CodeLabel = "code"
	}
	if len(c.DurationBuckets) == 0 {
		c.DurationBuckets = LatencyBuckets()
	}
	if len(c.SizeBuckets) == 0 {
		c.SizeBuckets = SizeBuckets()
	}
}

var _ collect.Collector = new(recorder)

type recorder struct {
	requests *prometheus.GaugeVec
	latency  *prometheus.HistogramVec
	size     *prometheus.HistogramVec
}

func NewRecorder(prom prometheus.Registerer, conf *Config) collect.Collector {
	if conf == nil {
		conf = new(Config)
	}

	conf.defaults()
	if prom == nil {
		prom = prometheus.DefaultRegisterer
	}

	r := &recorder{
		requests: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: conf.Namespace,
			Subsystem: conf.Subsystem,
			Name:      "active_requests",
			Help:      "Amount of simultaneous active requests.",
		}, []string{conf.HandlerLabel, conf.PathLabel}),

		latency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: conf.Namespace,
			Subsystem: conf.Subsystem,
			Name:      "duration_ms",
			Help:      "Duration (latency) of the requests.",
			Buckets:   conf.DurationBuckets,
		}, []string{conf.HandlerLabel, conf.PathLabel, conf.MethodLabel, conf.CodeLabel}),

		size: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: conf.Namespace,
			Subsystem: conf.Subsystem,
			Name:      "response_size_bytes",
			Help:      "Size of the response bodies.",
			Buckets:   conf.SizeBuckets,
		}, []string{conf.HandlerLabel, conf.PathLabel, conf.MethodLabel, conf.CodeLabel}),
	}
	prom.MustRegister(r.latency, r.size, r.requests)
	return r
}

func (r *recorder) Collect(ctx context.Context, met collect.Metrics, req *http.Request) {
	name := collect.HandlerName(ctx)
	code := strconv.Itoa(met.Code)

	r.requests.WithLabelValues(name, req.URL.Path).Add(float64(met.Traffic))

	r.latency.WithLabelValues(name, req.URL.Path, req.Method, code).
		Observe(float64(met.Duration.Milliseconds()))

	r.size.WithLabelValues(name, req.URL.Path, req.Method, code).
		Observe(float64(met.Written))
}
