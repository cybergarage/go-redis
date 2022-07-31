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

import (
	"time"

	"github.com/cybergarage/go-redis/redis"
)

func (server *Server) Set(ctx *redis.DBContext, key string, val string, opt redis.SetOption) (*redis.Message, error) {
	db := server.GetDatabase(ctx.ID())

	var oldVal []byte = nil
	hasOldRecord := false
	if opt.NX || opt.GET {
		var currRecord *Record
		currRecord, hasOldRecord = db.GetRecord(key)
		switch {
		case opt.NX:
			if hasOldRecord {
				return redis.NewIntegerMessage(0), nil
			}
		case opt.GET:
			if hasOldRecord {
				stringData, ok := currRecord.Data.(string)
				if ok {
					oldVal = []byte(stringData)
				}
			}
		}

	}

	record := &Record{
		Key:       key,
		Data:      val,
		Timestamp: time.Now(),
		TTL:       0,
	}
	db.SetRecord(record)

	switch {
	case opt.NX:
		return redis.NewIntegerMessage(1), nil
	case opt.GET:
		if hasOldRecord && oldVal != nil {
			return redis.NewBulkMessage(string(oldVal)), nil
		}
		return redis.NewNilMessage(), nil
	}

	return redis.NewOKMessage(), nil
}

func (server *Server) Get(ctx *redis.DBContext, key string, opt redis.GetOption) (*redis.Message, error) {
	db := server.GetDatabase(ctx.ID())
	record, ok := db.GetRecord(key)
	if !ok {
		return redis.NewNilMessage(), nil
	}
	stringData, ok := record.Data.(string)
	if ok {
		return redis.NewStringMessage(stringData), nil
	}
	return redis.NewNilMessage(), nil
}

func (server *Server) HSet(ctx *redis.DBContext, hash string, key string, val string, opt redis.HSetOption) (*redis.Message, error) {
	return nil, nil
}

func (server *Server) HGet(ctx *redis.DBContext, hash string, key string, opt redis.HGetOption) (*redis.Message, error) {
	db := server.GetDatabase(ctx.ID())
	record, ok := db.GetRecord(hash)
	if !ok {
		return redis.NewNilMessage(), nil
	}
	dict, ok := record.Data.(DictionaryRecord)
	if !ok {
		return redis.NewNilMessage(), nil
	}
	dictData, ok := dict[key]
	if !ok {
		return redis.NewNilMessage(), nil
	}
	return redis.NewStringMessage(dictData), nil
}
