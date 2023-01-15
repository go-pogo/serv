// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpheader

import (
	"strconv"
	"strings"
	"time"
)

const CacheControl = "Cache-Control"

// NoCache sets the CacheControl header to no-cache and the Expires header set
// to a date in the past.
func NoCache() HeaderModifier {
	return HeaderModifierFunc(func(h Header) {
		SetExpires(h, time.Unix(0, 0))
		h.Set(CacheControl, "private, max-age=0, no-cache")
	})
}

var _ HeaderModifier = new(CacheControlResponse)

// CacheControlResponse holds Cache-Control directives for responses that
// control caching in browsers and shared caches (e.g. Proxies, CDNs).
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control#response_directives
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching
type CacheControlResponse struct {
	// Private indicates that the response can be stored only in a private
	// cache (e.g. local caches in browsers).
	Private bool
	// Public indicates that the response can be stored in a shared cache.
	// Responses for requests with Authorization header fields must not be
	// stored in a shared cache; however, the public directive will cause such
	// responses to be stored in a shared cache.
	Public bool
	// MaxAge indicates that the response remains fresh until time.Duration
	// seconds after the response is generated.
	MaxAge time.Duration
	// SMaxAge indicates how long the response is fresh for (similar to MaxAge),
	// but it is specific to shared caches which will ignore max-age when it is
	// present.
	SMaxAge time.Duration
	// StaleWhileRevalidate indicates that the cache could reuse a stale
	// response while it revalidates it to a cache.
	StaleWhileRevalidate time.Duration
	// StaleIfError indicates that the cache can reuse a stale response when an
	// origin server responds with an error (500, 502, 503, or 504).
	StaleIfError time.Duration
	// MustRevalidate indicates that the response can be stored in caches and
	// can be reused while fresh. If the response becomes stale, it must be
	// validated with the origin server before reuse.
	// MustRevalidate is typically used with MaxAge.
	MustRevalidate bool
	// ProxyRevalidate is the equivalent of MustRevalidate, but specifically
	//for shared caches only.
	ProxyRevalidate bool
	// MustUnderstand indicates that a cache should store the response only if
	// it understands the requirements for caching based on status code.
	// It should be coupled with no-store for fallback behavior.
	MustUnderstand bool
	// NoCache indicates that the response can be stored in caches, but the
	// response must be validated with the origin server before each reuse,
	// even when the cache is disconnected from the origin server.
	NoCache bool
	// NoStore indicates that any caches of any kind (private or shared) should
	// not store this response.
	NoStore bool
	// NoTransform indicates that any intermediary (regardless of whether it
	// implements a cache) shouldn't transform the response contents.
	NoTransform bool
	// Immutable indicates that the response will not be updated while it's
	//fresh.
	Immutable bool
}

// SetCacheControl sets the Cache-Control directives from CacheControlResponse
// to h.
func SetCacheControl(h Header, c CacheControlResponse) { c.ModifyHeader(h) }

func (c CacheControlResponse) ModifyHeader(header Header) {
	cc := c.String()
	if cc == "" {
		header.Del(CacheControl)
	} else {
		header.Set(CacheControl, cc)
		if c.NoCache {
			// ensure backwards compatibility with older browsers
			// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Pragma
			header.Set("Pragma", "no-cache")
		}
	}
}

func (c CacheControlResponse) String() string {
	var sb strings.Builder
	if c.Private {
		write(&sb, "private")
	}
	if c.Public {
		write(&sb, "public")
	}

	if c.MaxAge > time.Second {
		writeDuration(&sb, "max-age=", c.MaxAge)
	} else if c.Private && c.NoCache {
		write(&sb, "max-age=0")
	}

	if c.SMaxAge > time.Second {
		writeDuration(&sb, "s-maxage=", c.SMaxAge)
	}
	if c.StaleWhileRevalidate > time.Second {
		writeDuration(&sb, "stale-while-revalidate=", c.StaleWhileRevalidate)
	}
	if c.StaleIfError > time.Second {
		writeDuration(&sb, "stale-if-error=", c.StaleIfError)
	}

	if c.MustRevalidate {
		write(&sb, "must-revalidate")
	}
	if c.ProxyRevalidate {
		write(&sb, "proxy-revalidate")
	}
	if c.MustUnderstand {
		write(&sb, "must-understand")
	}

	if c.NoCache {
		write(&sb, "no-cache")
	}
	if c.NoStore {
		write(&sb, "no-store")
	}
	if c.NoTransform {
		write(&sb, "no-transform")
	}
	if c.Immutable {
		write(&sb, "immutable")
	}

	return sb.String()
}

func write(sb *strings.Builder, str ...string) {
	if sb.Len() != 0 {
		sb.WriteString(", ")
	}
	for _, s := range str {
		sb.WriteString(s)
	}
}

func writeDuration(sb *strings.Builder, str string, d time.Duration) {
	write(sb, str, seconds(d))
}

func seconds(d time.Duration) string {
	return strconv.FormatInt(int64(d/time.Second), 10)
}
