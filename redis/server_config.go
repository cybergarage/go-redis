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
	portConfig  = "port"
	requirePass = "requirepass"
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
	portStr, ok := cfg.ConfigParameter(portConfig)
	if !ok {
		return DefaultPort
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return DefaultPort
	}
	return port
}

// SetRequirePass sets a password.
func (cfg *ServerConfig) SetRequirePass(password string) {
	cfg.SetConfig(requirePass, password)
}

// ConfigRequirePass returns a password.
func (cfg *ServerConfig) ConfigRequirePass() (bool, string) {
	passwd, ok := cfg.ConfigParameter(requirePass)
	if !ok {
		return false, ""
	}
	return true, passwd
}

// RemoveRequirePass removes a password.
func (cfg *ServerConfig) RemoveRequirePass() {
	cfg.RemoveConfig(requirePass)
}
