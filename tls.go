// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"crypto/tls"

	"github.com/go-pogo/errors"
)

type TLSConfig = tls.Config

func DefaultTLSConfig() *TLSConfig {
	return &tls.Config{
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,

		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},

		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}

func WithTLS(tc *TLSConfig, cl ...CertificateLoader) Option {
	return optionFunc(func(s *Server) error {
		if tc == nil {
			s.TLSConfig = DefaultTLSConfig()
		} else {
			s.TLSConfig = tc
		}

		var err error
		for _, x := range cl {
			errors.Append(&err, LoadCertificate(s.TLSConfig, x))
		}
		return err
	})
}

type CertificateLoader interface {
	LoadCertificate() (tls.Certificate, error)
}

func LoadCertificate(tc *tls.Config, l CertificateLoader) error {
	cert, err := l.LoadCertificate()
	if err != nil {
		return err
	}

	tc.Certificates = append(tc.Certificates, cert)
	return nil
}

// TLSKeyPair contains the paths to a public/private key pair of files.
type TLSKeyPair [2]string

// LoadCertificate reads and parses the key pair files with tls.LoadX509KeyPair.
// The files must contain PEM encoded data.
func (kp TLSKeyPair) LoadCertificate() (tls.Certificate, error) {
	return tls.LoadX509KeyPair(kp[0], kp[1])
}

type TLSPemBlocks [2][]byte

func (pb TLSPemBlocks) LoadCertificate() (tls.Certificate, error) {
	return tls.X509KeyPair(pb[0], pb[1])
}
