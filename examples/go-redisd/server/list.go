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

func (list *List) LPop(count int) ([]string, bool) {
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

func (list *List) LPush(elems []string) int {
	for _, elem := range elems {
		list.elements = append([]string{elem}, list.elements...)
	}
	return len(list.elements)
}

func (list *List) RPop(count int) ([]string, bool) {
	if count < 1 {
		return nil, false
	}
	elems := []string{}
	for n := 0; n < count; n++ {
		if len(list.elements) < 1 {
			continue
		}
		elems = append(elems, list.elements[len(list.elements)-1])
		list.elements = list.elements[:len(list.elements)-1]
	}
	return elems, true
}

func (list *List) RPush(elems []string) int {
	list.elements = append(list.elements, elems...)
	return len(list.elements)
}

func (list *List) Range(start int, stop int) []string {
	if start < 0 {
		start = len(list.elements) + start
	}
	if stop < 0 {
		stop = len(list.elements) + stop
	}
	elems := []string{}
	for n := start; n <= stop; n++ {
		if (n < 0) || ((len(list.elements) - 1) < n) {
			continue
		}
		elems = append(elems, list.elements[n])
	}
	return elems
}

func (list *List) Index(idx int) (string, bool) {
	if idx < 0 {
		idx = len(list.elements) + idx
	}
	if (idx < 0) || ((len(list.elements) - 1) < idx) {
		return "", false
	}
	return list.elements[idx], true
}

func (list *List) Len() int {
	return len(list.elements)
}

////////////////////////////////////////////////////////////
// List command handler
////////////////////////////////////////////////////////////

func (server *Server) pop(conn *redis.Conn, key string, count int, isLPop bool) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}

	if !db.HasRecord(key) {
		return redis.NewNilMessage(), nil
	}

	_, list, err := db.GetListRecord(key)
	if err != nil {
		return nil, err
	}

	var elems []string
	var ok bool
	if isLPop {
		elems, ok = list.LPop(count)
	} else {
		elems, ok = list.RPop(count)
	}

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

func (server *Server) push(conn *redis.Conn, key string, elems []string, opt redis.PushOption, isLPop bool) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}

	if opt.X {
		if !db.HasRecord(key) {
			return redis.NewIntegerMessage(0), nil
		}
	}

	_, list, err := db.GetListRecord(key)
	if err != nil {
		return nil, err
	}

	var cnt int
	if isLPop {
		cnt = list.LPush(elems)
	} else {
		cnt = list.RPush(elems)
	}

	return redis.NewIntegerMessage(cnt), nil
}

func (server *Server) LPop(conn *redis.Conn, key string, count int) (*redis.Message, error) {
	return server.pop(conn, key, count, true)
}

func (server *Server) LPush(conn *redis.Conn, key string, elems []string, opt redis.PushOption) (*redis.Message, error) {
	return server.push(conn, key, elems, opt, true)
}

func (server *Server) RPop(conn *redis.Conn, key string, count int) (*redis.Message, error) {
	return server.pop(conn, key, count, false)
}

func (server *Server) RPush(conn *redis.Conn, key string, elems []string, opt redis.PushOption) (*redis.Message, error) {
	return server.push(conn, key, elems, opt, false)
}

func (server *Server) LRange(conn *redis.Conn, key string, start int, stop int) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}

	_, list, err := db.GetListRecord(key)
	if err != nil {
		return nil, err
	}

	elems := list.Range(start, stop)
	arrayMsg := redis.NewArrayMessage()
	array, _ := arrayMsg.Array()
	for _, elem := range elems {
		array.Append(redis.NewBulkMessage(elem))
	}

	return arrayMsg, nil
}

func (server *Server) LIndex(conn *redis.Conn, key string, idx int) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}

	_, list, err := db.GetListRecord(key)
	if err != nil {
		return nil, err
	}

	elem, ok := list.Index(idx)
	if !ok {
		return redis.NewNilMessage(), nil
	}

	return redis.NewBulkMessage(elem), nil
}

func (server *Server) LLen(conn *redis.Conn, key string) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}

	_, list, err := db.GetListRecord(key)
	if err != nil {
		return nil, err
	}

	return redis.NewIntegerMessage(list.Len()), nil
}
