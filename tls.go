// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serv

import (
	"crypto/tls"
	"os"

	"github.com/go-pogo/errors"
)

// DefaultTLSConfig returns a modern preconfigured tls.Config.
func DefaultTLSConfig() *tls.Config {
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

type TLSOption interface {
	Apply(conf *tls.Config) error
}

func WithTLS(conf *tls.Config, opts ...TLSOption) Option {
	return optionFunc(func(s *Server) error {
		if conf == nil {
			s.TLSConfig = DefaultTLSConfig()
		} else {
			s.TLSConfig = conf
		}

		var err error
		for _, opt := range opts {
			errors.Append(&err, opt.Apply(conf))
		}
		return err
	})
}

var _ TLSOption = &TLSConfig{}

type TLSConfig struct {
	CaCertFile string `env:"TLS_CA_FILE" flag:"tlscacert"`
	CertFile   string `env:"TLS_CERT_FILE" flag:"tlscert"`
	KeyFile    string `env:"TLS_KEY_FILE" flag:"tlskey"`

	// InsecureSkipVerify should be used only for testing
	InsecureSkipVerify bool `env:"INSECURE_SKIP_VERIFY"`
}

func (tc TLSConfig) IsZero() bool {
	return tc.CertFile != "" && tc.KeyFile != ""
}

func (tc TLSConfig) Apply(conf *tls.Config) error {
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
	}.Apply(conf)
}

// CertificateLoader loads a tls.Certificate from any source.
type CertificateLoader interface {
	LoadCertificate() (*tls.Certificate, error)
}

var _ CertificateLoader = &TLSKeyPair{}
var _ TLSOption = &TLSKeyPair{}

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

func (kp TLSKeyPair) Apply(conf *tls.Config) error {
	if c, err := kp.LoadCertificate(); err != nil {
		return err
	} else if c != nil {
		conf.Certificates = append(conf.Certificates, *c)
	}

	return nil
}

var _ CertificateLoader = &TLSPemBlocks{}
var _ TLSOption = &TLSPemBlocks{}

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

func (pb TLSPemBlocks) Apply(conf *tls.Config) error {
	if c, err := pb.LoadCertificate(); err != nil {
		return err
	} else if c != nil {
		conf.Certificates = append(conf.Certificates, *c)
	}
	return nil
}
