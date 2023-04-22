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
	"sync"
	"time"

	"github.com/cybergarage/go-tracing/tracer"
)

// Conn represents a database connection.
type Conn struct {
	id int
	sync.Map
	ts time.Time
	tracer.SpanContext
}

func newConn() *Conn {
	return &Conn{
		id:          0,
		Map:         sync.Map{},
		ts:          time.Now(),
		SpanContext: nil,
	}
}

// SetDatabase sets th selected database number to the connection.
func (conn *Conn) SetDatabase(id int) {
	conn.id = id
}

// Database returns the current selected database number in the connection.
func (conn *Conn) Database() int {
	return conn.id
}

// Timestamp returns the creation time of the connection.
func (conn *Conn) Timestamp() time.Time {
	return conn.ts
}
