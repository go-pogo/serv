// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package otelaccesslog

import (
	"context"
	"github.com/go-pogo/serv/accesslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"net/url"
)

type exporter struct {
	log accesslog.Logger
}

// NewExporter creates a new otel SpanExporter which, when added to an
// otel provider, sends accesslog.Details derived from
// trace.SpanKindServer spans, to the provided accesslog.Logger.
//
//	var logger accesslog.DefaultLogger
//	tracer := tracesdk.NewTraceProvider(
//		tracesdk.WithBatcher(otelaccesslog.NewExporter(&logger)),
//	)
func NewExporter(log accesslog.Logger) tracesdk.SpanExporter {
	return &exporter{
		log: log,
	}
}

// ExportSpans exports a batch of trace.SpanKindServer spans. All other kinds
// are ignored.
func (exp *exporter) ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error {
	for _, span := range spans {
		if span.SpanKind() != trace.SpanKindServer {
			continue
		}

		var det accesslog.Details
		req := &http.Request{URL: &url.URL{}}

		det.ServerName = span.Name()
		det.StartTime = span.StartTime()
		det.Duration = span.EndTime().Sub(span.StartTime())

		for _, attr := range span.Attributes() {
			switch attr.Key {
			case semconv.CodeFunctionKey:
				det.HandlerName = attr.Value.AsString()
			case semconv.HTTPMethodKey:
				req.Method = attr.Value.AsString()
			case semconv.HTTPTargetKey:
				req.URL.Path = attr.Value.AsString()
			case semconv.HTTPSchemeKey:
				req.URL.Scheme = attr.Value.AsString()
			case semconv.HTTPFlavorKey:
				req.Proto = "HTTP/" + attr.Value.AsString()
			case semconv.NetPeerIPKey, "net.sock.peer.addr":
				if req.RemoteAddr == "" {
					req.RemoteAddr = attr.Value.AsString()
				} else {
					req.RemoteAddr = attr.Value.AsString() + ":" + req.RemoteAddr
				}
			case semconv.NetPeerPortKey, "net.sock.peer.port":
				if req.RemoteAddr == "" {
					req.RemoteAddr = attr.Value.AsString()
				} else {
					req.RemoteAddr += ":" + attr.Value.AsString()
				}
			case semconv.HTTPStatusCodeKey:
				det.StatusCode = int(attr.Value.AsInt64())
			case otelhttp.WroteBytesKey:
				det.BytesWritten = attr.Value.AsInt64()
			case otelhttp.RequestCount:
				det.RequestCount = attr.Value.AsInt64()
			}
		}

		exp.log.Log(ctx, det, req)
	}
	return nil
}

// Shutdown does nothing.
func (*exporter) Shutdown(context.Context) error { return nil }
