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
	"time"

	"github.com/cybergarage/go-redis/redis"
)

func (server *Server) Set(ctx *redis.DBContext, key string, val string, opt redis.SetOption) (*redis.Message, error) {
	db := server.GetDatabase(ctx.ID())

	oldVal := ""
	hasOldRecord := false
	if opt.GET {
		var currRecord *Record
		currRecord, hasOldRecord = db.GetRecord(key)
		if hasOldRecord {
			oldVal = string(currRecord.Data)
		}

	}

	record := &Record{
		Key:       key,
		Data:      []byte(val),
		Timestamp: time.Now(),
		TTL:       0,
	}
	db.SetRecord(record)

	if opt.GET {
		if hasOldRecord {
			return redis.NewBulkMessage(oldVal), nil
		}
		return redis.NewNilMessage(), nil
	}
	return redis.NewOKMessage(), nil
}

func (server *Server) Get(ctx *redis.DBContext, key string) (*redis.Message, error) {
	db := server.GetDatabase(ctx.ID())
	record, ok := db.GetRecord(key)
	if !ok {
		return redis.NewNilMessage(), nil
	}
	return redis.NewStringMessage(string(record.Data)), nil
}