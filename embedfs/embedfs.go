// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package embedfs

import (
	"github.com/go-pogo/errors"
	"io/fs"
	"net/http"
	"time"
)

type FS interface {
	fs.FS
	fs.ReadDirFS
	fs.ReadFileFS
}

type Option func(f *FileServer) error

var _ http.Handler = new(FileServer)

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

func (s *FileServer) ModTime() time.Time { return s.modTime }

func (s *FileServer) applyOpts(opts []Option) error {
	var err error
	for _, opt := range opts {
		errors.AppendInto(&err, opt(s))
	}
	return err
}

func (s *FileServer) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	s.handler.ServeHTTP(wri, req)
	if !s.modTime.IsZero() {
		wri.Header().Set("Last-Modified", s.modTime.Format(http.TimeFormat))
	}
}

func WithSubDir(dir string) Option {
	return func(s *FileServer) error {
		if dir == "" {
			return nil
		}

		sub, err := fs.Sub(s.FS, dir)
		if err != nil {
			return errors.WithStack(err)
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
