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
)

// o represents a server TLS configuration.
type TLSConfig interface {
	// SetTLSPort sets a listen port number for TLS.
	SetTLSPort(port int)
	// TLSPort returns a listen port number for TLS.
	TLSPort() int
	// IsTLSPortEnabled returns true if a listen port for TLS is enabled.
	IsTLSPortEnabled() bool

	// SetTLSCertFile sets a certificate file.
	SetTLSCertFile(certFile string) error
	// ConfigTLSCertFile returns a certificate file.
	ConfigTLSCertFile() (string, bool)
	// ConfigTLSCert returns a certificate.
	ConfigTLSCert() ([]byte, bool)

	// SetTLSKeyFile sets a key file.
	SetTLSKeyFile(keyFile string) error
	// ConfigTLSKeyFile returns a key file.
	ConfigTLSKeyFile() (string, bool)
	// ConfigTLSKey returns a key.
	ConfigTLSKey() ([]byte, bool)

	// SetTLSCaCertFile sets a CA certificate file.
	SetTLSCaCertFile(caCertFile string) error
	// ConfigTLSCACertFile returns a CA certificate file.
	ConfigTLSCACertFile() (string, bool)
	// ConfigTLSCACert returns a CA certificate.
	ConfigTLSCACert() ([]byte, bool)

	// SetTLSConfig sets a TLS configuration.
	SetTLSConfig(config *tls.Config)
	// TLSConfig returns a TLS configuration.
	TLSConfig() (*tls.Config, bool)
}

// Config represents a server configuration.
type Config interface {
	TLSConfig

	// SetPort sets a listen port number.
	SetPort(port int)
	// ConfigPort returns a listen port number.
	ConfigPort() int
	// IsPortEnabled returns true if a listen port is enabled.
	IsPortEnabled() bool

	// SetRequirePass sets a password.
	SetRequirePass(password string)
	// ConfigRequirePass returns a password.
	ConfigRequirePass() (string, bool)
	// RemoveRequirePass removes a password.
	RemoveRequirePass()
}
