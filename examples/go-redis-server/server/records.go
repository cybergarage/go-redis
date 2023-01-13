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

package server

import (
	"fmt"
	"sync"
)

// Records represents a database record map.
type Records struct {
	sync.Map
}

func NewRecords() *Records {
	return &Records{
		Map: sync.Map{},
	}
}

// Keys returns all key names.
func (rmap *Records) Keys() []string {
	keys := []string{}
	rmap.Range(func(key, value any) bool {
		skey, ok := key.(string)
		if ok {
			keys = append(keys, skey)
		}
		return true
	})
	return keys
}

// SetRecord sets the specified record into the records.
func (rmap *Records) SetRecord(record *Record) error {
	rmap.Store(record.Key, record)
	return nil
}

// HasRecord returns true if the database has the specified key record, otherwise false.
func (rmap *Records) HasRecord(key string) bool {
	_, ok := rmap.Load(key)
	return ok
}

// GetRecord gets a record with the specified key.
func (rmap *Records) GetRecord(key string) (*Record, bool) {
	v, ok := rmap.Load(key)
	if !ok {
		return nil, false
	}
	record, ok := v.(*Record)
	return record, ok
}

// RemoveRecord removes a record with the specified key.
func (rmap *Records) RemoveRecord(key string) error {
	if _, ok := rmap[key]; !ok {
		return fmt.Errorf("%w : %s", ErrNotFound, key)
	}
	delete(rmap, key)
	return nil
}

// RenameRecord renames the specified key record to the specified new record.
func (rmap *Records) RenameRecord(key string, newkey string) error {
	record, ok := rmap.GetRecord(key)
	if !ok {
		return fmt.Errorf("%w: %s", ErrNotFound, key)
	}
	record.Key = newkey
	err := rmap.SetRecord(record)
	if err != nil {
		return err
	}
	err = rmap.RemoveRecord(key)
	if err != nil {
		return err
	}
	return nil
}
