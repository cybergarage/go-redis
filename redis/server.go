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

	"github.com/cybergarage/go-authenticator/auth"
)

type ServerConf interface {
	// SetPort sets a listen port number.
	SetPort(port int)
	// ConfigPort returns a listen port number.
	ConfigPort() int
	// IsPortEnabled returns true if a listen port is enabled.
	IsPortEnabled() bool
	// SetTLSPort sets a listen port number for TLS.
	SetTLSPort(port int)
	// ConfigTLSPort returns a listen port number for TLS.
	ConfigTLSPort() int
	// IsTLSPortEnabled returns true if a listen port for TLS is enabled.
	IsTLSPortEnabled() bool

	// SetTLSCertFile sets a certificate file.
	SetTLSCertFile(certFile string) error
	// ConfigTLSCertFile returns a certificate file.
	// ConfigTLSCert returns a certificate.
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
	// ConfigTLSConfig returns a TLS configuration.
	ConfigTLSConfig() (*tls.Config, bool)

	// SetRequirePass sets a password.
	SetRequirePass(password string)

	// ConfigRequirePass returns a password.
	ConfigRequirePass() (string, bool)

	// RemoveRequirePass removes a password.
	RemoveRequirePass()
}

type Server interface {
	auth.Manager
	ServerConf

	// SetAuthCommandHandler sets a auth handler to handle auth commands.
	SetAuthCommandHandler(handler AuthCommandHandler)
	// SetCommandHandler sets a user handler to handle user commands.
	SetCommandHandler(handler UserCommandHandler)
	// RegisterExexutor sets a command executor.
	RegisterExexutor(cmd string, executor Executor)

	Start() error
	Stop() error
	Restart() error
}
