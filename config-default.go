// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import "time"

var defaultConfig = Config{
	ReadTimeout:       5 * time.Second,
	ReadHeaderTimeout: 2 * time.Second,
	WriteTimeout:      10 * time.Second,
	IdleTimeout:       120 * time.Second,
	ShutdownTimeout:   60 * time.Second,
	MaxHeaderBytes:    10240, // 10 KiB => 10 * data.Kibibyte
}
