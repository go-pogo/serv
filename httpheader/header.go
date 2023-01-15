// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpheader

import (
	"github.com/go-pogo/serv/middleware"
	"net/http"
	"strings"
	"time"
)

var _ Header = new(http.Header)

// Header is a http.Header compatible header implementation.
type Header interface {
	Add(key, value string)
	Set(key, value string)
	Get(key string) string
	Values(key string) []string
	Del(key string)
}

// HeaderModifier modifies header.
type HeaderModifier interface {
	ModifyHeader(header Header)
}

type HeaderModifierFunc func(h Header)

func (fn HeaderModifierFunc) ModifyHeader(h Header) { fn(h) }

// Modify header with the supplied HeaderModifier(s).
func Modify(header Header, modify ...HeaderModifier) {
	for _, m := range modify {
		m.ModifyHeader(header)
	}
}

var nopHttpHandler = func(_ http.ResponseWriter, _ *http.Request) {}

// Middleware wraps the Modify function with middleware.HandlerFunc.
func Middleware(modify ...HeaderModifier) middleware.Middleware {
	return middleware.HandlerFunc(func(wri http.ResponseWriter, req *http.Request) {
		header := wri.Header()
		for _, m := range modify {
			if mw, ok := m.(middleware.Middleware); ok {
				mw.Wrap(nopHttpHandler).ServeHTTP(wri, req)
			} else {
				m.ModifyHeader(header)
			}
		}
	})
}

const (
	ContentType = "Content-Type"

	Origin  = "Origin"
	Referer = "Referer"
	Cookie  = "Cookie"

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Vary
	Vary = "Vary"
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent
	UserAgent = "User-Agent"

	// cache
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Expires
	Expires = "Expires"
	//https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Last-Modified
	LastModified = "Last-Modified"

	// auth

	WWWAuthenticate = "WWW-Authenticate"
)

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept
const Accept = "Accept"

func SetAccept(h Header, accept ...string) { setMultiple(h, Accept, accept) }

func SetAcceptAny(h Header) { h.Set(Accept, "*/*") }

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Encoding
const AcceptEncoding = "Accept-Encoding"

func SetAcceptEncoding(h Header, accept ...string) { setMultiple(h, AcceptEncoding, accept) }

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Age
const Age = "Age"

func SetAge(h Header, age time.Duration) { h.Set(Age, seconds(age)) }

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Allow
const Allow = "Allow"

func SetAllow(h Header, allow ...string) { setMultiple(h, Allow, allow) }

func SetExpires(h Header, t time.Time) { h.Set(Expires, t.UTC().Format(http.TimeFormat)) }

func SetLastModified(h Header, t time.Time) { h.Set(LastModified, t.UTC().Format(http.TimeFormat)) }

func SetUserAgent(h Header, name, version string, comments ...string) {
	if version != "" {
		name += "/" + version
	}
	if len(comments) == 0 {
		h.Set(UserAgent, name)
	} else {
		comments = append(comments, "")
		copy(comments[1:], comments)
		comments[0] = name
		h.Set(UserAgent, strings.Join(comments, " "))
	}
}

func AllowDefaultMethods() []string {
	return []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPatch,
		http.MethodPut,
		http.MethodDelete,
		http.MethodOptions,
	}
}

// Copy header value associated with key from src to dst.
func Copy(dst, src Header, key ...string) {
	for _, k := range key {
		dst.Set(k, src.Get(k))
	}
}

func Delete(h Header, key ...string) {
	for _, k := range key {
		h.Del(k)
	}
}

func setMultiple(h Header, key string, value []string) {
	if n := len(value); n == 1 {
		h.Set(key, value[0])
	} else if n > 1 {
		h.Set(key, value[0])
		for _, v := range value[1:] {
			h.Add(key, v)
		}
	}
}
