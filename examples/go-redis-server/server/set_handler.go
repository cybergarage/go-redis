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

type Set []string

func (server *Server) SAdd(ctx *redis.DBContext, key string, members []string) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}

	record, sets, err := db.GetSetRecord(key)
	if err != nil {
		return nil, err
	}

	addedMemberCount := 0
	for _, member := range members {
		hasMember := false
		for _, set := range sets {
			if set == member {
				hasMember = true
				continue
			}
		}
		if hasMember {
			continue
		}
		sets = append(sets, member)
		addedMemberCount++
	}

	record.Data = sets

	return redis.NewIntegerMessage(addedMemberCount), nil
}

func (server *Server) SMembers(ctx *redis.DBContext, key string) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}

	_, sets, err := db.GetSetRecord(key)
	if err != nil {
		return nil, err
	}

	arrayMsg := redis.NewArrayMessage()
	array, _ := arrayMsg.Array()
	for _, set := range sets {
		array.Append(redis.NewBulkMessage(set))
	}

	return arrayMsg, nil
}

func (server *Server) SRem(ctx *redis.DBContext, key string, members []string) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}

	record, sets, err := db.GetSetRecord(key)
	if err != nil {
		return nil, err
	}

	removedMemberCount := 0
	for _, member := range members {
		for n, set := range sets {
			if set == member {
				sets = append(sets[:n], sets[n+1:]...)
				removedMemberCount++
				break
			}
		}
	}

	record.Data = sets

	return redis.NewIntegerMessage(removedMemberCount), nil
}
