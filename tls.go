// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"crypto/tls"
	"github.com/go-pogo/errors"
	"os"
)

// DefaultTLSConfig returns a modern preconfigured tls.Config.
func DefaultTLSConfig() *tls.Config {
	return &tls.Config{
		//PreferServerCipherSuites: true,
		MinVersion: tls.VersionTLS12,

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

type TLSOption interface {
	ApplyTo(conf *tls.Config) error
}

var _ TLSOption = (*TLSConfig)(nil)

type TLSConfig struct {
	CaCertFile string `env:"" flag:"tls-cacert"`
	CertFile   string `env:"" flag:"tls-cert"`
	KeyFile    string `env:"" flag:"tls-key"`

	// todo: implement mtls
	// todo: implement skip verify

	// VerifyClient enables mutual tls authentication.
	VerifyClient bool `env:""`
	// InsecureSkipVerify disabled all certificate verification and should only
	// be used for testing.
	InsecureSkipVerify bool `env:""`
}

func (tc TLSConfig) IsZero() bool {
	return tc.CertFile != "" && tc.KeyFile != ""
}

func (tc TLSConfig) ApplyTo(conf *tls.Config) error {
	if tc.CaCertFile != "" {
		data, err := os.ReadFile(tc.CaCertFile)
		if err != nil {
			return errors.WithStack(err)
		}
		conf.RootCAs.AppendCertsFromPEM(data)
	}
	if tc.InsecureSkipVerify {
		conf.InsecureSkipVerify = tc.InsecureSkipVerify
	}

	return TLSKeyPair{
		CertFile: tc.CaCertFile,
		KeyFile:  tc.KeyFile,
	}.ApplyTo(conf)
}

// CertificateLoader loads a tls.Certificate from any source.
type CertificateLoader interface {
	LoadCertificate() (*tls.Certificate, error)
}

var (
	_ CertificateLoader = (*TLSKeyPair)(nil)
	_ TLSOption         = (*TLSKeyPair)(nil)
)

// TLSKeyPair contains the paths to a public/private key pair of files.
type TLSKeyPair struct {
	CertFile string
	KeyFile  string
}

// LoadCertificate reads and parses the key pair files with tls.LoadX509KeyPair.
// The files must contain PEM encoded data.
func (kp TLSKeyPair) LoadCertificate() (*tls.Certificate, error) {
	if kp.CertFile == "" && kp.KeyFile == "" {
		return nil, nil
	}

	c, err := tls.LoadX509KeyPair(kp.CertFile, kp.KeyFile)
	return &c, errors.WithStack(err)
}

func (kp TLSKeyPair) ApplyTo(conf *tls.Config) error {
	if c, err := kp.LoadCertificate(); err != nil {
		return err
	} else if c != nil {
		conf.Certificates = append(conf.Certificates, *c)
	}

	return nil
}

var (
	_ CertificateLoader = (*TLSPemBlocks)(nil)
	_ TLSOption         = (*TLSPemBlocks)(nil)
)

// certPEMBlock, keyPEMBlock
type TLSPemBlocks struct {
	Cert []byte
	Key  []byte
}

func (pb TLSPemBlocks) LoadCertificate() (*tls.Certificate, error) {
	if len(pb.Cert) == 0 && len(pb.Key) == 0 {
		return nil, nil
	}

	c, err := tls.X509KeyPair(pb.Cert, pb.Key)
	return &c, errors.WithStack(err)
}

func (pb TLSPemBlocks) ApplyTo(conf *tls.Config) error {
	if c, err := pb.LoadCertificate(); err != nil {
		return err
	} else if c != nil {
		conf.Certificates = append(conf.Certificates, *c)
	}
	return nil
}
