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
	"errors"
	"sync"

	"github.com/google/uuid"
)

// ConnManager represents a connection map.
type ConnManager struct {
	m     map[uuid.UUID]*Conn
	mutex *sync.RWMutex
}

// NewConnManager returns a connection map.
func NewConnManager() *ConnManager {
	return &ConnManager{
		m:     map[uuid.UUID]*Conn{},
		mutex: &sync.RWMutex{},
	}
}

// AddConn adds the specified connection.
func (mgr *ConnManager) AddConn(c *Conn) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	uuid := c.UUID()
	mgr.m[uuid] = c
}

// ConnByUUID returns the connection by the specified UUID.
func (mgr *ConnManager) ConnByUUID(uuid uuid.UUID) (*Conn, bool) {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	c, ok := mgr.m[uuid]
	return c, ok
}

// DeleteConnByUID deletes the specified connection by the connection ID.
func (mgr *ConnManager) DeleteConnByUID(uuid uuid.UUID) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	delete(mgr.m, uuid)
}

// Start starts the connection manager.
func (mgr *ConnManager) Start() error {
	return nil
}

// Stop closes all connections.
func (mgr *ConnManager) Stop() error {
	var errs error
	for _, conn := range mgr.m {
		err := conn.Close()
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

// Length returns the included connection count.
func (mgr *ConnManager) Length() int {
	mgr.mutex.RLock()
	defer mgr.mutex.RUnlock()
	return len(mgr.m)
}
