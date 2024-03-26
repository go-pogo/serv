// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/go-pogo/errors"
	"os"
)

const ErrAppendRootCAFailure errors.Msg = "failed to append certificate to root ca pool"

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

var (
	_ Option    = (*TLSConfig)(nil)
	_ TLSOption = (*TLSConfig)(nil)
)

type TLSConfig struct {
	CaCertFile string `env:"" flag:"tls-cacert"`
	CertFile   string `env:"" flag:"tls-cert"`
	KeyFile    string `env:"" flag:"tls-key"`

	// todo: implement mtls
	// VerifyClient enables mutual tls authentication.
	VerifyClient bool `env:""`
	// InsecureSkipVerify disabled all certificate verification and should only
	// be used for testing. See tls.Config for additional information.
	InsecureSkipVerify bool `env:""`
}

func (tc TLSConfig) ApplyTo(conf *tls.Config) error {
	if tc.CaCertFile != "" {
		data, err := os.ReadFile(tc.CaCertFile)
		if err != nil {
			return errors.WithStack(err)
		}
		if conf.RootCAs == nil {
			if conf.RootCAs, err = x509.SystemCertPool(); err != nil {
				return errors.WithStack(err)
			}
		}
		if !conf.RootCAs.AppendCertsFromPEM(data) {
			return errors.New(ErrAppendRootCAFailure)
		}
	}

	conf.InsecureSkipVerify = tc.InsecureSkipVerify

	return TLSKeyPair{
		CertFile: tc.CertFile,
		KeyFile:  tc.KeyFile,
	}.ApplyTo(conf)
}

func (tc TLSConfig) apply(s *Server) error {
	if tc.CertFile == "" || tc.KeyFile == "" {
		return nil
	}

	s.TLSConfig = DefaultTLSConfig()
	return tc.ApplyTo(s.TLSConfig)
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
