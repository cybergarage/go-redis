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
	"github.com/cybergarage/go-authenticator/auth"
)

// TLSConfig represents a server TLS configuration.
type TLSConfig interface {
	auth.CertConfig

	// SetTLSPort sets a listen port number for TLS.
	SetTLSPort(port int)
	// TLSPort returns a listen port number for TLS.
	TLSPort() int
	// IsTLSPortEnabled returns true if a listen port for TLS is enabled.
	IsTLSPortEnabled() bool
}

// Config represents a server configuration.
type Config interface {
	TLSConfig

	// SetPort sets a listen port number.
	SetPort(port int)
	// Port returns a listen port number.
	Port() int
	// IsPortEnabled returns true if a listen port is enabled.
	IsPortEnabled() bool

	// SetRequirePass sets a password.
	SetRequirePass(password string)
	// ConfigRequirePass returns a password.
	ConfigRequirePass() (string, bool)
	// RemoveRequirePass removes a password.
	RemoveRequirePass()
}
