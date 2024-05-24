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
	"net"
	"sync"
	"time"

	"github.com/cybergarage/go-tracing/tracer"
)

// Conn represents a database connection.
type Conn struct {
	net.Conn
	id        DatabaseID
	authrized bool
	sync.Map
	ts time.Time
	tracer.Context
	tlsState *tls.ConnectionState
	username string
	password string
}

func newConnWith(conn net.Conn, tlsState *tls.ConnectionState) *Conn {
	return &Conn{
		Conn:      conn,
		authrized: false,
		id:        0,
		Map:       sync.Map{},
		ts:        time.Now(),
		Context:   nil,
		tlsState:  tlsState,
	}
}

// SetDatabase sets the selected database number to the connection.
func (conn *Conn) SetDatabase(id DatabaseID) {
	conn.id = id
}

// Database returns the current selected database number in the connection.
func (conn *Conn) Database() DatabaseID {
	return conn.id
}

// SetAuthrized sets the authrized flag to the connection.
func (conn *Conn) SetAuthrized(authrized bool) {
	conn.authrized = authrized
}

// IsAuthrized returns true if the connection is authrized.
func (conn *Conn) IsAuthrized() bool {
	return conn.authrized
}

// SetUserName sets the user name to the connection.
func (conn *Conn) SetUserName(username string) {
	conn.username = username
}

// UserName returns the user name and true if the connection has the user name.
func (conn *Conn) UserName() (string, bool) {
	return conn.username, 0 < len(conn.username)
}

// SetPassword sets the password to the connection.
func (conn *Conn) SetPassword(password string) {
	conn.password = password
}

// Password returns the password and true if the connection has the password.
func (conn *Conn) Password() (string, bool) {
	return conn.password, 0 < len(conn.password)
}

// Timestamp returns the creation time of the connection.
func (conn *Conn) Timestamp() time.Time {
	return conn.ts
}

// SetAuthrizedxt to the connection.
func (conn *Conn) SetSpanContext(span tracer.Context) {
	conn.Context = span
}

// SpanContext returns the span context of the connection.
func (conn *Conn) SpanContext() tracer.Context {
	return conn.Context
}

// IsTLSConnection return true if the connection is enabled TLS.
func (conn *Conn) IsTLSConnection() bool {
	return conn.tlsState != nil
}

// TLSConnectionState returns the TLS connection state.
func (conn *Conn) TLSConnectionState() (*tls.ConnectionState, bool) {
	return conn.tlsState, conn.tlsState != nil
}
