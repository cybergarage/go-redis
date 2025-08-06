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

////////////////////////////////////////////////////////////
// Hash
////////////////////////////////////////////////////////////

type Hash map[string]string

// nolint: ifshort
func (hash Hash) Set(field string, val string, opt redis.HSetOption) int {
	_, hasKey := hash[field]
	if opt.NX && hasKey {
		return 0
	}

	hash[field] = val

	if hasKey {
		return 0
	}

	return 1
}

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

////////////////////////////////////////////////////////////
// Hash command handler
////////////////////////////////////////////////////////////

func (server *Server) HDel(conn *redis.Conn, key string, fields []string) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}

	record, ok := db.GetRecord(key)
	if !ok {
		return redis.NewIntegerMessage(0), nil
	}

	hash, ok := record.Data.(Hash)
	if !ok {
		return redis.NewIntegerMessage(0), nil
	}

	return redis.NewIntegerMessage(hash.Del(fields)), nil
}

// nolint: ifshort
func (server *Server) HSet(conn *redis.Conn, key string, field string, val string, opt redis.HSetOption) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}

	var hash Hash

	record, hasRecord := db.GetRecord(key)
	if hasRecord {
		var ok bool

		hash, ok = record.Data.(Hash)
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

	return redis.NewIntegerMessage(hash.Set(field, val, opt)), nil
}

func (server *Server) HGet(conn *redis.Conn, key string, field string) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}

	record, ok := db.GetRecord(key)
	if !ok {
		return redis.NewNilMessage(), nil
	}

	hash, ok := record.Data.(Hash)
	if !ok {
		return redis.NewNilMessage(), nil
	}

	hashData, ok := hash[field]
	if !ok {
		return redis.NewNilMessage(), nil
	}

	return redis.NewStringMessage(hashData), nil
}

func (server *Server) HGetAll(conn *redis.Conn, key string) (*redis.Message, error) {
	arrayMsg := redis.NewArrayMessage()

	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}

	record, ok := db.GetRecord(key)
	if !ok {
		return arrayMsg, nil
	}

	hash, ok := record.Data.(Hash)
	if !ok {
		return arrayMsg, nil
	}

	array, _ := arrayMsg.Array()
	for key, val := range hash {
		array.Append(redis.NewBulkMessage(key))
		array.Append(redis.NewBulkMessage(val))
	}

	return arrayMsg, nil
}
