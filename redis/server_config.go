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

import "strconv"

const (
	portConfig    = "port"
	requirePass   = "requirepass"
	tlsPortConfig = "tls-port"
	tlsCertFile   = "tls-cert-file"
	tlsKeyFile    = "tls-key-file"
	tlsCaCertFile = "tls-ca-cert-file"
)

// ServerConfig is a configuration for the Redis server.
type ServerConfig struct {
	*Config
}

// NewDefaultServerConfig returns a default server configuration.
func NewDefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Config: newConfig(),
	}
}

// SetPort sets a listen port number.
func (cfg *ServerConfig) SetPort(port int) {
	cfg.SetConfig(portConfig, strconv.Itoa(port))
}

// ConfigPort returns a listen port number.
func (cfg *ServerConfig) ConfigPort() int {
	port, ok := cfg.ConfigInteger(portConfig)
	if !ok {
		return DefaultTLSPort
	}
	return port
}

// IsPortEnabled returns true if a listen port is enabled.
func (cfg *ServerConfig) IsPortEnabled() bool {
	port, ok := cfg.ConfigInteger(portConfig)
	if !ok {
		return false
	}
	return (0 < port)
}

// SetTLSPort sets a listen port number for TLS.
func (cfg *ServerConfig) SetTLSPort(port int) {
	cfg.SetConfig(tlsPortConfig, strconv.Itoa(port))
}

// ConfigTLSPort returns a listen port number for TLS.
func (cfg *ServerConfig) ConfigTLSPort() int {
	port, ok := cfg.ConfigInteger(tlsPortConfig)
	if !ok {
		return DefaultTLSPort
	}
	return port
}

// IsTLSPortEnabled returns true if a listen port for TLS is enabled.
func (cfg *ServerConfig) IsTLSPortEnabled() bool {
	port, ok := cfg.ConfigInteger(tlsPortConfig)
	if !ok {
		return false
	}
	return (0 < port)
}

// SetTLSCertFile sets a certificate file.
func (cfg *ServerConfig) SetTLSCertFile(certFile string) {
	cfg.SetConfig(tlsCertFile, certFile)
}

// ConfigTLSCertFile returns a certificate file.
func (cfg *ServerConfig) ConfigTLSCertFile() (string, bool) {
	return cfg.ConfigString(tlsCertFile)
}

// SetTLSKeyFile sets a key file.
func (cfg *ServerConfig) SetTLSKeyFile(keyFile string) {
	cfg.SetConfig(tlsKeyFile, keyFile)
}

// ConfigTLSKeyFile returns a key file.
func (cfg *ServerConfig) ConfigTLSKeyFile() (string, bool) {
	return cfg.ConfigString(tlsKeyFile)
}

// SetTLSCaCertFile sets a CA certificate file.
func (cfg *ServerConfig) SetTLSCaCertFile(caCertFile string) {
	cfg.SetConfig(tlsCaCertFile, caCertFile)
}

// ConfigTLSCaCertFile returns a CA certificate file.
func (cfg *ServerConfig) ConfigTLSCaCertFile() (string, bool) {
	return cfg.ConfigString(tlsCaCertFile)
}

// SetRequirePass sets a password.
func (cfg *ServerConfig) SetRequirePass(password string) {
	cfg.SetConfig(requirePass, password)
}

// ConfigRequirePass returns a password.
func (cfg *ServerConfig) ConfigRequirePass() (string, bool) {
	return cfg.ConfigString(requirePass)
}

// RemoveRequirePass removes a password.
func (cfg *ServerConfig) RemoveRequirePass() {
	cfg.RemoveConfig(requirePass)
}
