// Copyright (c) 2021, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func ExampleParsePort() {
	port, err := ParsePort(":8080")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(port)

	// Output:
	// 8080
}

func ExampleSplitHostPort() {
	host, port, err := SplitHostPort("localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(host, port)

	// Output:
	// localhost 8080
}

func TestParsePort(t *testing.T) {
	tests := []struct {
		input    string
		wantPort Port
		wantErr  error
	}{
		{
			input:   "",
			wantErr: ErrMissingPort,
		},
		{
			input:    "443",
			wantPort: 443,
		},
		{
			input:    ":8080",
			wantPort: 8080,
		},
		{
			input:   "localhost:123",
			wantErr: ErrInvalidFormat,
		},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			have, err := ParsePort(tc.input)
			assert.Equal(t, tc.wantPort, have)
			if tc.wantErr == nil {
				assert.Nil(t, err)
			} else {
				// assert.ErrorIs(t, err, ParseError)
				assert.ErrorIs(t, err, tc.wantErr)
			}
		})
	}
}

func TestSplitHostPort(t *testing.T) {
	tests := []struct {
		input    string
		wantHost string
		wantPort Port
		wantErr  error
	}{
		{
			input:   "",
			wantErr: ErrMissingPort,
		},
		{
			input:    ":8080",
			wantPort: 8080,
		},
		{
			input:   "localhost",
			wantErr: ErrMissingPort,
		},
		{
			input:    "localhost:4040",
			wantHost: "localhost",
			wantPort: 4040,
		},
		{
			input:   "[::1]",
			wantErr: ErrMissingPort,
		},
		{
			input:    "[::1]:123",
			wantHost: "::1",
			wantPort: 123,
		},
		{
			input:    "[::1%lo0]:456",
			wantHost: "::1%lo0",
			wantPort: 456,
		},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			haveHost, havePort, err := SplitHostPort(tc.input)
			assert.Equal(t, tc.wantHost, haveHost)
			assert.Equal(t, tc.wantPort, havePort)

			if tc.wantErr == nil {
				assert.Nil(t, err)
			} else {
				// assert.ErrorIs(t, err, newServer(ParseError))
				assert.ErrorIs(t, err, tc.wantErr)
			}
		})
	}
}
