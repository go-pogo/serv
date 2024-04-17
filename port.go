// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"encoding"
	"flag"
	"fmt"
	"github.com/go-pogo/errors"
	"net"
	"strconv"
	"strings"
)

const (
	ErrMissingPort   errors.Msg = "missing port"
	ErrInvalidFormat errors.Msg = "invalid format"
)

type PortParseError struct {
	// Cause is the underlying error. It is never nil.
	Cause error
	// Input string that triggered the error
	Input string
}

func (p *PortParseError) Unwrap() error { return p.Cause }

func (p *PortParseError) Error() string {
	// use an explicit type check because the error might be wrapped
	//goland:noinspection GoTypeAssertionOnErrors
	if e, ok := p.Cause.(errors.Msg); ok {
		return fmt.Sprintf("`%s`: %s", p.Input, e.Error())
	}
	return p.Cause.Error()
}

var (
	_ Option                   = (*Port)(nil)
	_ encoding.TextMarshaler   = (*Port)(nil)
	_ encoding.TextUnmarshaler = (*Port)(nil)
	_ flag.Value               = (*Port)(nil)
)

// Port represents a network port.
type Port uint16

// ParsePort parses string s into a [Port]. A [PortParseError] containing a
// [PortParseError.Cause] is returned when an error is encountered.
func ParsePort(s string) (Port, error) {
	if s == "" {
		return 0, errors.WithStack(&PortParseError{
			Cause: ErrMissingPort,
			Input: s,
		})
	}

	if i := strings.IndexRune(s, ':'); i == 0 {
		s = s[1:]
	} else if i > 0 {
		return 0, errors.WithStack(&PortParseError{
			Cause: ErrInvalidFormat,
			Input: s,
		})
	}

	x, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0, errors.WithStack(&PortParseError{
			Cause: ErrMissingPort,
			Input: s,
		})
	}
	return Port(x), nil
}

// SplitHostPort uses [net.SplitHostPort] to split a network address of the form
// "host:port", "host%zone:port", "[host]:port" or "[host%zone]:port" into host
// or host%zone and [Port]. A [PortParseError] is returned when an error is
// encountered.
func SplitHostPort(hostport string) (string, Port, error) {
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		if isMissingPort(err) {
			err = ErrMissingPort
		}
		return "", 0, errors.WithStack(&PortParseError{
			Cause: err,
			Input: hostport,
		})
	}

	p, err := ParsePort(port)
	return host, p, err
}

// JoinHostPort uses [net.JoinHostPort] to combine host and port into a network
// address of the form "host:port". If host contains a colon, as found in
// literal IPv6 addresses, then [JoinHostPort] returns "[host]:port".
func JoinHostPort(host string, port Port) string {
	return net.JoinHostPort(host, port.String())
}

// Set parses string s into the Port using [ParsePort].
// This method implements the [flag.Value] interface.
func (p *Port) Set(s string) (err error) {
	if s == "" {
		return nil
	}

	*p, err = ParsePort(s)
	return err
}

// UnmarshalText unmarshals text into [Port] using [ParsePort].
// This method implements the [encoding.TextUnmarshaler] interface.
func (p *Port) UnmarshalText(text []byte) (err error) {
	if len(text) == 0 {
		return nil
	}

	*p, err = ParsePort(string(text))
	return err
}

// MarshalText marshals Port into a byte slice using [Port.String].
// This method implements the [encoding.TextMarshaler] interface.
func (p Port) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// String returns the port as a formatted string using [strconv.FormatUint].
func (p Port) String() string {
	if p == 0 {
		return ""
	}
	return strconv.FormatUint(uint64(p), 10)
}

// Addr returns the port as an address string which can be used as value in
// [Server.Addr] or [http.Server.Addr].
func (p Port) Addr() string {
	if p == 0 {
		return ""
	}
	return ":" + strconv.FormatUint(uint64(p), 10)
}

func (p Port) apply(s *Server) error {
	if s.Addr == "" {
		s.Addr = p.Addr()
		return nil
	}
	if !strings.ContainsRune(s.Addr, ':') {
		s.Addr += p.Addr()
		return nil
	}

	host, _, err := net.SplitHostPort(s.Addr)
	if err != nil && isMissingPort(err) {
		host = s.Addr
	}

	s.Addr = host + p.Addr()
	return nil
}

func isMissingPort(err error) bool {
	var addrErr *net.AddrError
	return errors.As(err, &addrErr) && addrErr.Err == "missing port in address"
}
