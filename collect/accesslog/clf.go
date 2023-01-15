// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package accesslog

import (
	"io"
	"strconv"
	"strings"

	"github.com/go-pogo/errors"
	"github.com/go-pogo/writing"
)

// https://httpd.apache.org/docs/current/logs.html#common
type ClfFormatter struct{ Entry }

func (clf *ClfFormatter) Username() string {
	if u := clf.Entry.Username(); u != "" {
		return u
	}
	return "-"
}

func (clf *ClfFormatter) String() string {
	var b strings.Builder
	clf.writeTo(&b)
	return b.String()
}

func (clf *ClfFormatter) WriteTo(w io.Writer) (n int64, err error) {
	sw := writing.ToCountingStringWriter(w)
	clf.writeTo(sw)
	_, _ = sw.WriteString("\n")
	return int64(sw.Count()), errors.Join(sw.Errors()...)
}

func (clf *ClfFormatter) writeTo(sw writing.StringWriter) {
	_, _ = sw.WriteString(clf.RemoteAddr())
	_, _ = sw.WriteString(" - ")
	_, _ = sw.WriteString(clf.Username())
	_, _ = sw.WriteString(" [")
	_, _ = sw.WriteString(clf.Timestamp())
	_, _ = sw.WriteString("] \"")
	_, _ = sw.WriteString(clf.Request.Method)
	_, _ = sw.WriteString(" ")
	_, _ = sw.WriteString(clf.RequestURI())
	_, _ = sw.WriteString(" ")
	_, _ = sw.WriteString(clf.Request.Proto)
	_, _ = sw.WriteString("\" ")
	_, _ = sw.WriteString(strconv.Itoa(clf.Metrics.Code))
	_, _ = sw.WriteString(" ")
	_, _ = sw.WriteString(strconv.FormatInt(clf.Metrics.Written, 10))
}
