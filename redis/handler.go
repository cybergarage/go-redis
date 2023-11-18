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

package redis

// ConnectionManagementCommandHandler represents a hander interface for connection management commands.
type ConnectionManagementCommandHandler interface {
	Ping(conn *Conn, arg string) (*Message, error)
	Echo(conn *Conn, arg string) (*Message, error)
	Select(conn *Conn, index int) (*Message, error)
	Quit(conn *Conn) (*Message, error)
}

// ServerManagementCommandHandler represents a hander interface for server management commands.
type ServerManagementCommandHandler interface {
	ConfigSet(conn *Conn, params map[string]string) (*Message, error)
	ConfigGet(conn *Conn, keys []string) (*Message, error)
}

// GenericCommandHandler represents a hander interface for genelic commands.
type GenericCommandHandler interface {
	Del(conn *Conn, keys []string) (*Message, error)
	Exists(conn *Conn, keys []string) (*Message, error)
	Expire(conn *Conn, key string, opt ExpireOption) (*Message, error)
	Keys(conn *Conn, pattern string) (*Message, error)
	Rename(conn *Conn, key string, newkey string, opt RenameOption) (*Message, error)
	Type(conn *Conn, key string) (*Message, error)
	TTL(conn *Conn, key string) (*Message, error)
}

// StringCommandHandler represents a core command hander interface for string commands.
// APPEND, DECR, DECRBY, GETRANGE, GETSET, INCR, INCRBY, MGET, MSET, MSETNX, SETRANGE, STRLEN commands are implemented by the StringCommandHandler.
type StringCommandHandler interface {
	// Set represents a handler interface for SET, SETNX, SETEX, PSETEX, MSET and MSETNX commands.
	Set(conn *Conn, key string, val string, opt SetOption) (*Message, error)
	// Get represents a handler interface for GET and MGET commands.
	Get(conn *Conn, key string) (*Message, error)
}

// HashCommandHandler represents a core command hander interface for hash commands.
type HashCommandHandler interface {
	HDel(conn *Conn, key string, fields []string) (*Message, error)
	HSet(conn *Conn, key string, field string, val string, opt HSetOption) (*Message, error)
	HGet(conn *Conn, key string, field string) (*Message, error)
	HGetAll(conn *Conn, key string) (*Message, error)
	HMSet(conn *Conn, key string, dict map[string]string) (*Message, error)
	HMGet(conn *Conn, key string, fields []string) (*Message, error)
}

// ListCommandHandler represents a core command hander interface for list commands.
type ListCommandHandler interface {
	LPush(conn *Conn, key string, elements []string, opt PushOption) (*Message, error)
	RPush(conn *Conn, key string, elements []string, opt PushOption) (*Message, error)
	LPop(conn *Conn, key string, count int) (*Message, error)
	RPop(conn *Conn, key string, count int) (*Message, error)
	LRange(conn *Conn, key string, start int, stop int) (*Message, error)
	LIndex(conn *Conn, key string, index int) (*Message, error)
	LLen(conn *Conn, key string) (*Message, error)
}

// SetCommandHandler represents a core command hander interface for set commands.
type SetCommandHandler interface {
	SAdd(conn *Conn, key string, members []string) (*Message, error)
	SMembers(conn *Conn, key string) (*Message, error)
	SRem(conn *Conn, key string, members []string) (*Message, error)
}

// ZSetMember represents a parameter for ZSetCommandHandler.
type ZSetMember struct {
	Score  float64
	Member string
}

// ZSetCommandHandler represents a core command hander interface for zset commands.
type ZSetCommandHandler interface {
	ZAdd(conn *Conn, key string, members []*ZSetMember, opt ZAddOption) (*Message, error)
	ZRange(conn *Conn, key string, start int, stop int, opt ZRangeOption) (*Message, error)
	ZRangeByScore(conn *Conn, key string, min float64, max float64, opt ZRangeOption) (*Message, error)
	ZRem(conn *Conn, key string, members []string) (*Message, error)
	ZScore(conn *Conn, key string, member string) (*Message, error)
	ZIncBy(conn *Conn, key string, inc float64, member string) (*Message, error)
}

// AuthCommandHandler represents a hander interface for authentication commands.
type AuthCommandHandler interface {
	Auth(conn *Conn, username string, password string) (*Message, error)
}

// UserCommandHandler represents a command hander interface for user commands.
type UserCommandHandler interface {
	GenericCommandHandler
	StringCommandHandler
	HashCommandHandler
	ListCommandHandler
	SetCommandHandler
	ZSetCommandHandler
}

// SystemCommandHandler represents a hander interface for system commands.
type SystemCommandHandler interface {
	ConnectionManagementCommandHandler
	ServerManagementCommandHandler
}
