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

package server

import "fmt"

// Records represents a database record map.
type Records map[string]*Record

// Keys returns all key names.
func (rmap Records) Keys(record *Record) []string {
	keys := []string{}
	for key := range rmap {
		keys = append(keys, key)
	}
	return keys
}

// SetRecord sets the specified record into the records.
func (rmap Records) SetRecord(record *Record) error {
	rmap[record.Key] = record
	return nil
}

// GetRecord gets a record with the specified key.
func (rmap Records) GetRecord(key string) (*Record, bool) {
	record, ok := rmap[key]
	return record, ok
}

// RemoveRecord removes a record with the specified key.
func (rmap Records) RemoveRecord(key string) error {
	if _, ok := rmap[key]; !ok {
		return fmt.Errorf("%w : %s", ErrNotFound, key)
	}
	delete(rmap, key)
	return nil
}
