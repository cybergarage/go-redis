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
	"time"
)

// Database represents a database.
type Database struct {
	ID int
	Records
}

// NewDatabaseWithID returns a new database with the specified ID.
func NewDatabaseWithID(id int) *Database {
	return &Database{
		ID:      id,
		Records: Records{},
	}
}

func (db *Database) GetListRecord(key string) (*Record, List, error) {
	var list List
	record, hasRecord := db.GetRecord(key)
	if hasRecord {
		var ok bool
		list, ok = record.Data.(List)
		if !ok {
			return nil, nil, fmt.Errorf(errorInvalidStoredDataType, record.Data)
		}
	}
	if !hasRecord {
		list = List{}
		record = &Record{
			Key:       key,
			Data:      list,
			Timestamp: time.Now(),
			TTL:       0,
		}
		db.SetRecord(record)
	}
	return record, list, nil
}

func (db *Database) GetSetRecord(key string) (*Record, *Set, error) {
	var set *Set
	record, hasRecord := db.GetRecord(key)
	if hasRecord {
		var ok bool
		set, ok = record.Data.(*Set)
		if !ok {
			return nil, nil, fmt.Errorf(errorInvalidStoredDataType, record.Data)
		}
	}
	if !hasRecord {
		set = NewSet()
		record = &Record{
			Key:       key,
			Data:      set,
			Timestamp: time.Now(),
			TTL:       0,
		}
		db.SetRecord(record)
	}
	return record, set, nil
}
