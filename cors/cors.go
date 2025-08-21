// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cors

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	AccessControlAllowCredentialsKey = "Access-Control-Allow-Credentials"
	AccessControlAllowHeadersKey     = "Access-Control-Allow-Headers"
	AccessControlAllowMethodsKey     = "Access-Control-Allow-Methods"
	AccessControlAllowOriginKey      = "Access-Control-Allow-Origin"
	AccessControlExposeHeadersKey    = "Access-Control-Expose-Headers"
	AccessControlMaxAgeKey           = "Access-Control-Max-Age"
	AccessControlRequestHeadersKey   = "Access-Control-Request-Headers"
	AccessControlRequestMethodKey    = "Access-Control-Request-Method"

	originKey  = "Origin"
	refererKey = "Referer"
	varyKey    = "Vary"
)

func Middleware(ac AccessControl) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return NewHandler(next, ac)
	}
}

// AccessControl
// https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
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

var _ http.Handler = (*handler)(nil)

type handler struct {
	AccessControl
	next http.Handler
}

func NewHandler(next http.Handler, ac AccessControl) http.Handler {
	return &handler{
		AccessControl: ac,
		next:          next,
	}
}

func (ac *handler) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	var origin string
	if corsUnsafe {
		origin = req.Header.Get(originKey)
		if origin == "" {
			origin = req.Header.Get(refererKey)
		}
	}
	ac.modifyHeader(wri.Header(), origin)
	ac.next.ServeHTTP(wri, req)
}

func (ac *handler) modifyHeader(h http.Header, origin string) {
	if ac.AllowCredentials {
		h.Set(AccessControlAllowCredentialsKey, "true")
	}

	if len(ac.AllowHeaders) != 0 {
		h.Set(AccessControlAllowHeadersKey, joinOrWildcard(ac.AllowHeaders))
	}

	if len(ac.AllowMethods) != 0 {
		h.Set(AccessControlAllowMethodsKey, strings.Join(ac.AllowMethods, ", "))
	}

	if allow := ac.resolveAllowOrigin(origin); allow != "" {
		h.Set(AccessControlAllowOriginKey, allow)
		h.Add(varyKey, "Origin")
	}

	if len(ac.ExposeHeaders) != 0 {
		h.Set(AccessControlExposeHeadersKey, joinOrWildcard(ac.ExposeHeaders))
	}

	if ac.MaxAge > 0 {
		if secs := math.Round(ac.MaxAge.Seconds()); secs >= 1.0 {
			h.Set(AccessControlMaxAgeKey, strconv.FormatFloat(secs, 'f', 0, 64))
		}
	}
}

func (ac *handler) resolveAllowOrigin(origin string) string {
	if n := len(ac.AllowOrigin); n == 0 {
		return ""
	} else if n == 1 {
		res := ac.AllowOrigin[0]
		if res == "*" && ac.AllowCredentials && origin != "" {
			res = origin
		}
		return res
	}

	wildcard := contains(ac.AllowOrigin, "*")
	if wildcard && ac.AllowCredentials && origin != "" {
		return origin
	}

	var res string
	if wildcard {
		res = "*"
	}
	if origin != "" {
		for _, ao := range ac.AllowOrigin {
			if ao == origin {
				res = origin
				break
			}
		}
	}
	return res
}

func joinOrWildcard(elem []string) string {
	if contains(elem, "*") {
		return "*"
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
