// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpheader

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCacheControlResponse_String(t *testing.T) {
	tests := map[string]CacheControlResponse{
		"":                                {},
		"max-age=604800":                  {MaxAge: time.Hour * 24 * 7},
		"s-maxage=604800":                 {SMaxAge: time.Hour * 24 * 7},
		"max-age=3600, must-revalidate":   {MaxAge: time.Hour, MustRevalidate: true},
		"max-age=10, proxy-revalidate":    {MaxAge: time.Second * 10, ProxyRevalidate: true},
		"no-cache":                        {NoCache: true},
		"no-store":                        {NoStore: true},
		"private":                         {Private: true},
		"public":                          {Public: true},
		"public, max-age=7200":            {MaxAge: time.Hour * 2, Public: true},
		"must-understand, no-store":       {MustUnderstand: true, NoStore: true},
		"no-transform":                    {NoTransform: true},
		"public, max-age=3600, immutable": {MaxAge: time.Hour, Public: true, Immutable: true},
		"max-age=604800, stale-while-revalidate=86400": {MaxAge: time.Hour * 24 * 7, StaleWhileRevalidate: time.Hour * 24},
		"max-age=604800, stale-if-error=86400":         {MaxAge: time.Hour * 24 * 7, StaleIfError: time.Hour * 24},
	}

	for want, have := range tests {
		t.Run(want, func(t *testing.T) {
			assert.Equal(t, want, have.String())
		})
	}
}
