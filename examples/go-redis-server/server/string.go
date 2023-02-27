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

func (server *Server) Set(conn *redis.Conn, key string, val string, opt redis.SetOption) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}

	var oldVal []byte
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

func (server *Server) Get(conn *redis.Conn, key string) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}
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

func (server *Server) MSet(conn *redis.Conn, dict map[string]string, opt redis.MSetOption) (*redis.Message, error) {
	if opt.NX {
		for key := range dict {
			res, err := server.Get(conn, key)
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
		if _, err := server.Set(conn, key, val, setOpt); err != nil {
			return nil, err
		}
	}

	if opt.NX {
		return redis.NewIntegerMessage(1), nil
	}
	return redis.NewOKMessage(), nil
}

func (server *Server) MGet(conn *redis.Conn, keys []string) (*redis.Message, error) {
	arrayMsg := redis.NewArrayMessage()
	array, _ := arrayMsg.Array()
	for _, key := range keys {
		msg, err := server.Get(conn, key)
		if err != nil {
			return nil, err
		}
		array.Append(msg)
	}
	return arrayMsg, nil
}
