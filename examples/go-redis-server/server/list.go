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
	"fmt"

	"github.com/cybergarage/go-redis/redis"
)

////////////////////////////////////////////////////////////
// List
////////////////////////////////////////////////////////////

type List struct {
	elements []string
}

func NewList() *List {
	return &List{
		elements: []string{},
	}
}

func (list *List) Pop(count int) ([]string, bool) {
	if count < 1 {
		return nil, false
	}

	elems := []string{}
	for n := 0; n < count; n++ {
		if len(list.elements) < 1 {
			continue
		}
		elems = append(elems, list.elements[0])
		list.elements = list.elements[1:]
	}

	return elems, true
}

////////////////////////////////////////////////////////////
// List command handler
////////////////////////////////////////////////////////////

func (server *Server) LPop(ctx *redis.DBContext, key string, count int) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}

	if !db.HasRecord(key) {
		return redis.NewNilMessage(), nil
	}

	record, list, err := db.GetListRecord(key)
	if err != nil {
		return nil, err
	}

	elems, ok := list.Pop(count)
	if !ok || len(elems) == 0 {
		return redis.NewNilMessage(), nil
	}

	if count == 1 {
		if len(elems) < 1 {
			return redis.NewNilMessage(), nil
		}
		return redis.NewBulkMessage(elems[0]), nil
	}

	arrayMsg := redis.NewArrayMessage()
	array, _ := arrayMsg.Array()
	for _, elem := range elems {
		array.Append(redis.NewBulkMessage(elem))
	}

	return arrayMsg, nil
}

func (server *Server) LPush(ctx *redis.DBContext, key string, elems []string, opt redis.PushOption) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}

	if opt.X {
		if !db.HasRecord(key) {
			return redis.NewIntegerMessage(0), nil
		}
	}

	record, list, err := db.GetListRecord(key)
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
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}

	if opt.X {
		if !db.HasRecord(key) {
			return redis.NewIntegerMessage(0), nil
		}
	}

	record, list, err := db.GetListRecord(key)
	if err != nil {
		return nil, err
	}

	for _, elem := range elems {
		list = append(list, elem)
	}

	record.Data = list

	return redis.NewIntegerMessage(len(list)), nil
}

func (server *Server) RPop(ctx *redis.DBContext, key string, count int) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}

	if !db.HasRecord(key) {
		return redis.NewNilMessage(), nil
	}

	record, list, err := db.GetListRecord(key)
	if err != nil {
		return nil, err
	}

	if count < 1 {
		return redis.NewNilMessage(), nil
	}

	if count == 1 {
		if len(list) < 1 {
			return redis.NewNilMessage(), nil
		}
		msg := redis.NewBulkMessage(list[len(list)-1])
		record.Data = list[:len(list)-1]
		return msg, nil
	}

	arrayMsg := redis.NewArrayMessage()
	array, _ := arrayMsg.Array()
	for n := 0; n < count; n++ {
		if len(list) < 1 {
			continue
		}
		array.Append(redis.NewBulkMessage(list[len(list)-1]))
		list = list[:len(list)-1]
	}
	record.Data = list

	return arrayMsg, nil
}

func (server *Server) LRange(ctx *redis.DBContext, key string, start int, stop int) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}

	_, list, err := db.GetListRecord(key)
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

func (server *Server) LIndex(ctx *redis.DBContext, key string, idx int) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}

	_, list, err := db.GetListRecord(key)
	if err != nil {
		return nil, err
	}

	if idx < 0 {
		idx = len(list) + idx
	}

	if (idx < 0) || ((len(list) - 1) < idx) {
		return redis.NewNilMessage(), nil
	}

	return redis.NewBulkMessage(list[idx]), nil
}

func (server *Server) LLen(ctx *redis.DBContext, key string) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}

	record, hasRecord := db.GetRecord(key)
	if !hasRecord {
		return redis.NewIntegerMessage(0), nil
	}

	list, ok := record.Data.(List)
	if !ok {
		return redis.NewErrorMessage(fmt.Errorf(errorInvalidStoredDataType, record.Data)), nil
	}

	return redis.NewIntegerMessage(len(list)), nil
}
