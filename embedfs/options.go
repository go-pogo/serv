// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package embedfs

import (
	"io/fs"
	"net/http"
	"time"

	"github.com/go-pogo/errors"
)

const ErrInvalidSubDir errors.Msg = "invalid sub directory"

func WithSubDir(dir string) Option {
	return func(s *FileServer) error {
		if dir == "" {
			return nil
		}

		sub, err := fs.Sub(s.FS, dir)
		if err != nil {
			return errors.Wrap(err, ErrInvalidSubDir)
		}

		s.handler = http.FileServer(http.FS(sub))
		return nil
	}
}

func WithModTime(t time.Time) Option {
	return func(s *FileServer) error {
		s.modTime = t.UTC()
		return nil
	}
}
