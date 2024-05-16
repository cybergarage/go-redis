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
	"os"
)

// NewTLSConfigFrom returns a new TLS configuration from the specified server configuration.
func NewTLSConfigFrom(config *ServerConfig) (*tls.Config, error) {
	certFile, ok := config.ConfigTLSCertFile()
	if !ok {
		return nil, errors.New("no certificate file")
	}
	keyFile, ok := config.ConfigTLSKeyFile()
	if !ok {
		return nil, errors.New("no key file")
	}
	serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	caFile, ok := config.ConfigTLSCaCertFile()
	if ok {
		rootCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, err
		}
		certPool.AppendCertsFromPEM(rootCert)
	}
	return &tls.Config{ // nolint: exhaustruct
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    certPool,
		RootCAs:      certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}, nil
}
