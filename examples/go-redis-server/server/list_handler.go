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

type List []string

func (server *Server) getDatabaseListRecord(ctx *redis.DBContext, key string) (*Record, List, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, nil, err
	}

	var list List
	record, hasRecord := db.GetRecord(key)
	if hasRecord {
		var ok bool
		list, ok = record.Data.(List)
		if !ok {
			hasRecord = false
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

func (server *Server) LPush(ctx *redis.DBContext, key string, elems []string, opt redis.PushOption) (*redis.Message, error) {
	record, list, err := server.getDatabaseListRecord(ctx, key)
	if err != nil {
		return nil, err
	}

	for _, elem := range elems {
		list = append([]string{elem}, list...)
	}

	record.Data = list

	return redis.NewIntegerMessage(len(list)), nil
}

func (server *Server) RPush(ctx *redis.DBContext, key string, elems []string, opt redis.PushOption) (*redis.Message, error) {
	return nil, nil
}

func (server *Server) LPop(ctx *redis.DBContext, key string, count int) (*redis.Message, error) {
	return nil, nil
}

func (server *Server) RPop(ctx *redis.DBContext, key string, count int) (*redis.Message, error) {
	return nil, nil
}

func (server *Server) LRange(ctx *redis.DBContext, key string, start int, stop int) (*redis.Message, error) {
	_, list, err := server.getDatabaseListRecord(ctx, key)
	if err != nil {
		return nil, err
	}

	if start < 0 {
		start = len(list) + start
	}
	if stop < 0 {
		stop = len(list) + stop
	}

	arrayMsg := redis.NewArrayMessage()
	array, _ := arrayMsg.Array()
	for n := start; n <= stop; n++ {
		if (n < 0) || ((len(list) - 1) < n) {
			continue
		}
		array.Append(redis.NewBulkMessage(list[n]))
	}

	return arrayMsg, nil
}

func (server *Server) LIndex(ctx *redis.DBContext, key string, index int) (*redis.Message, error) {
	return nil, nil
}

func (server *Server) LLen(ctx *redis.DBContext, key string, index int) (*redis.Message, error) {
	return nil, nil
}
