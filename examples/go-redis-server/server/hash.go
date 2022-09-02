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

type Hash map[string]string

func (hash Hash) Del(fields []string) int {
	removedFields := 0
	for _, field := range fields {
		_, ok := hash[field]
		if !ok {
			continue
		}
		delete(hash, field)
		removedFields++
	}
	return removedFields
}

func (server *Server) HDel(ctx *redis.DBContext, key string, fields []string) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}
	record, ok := db.GetRecord(key)
	if !ok {
		return redis.NewIntegerMessage(0), nil
	}
	dict, ok := record.Data.(Hash)
	if !ok {
		return redis.NewIntegerMessage(0), nil
	}
	return redis.NewIntegerMessage(dict.Del(fields)), nil
}

// nolint: ifshort
func (server *Server) HSet(ctx *redis.DBContext, key string, field string, val string, opt redis.HSetOption) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}

	var dict Hash
	record, hasRecord := db.GetRecord(key)
	if hasRecord {
		var ok bool
		dict, ok = record.Data.(Hash)
		if !ok {
			hasRecord = false
		}
	}
	if !hasRecord {
		dict := Hash{}
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
	if opt.NX && hasKey {
		return redis.NewIntegerMessage(0), nil
	}

	dict[field] = val
	if hasKey {
		return redis.NewIntegerMessage(0), nil
	}
	return redis.NewIntegerMessage(1), nil
}

func (server *Server) HGet(ctx *redis.DBContext, key string, field string) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}
	record, ok := db.GetRecord(key)
	if !ok {
		return redis.NewNilMessage(), nil
	}
	dict, ok := record.Data.(Hash)
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

	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}
	record, ok := db.GetRecord(key)
	if !ok {
		return arrayMsg, nil
	}

	dict, ok := record.Data.(Hash)
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

func (server *Server) HMSet(ctx *redis.DBContext, key string, dict map[string]string) (*redis.Message, error) {
	hsetOpt := redis.HSetOption{
		NX: false,
	}
	for field, val := range dict {
		if _, err := server.HSet(ctx, key, field, val, hsetOpt); err != nil {
			return nil, err
		}
	}
	return redis.NewOKMessage(), nil
}

func (server *Server) HMGet(ctx *redis.DBContext, key string, fields []string) (*redis.Message, error) {
	arrayMsg := redis.NewArrayMessage()
	array, _ := arrayMsg.Array()
	for _, field := range fields {
		msg, err := server.HGet(ctx, key, field)
		if err != nil {
			return nil, err
		}
		array.Append(msg)
	}
	return arrayMsg, nil
}
