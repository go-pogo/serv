// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !cors_unsafe
// +build !cors_unsafe

package httpheader

import (
	"net/http"
)

func extractOrigin(http.Header) string { return "" }
