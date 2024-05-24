// Copyright (C) 2019 The go-redis Authors. All rights reserved.
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

package auth

import (
	"crypto/tls"
	"net"
)

// Conn represents a connection.
type Conn interface {
	net.Conn

	// UserName returns the user name and true if the connection has the user name.
	UserName() (string, bool)
	// Password returns the password and true if the connection has the password.
	Password() (string, bool)

	// IsTLSConnection return true if the connection is enabled TLS.
	IsTLSConnection() bool
	// TLSConnectionState returns the TLS connection state.
	TLSConnectionState() (*tls.ConnectionState, bool)
}
