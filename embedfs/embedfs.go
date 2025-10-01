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

type FS interface {
	fs.FS
	fs.ReadDirFS
	fs.ReadFileFS
}

type Option func(f *FileServer) error

var _ http.Handler = (*FileServer)(nil)

type FileServer struct {
	FS
	handler http.Handler
	modTime time.Time
}

func New(embed FS, opts ...Option) (*FileServer, error) {
	s := &FileServer{FS: embed}
	if err := s.applyOpts(opts); err != nil {
		return nil, err
	}
	if s.handler == nil {
		s.handler = http.FileServer(http.FS(s.FS))
	}

	return s, nil
}

func (s *FileServer) applyOpts(opts []Option) error {
	var err error
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		err = errors.Append(err, opt(s))
	}
	return err
}

func (s *FileServer) ModTime() time.Time { return s.modTime }

func (s *FileServer) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	s.handler.ServeHTTP(wri, req)
	if !s.modTime.IsZero() {
		wri.Header().Set("Last-Modified", s.modTime.Format(http.TimeFormat))
	}
}
