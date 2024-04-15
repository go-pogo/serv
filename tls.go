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
	// CACertFile is the path to the root certificate authority file. It is used
	// to verify the client's (whom connect to the server) certificate.
	CACertFile string `env:"" flag:"tls-ca"`
	// CertFile is the path to the server's certificate file.
	CertFile string `env:"" flag:"tls-cert"`
	// KeyFile is the path to the server's private key file.
	KeyFile string `env:"" flag:"tls-key"`

	// VerifyClient enables mutual tls authentication.
	VerifyClient bool `env:""`
	// InsecureSkipVerify disabled all certificate verification and should only
	// be used for testing. See tls.Config for additional information.
	InsecureSkipVerify bool `env:""`
}

func (tc TLSConfig) ApplyTo(conf *tls.Config) error {
	if conf == nil {
		return nil
	}

	if tc.CACertFile != "" {
		data, err := os.ReadFile(tc.CACertFile)
		if err != nil {
			return errors.WithStack(err)
		}
		if conf.ClientCAs == nil {
			conf.ClientCAs = x509.NewCertPool()
		}
		if !conf.ClientCAs.AppendCertsFromPEM(data) {
			return errors.New(ErrAppendRootCAFailure)
		}
	}

	conf.InsecureSkipVerify = tc.InsecureSkipVerify
	if tc.VerifyClient {
		conf.ClientAuth = tls.RequireAndVerifyClientCert
	}

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

// GetCertificate can be used in tls.Config to load a certificate when it's
// requested for.
func GetCertificate(cl CertificateLoader) func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		cert, err := cl.LoadCertificate()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return cert, nil
	}
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
	if conf == nil {
		return nil
	}
	if conf.GetCertificate == nil {
		conf.GetCertificate = GetCertificate(kp)
		return nil
	}

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
	if conf == nil {
		return nil
	}
	if c, err := pb.LoadCertificate(); err != nil {
		return err
	} else if c != nil {
		conf.Certificates = append(conf.Certificates, *c)
	}
	return nil
}
