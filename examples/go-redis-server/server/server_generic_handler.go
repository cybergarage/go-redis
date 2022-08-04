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

func (server *Server) Exists(ctx *redis.DBContext, keys []string) (*redis.Message, error) {
	db := server.GetDatabase(ctx.ID())
	existCount := 0
	for _, key := range keys {
		_, ok := db.GetRecord(key)
		if ok {
			existCount++
		}
	}
	return redis.NewIntegerMessage(existCount), nil
}

func (server *Server) Type(ctx *redis.DBContext, key string) (*redis.Message, error) {
	db := server.GetDatabase(ctx.ID())
	record, ok := db.GetRecord(key)
	if !ok {
		return redis.NewStringMessage("none"), nil
	}
	switch record.Data.(type) {
	case string:
		return redis.NewStringMessage("string"), nil
	case HashData:
		return redis.NewStringMessage("hash"), nil
	}
	return redis.NewStringMessage("none"), nil
}
