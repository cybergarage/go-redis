// Copyright (C) 2022 Satoshi Konno All rights reserved.
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
	"net"
	"sync"
	"time"

	"github.com/cybergarage/go-tracing/tracer"
)

// Conn represents a database connection.
type Conn struct {
	net.Conn
	id        int
	authrized bool
	sync.Map
	ts time.Time
	tracer.Context
}

func newConnWith(conn net.Conn) *Conn {
	return &Conn{
		Conn:      conn,
		authrized: false,
		id:        0,
		Map:       sync.Map{},
		ts:        time.Now(),
		Context:   nil,
	}
}

// SetDatabase sets the selected database number to the connection.
func (conn *Conn) SetDatabase(id int) {
	conn.id = id
}

// Database returns the current selected database number in the connection.
func (conn *Conn) Database() int {
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
