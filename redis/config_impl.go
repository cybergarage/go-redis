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
	"strconv"

	"github.com/cybergarage/go-authenticator/auth/tls"
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
	tls.CertConfig
}

// newDefaultServerConfig returns a default server configuration.
func newDefaultServerConfig() *serverConfig {
	return &serverConfig{
		configMap:  newConfig(),
		CertConfig: tls.NewCertConfig()}
}

// SetPort sets a listen port number.
func (cfg *serverConfig) SetPort(port int) {
	cfg.SetConfig(portConfig, strconv.Itoa(port))
}

// Port returns a listen port number.
func (cfg *serverConfig) Port() int {
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
