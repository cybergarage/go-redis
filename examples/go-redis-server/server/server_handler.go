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

func (server *Server) Del(ctx *redis.DBContext, keys []string) (*redis.Message, error) {
	db := server.GetDatabase(ctx.ID())
	removedCount := 0
	for _, key := range keys {
		err := db.RemoveRecord(key)
		if err == nil {
			removedCount++
		}
	}
	return redis.NewIntegerMessage(removedCount), nil
}

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

// nolint: ifshort
func (server *Server) HSet(ctx *redis.DBContext, key string, field string, val string, opt redis.HSetOption) (*redis.Message, error) {
	db := server.GetDatabase(ctx.ID())
	var dict DictionaryRecord
	record, hasRecord := db.GetRecord(key)
	if hasRecord {
		var ok bool
		dict, ok = record.Data.(DictionaryRecord)
		if !ok {
			hasRecord = false
		}
	}
	if !hasRecord {
		dict := DictionaryRecord{}
		dict[field] = val
		record := &Record{
			Key:       key,
			Data:      dict,
			Timestamp: time.Now(),
			TTL:       0,
		}
		db.SetRecord(record)
		return redis.NewIntegerMessage(1), nil
	}

	_, hasKey := dict[field]
	dict[field] = val
	if hasKey {
		return redis.NewIntegerMessage(0), nil
	}
	return redis.NewIntegerMessage(1), nil
}

func (server *Server) HGet(ctx *redis.DBContext, key string, field string, opt redis.HGetOption) (*redis.Message, error) {
	db := server.GetDatabase(ctx.ID())
	record, ok := db.GetRecord(key)
	if !ok {
		return redis.NewNilMessage(), nil
	}
	dict, ok := record.Data.(DictionaryRecord)
	if !ok {
		return redis.NewNilMessage(), nil
	}
	dictData, ok := dict[field]
	if !ok {
		return redis.NewNilMessage(), nil
	}
	return redis.NewStringMessage(dictData), nil
}

func (server *Server) HGetAll(ctx *redis.DBContext, key string) (*redis.Message, error) {
	arrayMsg := redis.NewArrayMessage()

	db := server.GetDatabase(ctx.ID())
	record, ok := db.GetRecord(key)
	if !ok {
		return arrayMsg, nil
	}

	dict, ok := record.Data.(DictionaryRecord)
	if !ok {
		return arrayMsg, nil
	}

	array, _ := arrayMsg.Array()
	for key, val := range dict {
		array.Append(redis.NewBulkMessage(key))
		array.Append(redis.NewBulkMessage(val))
	}

	return arrayMsg, nil
}

func (server *Server) MSet(ctx *redis.DBContext, dict map[string]string, opt redis.MSetOption) (*redis.Message, error) {
	if opt.NX {
		getOpt := redis.GetOption{}
		for key := range dict {
			res, err := server.Get(ctx, key, getOpt)
			if err != nil {
				return nil, err
			}
			if !res.IsNil() {
				return redis.NewIntegerMessage(0), nil
			}
		}
	}

	now := time.Now()
	setOpt := redis.SetOption{
		NX:      true,
		XX:      false,
		EX:      0,
		PX:      0,
		EXAT:    now,
		PXAT:    now,
		KEEPTTL: false,
		GET:     false,
	}
	for key, val := range dict {
		if _, err := server.Set(ctx, key, val, setOpt); err != nil {
			return nil, err
		}
	}

	if opt.NX {
		return redis.NewIntegerMessage(1), nil
	}
	return redis.NewOKMessage(), nil
}

func (server *Server) MGet(ctx *redis.DBContext, keys []string, opt redis.MGetOption) (*redis.Message, error) {
	getOpt := redis.GetOption{}
	arrayMsg := redis.NewArrayMessage()
	array, _ := arrayMsg.Array()
	for _, key := range keys {
		msg, err := server.Get(ctx, key, getOpt)
		if err != nil {
			return nil, err
		}
		array.Append(msg)
	}
	return arrayMsg, nil
}

func (server *Server) HMSet(ctx *redis.DBContext, key string, dict map[string]string, opt redis.HMSetOption) (*redis.Message, error) {
	hsetOpt := redis.HSetOption{}
	for field, val := range dict {
		if _, err := server.HSet(ctx, key, field, val, hsetOpt); err != nil {
			return nil, err
		}
	}
	return redis.NewOKMessage(), nil
}

func (server *Server) HMGet(ctx *redis.DBContext, key string, fields []string, opt redis.HMGetOption) (*redis.Message, error) {
	hgetOpt := redis.HGetOption{}
	arrayMsg := redis.NewArrayMessage()
	array, _ := arrayMsg.Array()
	for _, field := range fields {
		msg, err := server.HGet(ctx, key, field, hgetOpt)
		if err != nil {
			return nil, err
		}
		array.Append(msg)
	}
	return arrayMsg, nil
}
