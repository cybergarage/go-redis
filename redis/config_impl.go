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
	"os"
	"strconv"
)

const (
	portConfig    = "port"
	requirePass   = "requirepass"
	tlsPortConfig = "tls-port"
	tlsCertFile   = "tls-cert-file"
	tlsKeyFile    = "tls-key-file"
	tlsCACertFile = "tls-ca-cert-file"
)

// serverConfig is a configuration for the Redis server.
type serverConfig struct {
	*configMap
	ServerCert []byte
	ServerKey  []byte
	CACerts    []byte
	tlsConfig  *tls.Config
}

// newDefaultServerConfig returns a default server configuration.
func newDefaultServerConfig() *serverConfig {
	return &serverConfig{
		configMap:  newConfig(),
		ServerCert: nil,
		ServerKey:  nil,
		CACerts:    nil,
		tlsConfig:  nil,
	}
}

// SetPort sets a listen port number.
func (cfg *serverConfig) SetPort(port int) {
	cfg.SetConfig(portConfig, strconv.Itoa(port))
}

// ConfigPort returns a listen port number.
func (cfg *serverConfig) ConfigPort() int {
	port, ok := cfg.ConfigInteger(portConfig)
	if !ok {
		return DefaultTLSPort
	}
	return port
}

// IsPortEnabled returns true if a listen port is enabled.
func (cfg *serverConfig) IsPortEnabled() bool {
	port, ok := cfg.ConfigInteger(portConfig)
	if !ok {
		return false
	}
	return (0 < port)
}

// SetTLSPort sets a listen port number for TLS.
func (cfg *serverConfig) SetTLSPort(port int) {
	cfg.SetConfig(tlsPortConfig, strconv.Itoa(port))
}

// TLSPort returns a listen port number for TLS.
func (cfg *serverConfig) TLSPort() int {
	port, ok := cfg.ConfigInteger(tlsPortConfig)
	if !ok {
		return DefaultTLSPort
	}
	return port
}

// IsTLSPortEnabled returns true if a listen port for TLS is enabled.
func (cfg *serverConfig) IsTLSPortEnabled() bool {
	port, ok := cfg.ConfigInteger(tlsPortConfig)
	if !ok {
		return false
	}
	return (0 < port)
}

// SetTLSCertFile sets a certificate file.
func (cfg *serverConfig) SetTLSCertFile(certFile string) error {
	cert, err := os.ReadFile(certFile)
	if err != nil {
		return err
	}
	cfg.SetConfig(tlsCertFile, certFile)
	cfg.ServerCert = cert
	return nil
}

// ConfigTLSCertFile returns a certificate file.
func (cfg *serverConfig) ConfigTLSCertFile() (string, bool) {
	return cfg.ConfigString(tlsCertFile)
}

// ConfigTLSCert returns a certificate.
func (cfg *serverConfig) ConfigTLSCert() ([]byte, bool) {
	if cfg.ServerCert == nil {
		return nil, false
	}
	return cfg.ServerCert, true
}

// SetTLSKeyFile sets a key file.
func (cfg *serverConfig) SetTLSKeyFile(keyFile string) error {
	key, err := os.ReadFile(keyFile)
	if err != nil {
		return err
	}
	cfg.SetConfig(tlsKeyFile, keyFile)
	cfg.ServerKey = key
	return nil
}

// ConfigTLSKeyFile returns a key file.
func (cfg *serverConfig) ConfigTLSKeyFile() (string, bool) {
	return cfg.ConfigString(tlsKeyFile)
}

// ConfigTLSKey returns a key.
func (cfg *serverConfig) ConfigTLSKey() ([]byte, bool) {
	if cfg.ServerKey == nil {
		return nil, false
	}
	return cfg.ServerKey, true
}

// SetTLSCaCertFile sets a CA certificate file.
func (cfg *serverConfig) SetTLSCaCertFile(caCertFile string) error {
	rootCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return err
	}
	cfg.SetConfig(tlsCACertFile, caCertFile)
	cfg.CACerts = rootCert
	return nil
}

// ConfigTLSCACertFile returns a CA certificate file.
func (cfg *serverConfig) ConfigTLSCACertFile() (string, bool) {
	return cfg.ConfigString(tlsCACertFile)
}

// ConfigTLSCACert returns a CA certificate.
func (cfg *serverConfig) ConfigTLSCACert() ([]byte, bool) {
	if cfg.CACerts == nil {
		return nil, false
	}
	return cfg.CACerts, true
}

// SetTLSConfig sets a TLS configuration.
func (cfg *serverConfig) SetTLSConfig(tlsConfig *tls.Config) {
	cfg.tlsConfig = tlsConfig
}

// TLSConfig returns a TLS configuration.
func (cfg *serverConfig) TLSConfig() (*tls.Config, bool) {
	if cfg.tlsConfig == nil {
		return nil, false
	}
	return cfg.tlsConfig, true
}

// SetRequirePass sets a password.
func (cfg *serverConfig) SetRequirePass(password string) {
	cfg.SetConfig(requirePass, password)
}

// ConfigRequirePass returns a password.
func (cfg *serverConfig) ConfigRequirePass() (string, bool) {
	return cfg.ConfigString(requirePass)
}

// RemoveRequirePass removes a password.
func (cfg *serverConfig) RemoveRequirePass() {
	cfg.RemoveConfig(requirePass)
}
