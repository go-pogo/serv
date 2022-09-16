// Copyright (c) 2022, Roel Schut. All rights reserved.
// applyOptions of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cors

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-pogo/serv"
	"github.com/go-pogo/serv/httpheader"
)

type AccessControl struct {
	// AllowCredentials sets the Access-Control-Allow-Credentials response
	// header which tells browsers whether to expose the response to the
	// frontend JavaScript code.
	// When a request's credentials mode is `include`, browsers will only
	// expose the response if the Access-Control-Allow-Credentials header value
	// is true.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials
	AllowCredentials bool

	// AllowHeaders sets the Access-Control-Allow-Headers response header which
	// is used in response to a preflight request which includes the
	// Access-Control-Request-Headers to indicate which HTTP headers can be
	// used during the actual request.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers
	AllowHeaders []string

	// AllowMethods sets the Access-Control-Allow-Methods response header which
	// specifies one or more methods allowed when accessing a resource in
	// response to a preflight request.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods
	AllowMethods []string

	// AllowOrigin sets the Access-Control-Allow-Origin response header which
	// indicates whether the response can be shared with requesting code from
	// the given origin.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin
	AllowOrigin []string

	// ExposeHeaders sets the Access-Control-Expose-Headers response header
	// which allows a server to indicate which response headers should be made
	// available to scripts running in the browser, in response to a
	// cross-origin request.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers
	ExposeHeaders []string

	// MaxAge sets the Access-Control-Max-Age response header which indicates
	// how long the results of a preflight request (that is the information
	// contained in the Access-Control-Allow-Methods and
	// Access-Control-Allow-Headers headers) can be cached.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
	MaxAge time.Duration
}

const (
	wildcard = "*"
	null     = "null"
)

func (ac AccessControl) ModifyHeader(h serv.Header) {
	if ac.AllowCredentials {
		h.Set(httpheader.AccessControlAllowCredential, "true")
	}

	if len(ac.AllowHeaders) != 0 {
		h.Set(httpheader.AccessControlAllowHeaders, joinOrWildcard(ac.AllowHeaders))
	}

	if len(ac.AllowMethods) != 0 {
		h.Set(httpheader.AccessControlAllowMethods, strings.Join(ac.AllowMethods, ", "))
	}

	if len(ac.AllowOrigin) != 0 {
		if contains(ac.AllowOrigin, wildcard) {
			h.Set(httpheader.AccessControlAllowOrigin, wildcard)
		} else if contains(ac.AllowOrigin, null) {
			h.Set(httpheader.AccessControlAllowOrigin, null)
		} else {
			h.Set(httpheader.AccessControlAllowOrigin, strings.Join(ac.AllowOrigin, ", "))
			h.Set("Vary", "Origin")
		}
	}

	if len(ac.ExposeHeaders) != 0 {
		h.Set(httpheader.AccessControlExposeHeaders, joinOrWildcard(ac.ExposeHeaders))
	}

	if ac.MaxAge > 0 {
		if secs := math.Round(ac.MaxAge.Seconds()); secs >= 1.0 {
			h.Set(httpheader.AccessControlMaxAge, strconv.FormatFloat(secs, 'f', 0, 64))
		}
	}
}

func joinOrWildcard(elem []string) string {
	if contains(elem, wildcard) {
		return wildcard
	}

	return strings.Join(elem, ", ")
}

func contains(elems []string, val string) bool {
	for _, elem := range elems {
		if elem == val {
			return true
		}
	}
	return false
}
