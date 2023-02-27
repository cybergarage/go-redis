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
	"github.com/cybergarage/go-redis/redis"
)

////////////////////////////////////////////////////////////
// Set
////////////////////////////////////////////////////////////

type Set struct {
	members []string
}

func NewSet() *Set {
	return &Set{
		members: []string{},
	}
}

func (set *Set) Add(members []string) int {
	addedMemberCount := 0
	for _, member := range members {
		hasMember := false
		for _, set := range set.members {
			if set == member {
				hasMember = true
				continue
			}
		}
		if hasMember {
			continue
		}
		set.members = append(set.members, member)
		addedMemberCount++
	}
	return addedMemberCount
}

func (set *Set) Rem(members []string) int {
	removedMemberCount := 0
	for _, rm := range members {
		for n, m := range set.members {
			if m == rm {
				set.members = append(set.members[:n], set.members[n+1:]...)
				removedMemberCount++
				break
			}
		}
	}
	return removedMemberCount
}

func (set *Set) Members() []string {
	return set.members
}

////////////////////////////////////////////////////////////
// Set command handler
////////////////////////////////////////////////////////////

func (server *Server) SAdd(conn *redis.Conn, key string, members []string) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}
	_, set, err := db.GetSetRecord(key)
	if err != nil {
		return nil, err
	}
	return redis.NewIntegerMessage(set.Add(members)), nil
}

func (server *Server) SMembers(conn *redis.Conn, key string) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}
	_, set, err := db.GetSetRecord(key)
	if err != nil {
		return nil, err
	}
	arrayMsg := redis.NewArrayMessage()
	array, _ := arrayMsg.Array()
	for _, m := range set.Members() {
		array.Append(redis.NewBulkMessage(m))
	}

	return arrayMsg, nil
}

func (server *Server) SRem(conn *redis.Conn, key string, members []string) (*redis.Message, error) {
	db, err := server.GetDatabase(conn.Database())
	if err != nil {
		return nil, err
	}
	_, set, err := db.GetSetRecord(key)
	if err != nil {
		return nil, err
	}
	return redis.NewIntegerMessage(set.Rem(members)), nil
}
