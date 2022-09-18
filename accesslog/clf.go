// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-pogo/errors"
	"github.com/go-pogo/writing"
)

func RemoteAddr(r *http.Request) string {
	addr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		addr = r.RemoteAddr
	}
	return addr
}

// https://www.rfc-editor.org/rfc/rfc7540#section-8.3
func RequestURI(r *http.Request) string {
	var uri string
	if r.ProtoMajor == 2 && r.Method == "CONNECT" {
		uri = r.Host
	} else {
		uri = r.RequestURI
	}
	if uri == "" {
		uri = r.URL.RequestURI()
	}
	return uri
}

func Username(r *http.Request) string {
	if r.URL != nil && r.URL.User != nil {
		if user := r.URL.User.Username(); user != "" {
			return user
		}
	}
	return ""
}

// https://httpd.apache.org/docs/current/logs.html#common
type CLF Access

func (clf *CLF) Username() string {
	if u := Username(clf.Request); u != "" {
		return u
	}
	return "-"
}

func (clf *CLF) Timestamp() string {
	return clf.Time.Format("02/Jan/2006:15:04:05 -0700")
}

func (clf *CLF) String() string {
	var b strings.Builder
	clf.writeTo(&b)
	return b.String()
}

func (clf *CLF) WriteTo(w io.Writer) (n int64, err error) {
	sw := writing.ToCountingStringWriter(w)
	clf.writeTo(sw)
	_, _ = sw.WriteString("\n")
	return int64(sw.Count()), errors.Combine(sw.Errors()...)
}

func (clf *CLF) writeTo(sw writing.StringWriter) {
	_, _ = sw.WriteString(RemoteAddr(clf.Request))
	_, _ = sw.WriteString(" - ")
	_, _ = sw.WriteString(clf.Username())
	_, _ = sw.WriteString(" [")
	_, _ = sw.WriteString(clf.Timestamp())
	_, _ = sw.WriteString("] \"")
	_, _ = sw.WriteString(clf.Request.Method)
	_, _ = sw.WriteString(" ")
	_, _ = sw.WriteString(RequestURI(clf.Request))
	_, _ = sw.WriteString(" ")
	_, _ = sw.WriteString(clf.Request.Proto)
	_, _ = sw.WriteString("\" ")
	_, _ = sw.WriteString(strconv.Itoa(clf.StatusCode))
	_, _ = sw.WriteString(" ")
	_, _ = sw.WriteString(strconv.FormatInt(clf.Size, 10))
}
