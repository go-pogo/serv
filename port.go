// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"encoding"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/go-pogo/errors"
)

const (
	ErrMissingPort   errors.Msg = "missing port"
	ErrInvalidFormat errors.Msg = "invalid format"

	colon = ':'
)

type ParseError struct {
	Err  error
	Text string
}

func (p *ParseError) Unwrap() error { return p.Err }

func (p *ParseError) Error() string {
	if e, ok := p.Err.(errors.Msg); ok {
		return fmt.Sprintf("`%s`: %s", p.Text, e.Error())
	}
	return p.Err.Error()
}

var (
	_ encoding.TextMarshaler   = new(Port)
	_ encoding.TextUnmarshaler = new(Port)
)

type Port uint16

func ParsePort(s string) (Port, error) {
	if s == "" {
		return 0, errors.WithStack(&ParseError{
			Err:  ErrMissingPort,
			Text: s,
		})
	}

	if i := strings.IndexRune(s, colon); i == 0 {
		s = s[1:]
	} else if i > 0 {
		return 0, errors.WithStack(&ParseError{
			Err:  ErrInvalidFormat,
			Text: s,
		})
	}

	x, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0, errors.WithStack(&ParseError{
			Err:  ErrMissingPort,
			Text: s,
		})
	}
	return Port(x), nil
}

func SplitHostPort(hostport string) (string, Port, error) {
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		if addrErr, ok := err.(*net.AddrError); ok &&
			addrErr.Err == "missing port in address" {
			err = ErrMissingPort
		}
		return "", 0, errors.WithStack(&ParseError{
			Err:  err,
			Text: hostport,
		})
	}

	p, err := ParsePort(port)
	return host, p, err
}

func (p *Port) UnmarshalText(text []byte) (err error) {
	*p, err = ParsePort(string(text))
	return err
}

func (p Port) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p Port) String() string {
	if p == 0 {
		return ""
	}
	return strconv.FormatUint(uint64(p), 10)
}

func (p Port) Addr() string {
	if p == 0 {
		return ""
	}
	return string(colon) + strconv.FormatUint(uint64(p), 10)
}

func (p Port) applyTo(s *Server) error {
	if s.server.Addr == "" {
		s.server.Addr = p.Addr()
		return nil
	}
	if !strings.ContainsRune(s.server.Addr, colon) {
		s.server.Addr += p.Addr()
		return nil
	}

	host, _, err := net.SplitHostPort(s.server.Addr)

	var addrErr net.AddrError
	if errors.As(err, &addrErr) && addrErr.Err == "missing port in address" {
		host = s.server.Addr
	}

	s.server.Addr = host + p.Addr()
	return nil
}
