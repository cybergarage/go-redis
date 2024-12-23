// Copyright (C) 2022 The go-redis Authors All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redis

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
)

// NewTLSConfigFrom returns a new TLS configuration from the specified server configuration.
func NewTLSConfigFrom(config Config) (*tls.Config, error) {
	tlsConfig, ok := config.TLSConfig()
	if ok {
		return tlsConfig, nil
	}
	cert, ok := config.ConfigTLSCert()
	if !ok {
		return nil, errors.New("no server certificate")
	}
	key, ok := config.ConfigTLSKey()
	if !ok {
		return nil, errors.New("no server key")
	}
	serverCert, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	caCert, ok := config.ConfigTLSCACert()
	if ok {
		certPool.AppendCertsFromPEM(caCert)
	}
	return &tls.Config{ // nolint: exhaustruct
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    certPool,
		RootCAs:      certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}, nil
}
