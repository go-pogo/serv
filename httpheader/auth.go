// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpheader

import "encoding/base64"

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Authorization
const Authorization = "Authorization"

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#authentication_schemes
type AuthScheme interface {
	Scheme() string
	String() string
}

func SetAuthorization(h Header, auth AuthScheme) {
	h.Set(Authorization, auth.String())
}

var _ AuthScheme = new(BasicAuth)

type BasicAuth struct {
	Username string
	Password string
}

func (ba BasicAuth) Scheme() string { return "Basic" }

func (ba BasicAuth) String() string {
	auth := ba.Username + ":" + ba.Password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

var _ AuthScheme = new(BearerAuth)

type BearerAuth string

func (ba BearerAuth) Token() string  { return string(ba) }
func (ba BearerAuth) Scheme() string { return "Bearer" }
func (ba BearerAuth) String() string { return "Bearer " + string(ba) }
