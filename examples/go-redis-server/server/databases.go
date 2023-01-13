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

import "sync"

// Databases represents a database map.
type Databases struct {
	sync.Map
}

func NewDatabases() *Databases {
	return &Databases{
		Map: sync.Map{},
	}
}

// GetDatabase returns the database with the specified ID.
func (dbs Databases) GetDatabase(id int) (*Database, bool) {
	db, ok := dbs.Load(id)
	if !ok {
		return nil, false
	}
	return db.(*Database), true
}
