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

package redis

// GenericCommandHandler represents a hander interface for genelic commands.
type GenericCommandHandler interface {
	// 1.0.0
	Del(ctx *DBContext, keys []string) (*Message, error)
	Exists(ctx *DBContext, keys []string) (*Message, error)
	Expire(ctx *DBContext, key string, opt ExpireOption) (*Message, error)
	Keys(ctx *DBContext, pattern string) (*Message, error)
	Rename(ctx *DBContext, key string, newkey string, opt RenameOption) (*Message, error)
	Type(ctx *DBContext, key string) (*Message, error)
	TTL(ctx *DBContext, key string) (*Message, error)
}

// StringHandler represents a core command hander interface for string commands.
type StringCommandHandler interface {
	Set(ctx *DBContext, key string, val string, opt SetOption) (*Message, error)
	Get(ctx *DBContext, key string) (*Message, error)
	MSet(ctx *DBContext, dict map[string]string, opt MSetOption) (*Message, error)
	MGet(ctx *DBContext, keys []string) (*Message, error)
}

// HashCommandHandler represents a core command hander interface for hash commands.
type HashCommandHandler interface {
	HDel(ctx *DBContext, key string, fields []string) (*Message, error)
	HSet(ctx *DBContext, key string, field string, val string, opt HSetOption) (*Message, error)
	HGet(ctx *DBContext, key string, field string) (*Message, error)
	HGetAll(ctx *DBContext, key string) (*Message, error)
	HMSet(ctx *DBContext, key string, dict map[string]string) (*Message, error)
	HMGet(ctx *DBContext, key string, fields []string) (*Message, error)
}

// ListCommandHandler represents a core command hander interface for list commands.
type ListCommandHandler interface {
	LPush(ctx *DBContext, key string, elements []string, opt PushOption) (*Message, error)
	RPush(ctx *DBContext, key string, elements []string, opt PushOption) (*Message, error)
	LPop(ctx *DBContext, key string, count int) (*Message, error)
	RPop(ctx *DBContext, key string, count int) (*Message, error)
	LRange(ctx *DBContext, key string, start int, stop int) (*Message, error)
	LIndex(ctx *DBContext, key string, index int) (*Message, error)
	LLen(ctx *DBContext, key string) (*Message, error)
}

// SetCommandHandler represents a core command hander interface for set commands.
type SetCommandHandler interface {
	SAdd(ctx *DBContext, key string, members []string) (*Message, error)
	SMembers(ctx *DBContext, key string) (*Message, error)
	SRem(ctx *DBContext, key string, members []string) (*Message, error)
}

type ZSetMember struct {
	Score float64
	Data  string
}

// ZSetCommandHandler represents a core command hander interface for zset commands.
type ZSetCommandHandler interface {
	ZAdd(ctx *DBContext, key string, members []*ZSetMember, opt ZAddOption) (*Message, error)
	ZRange(ctx *DBContext, key string, start int, stop int, opt ZRangeOption) (*Message, error)
	ZRangeByScore(ctx *DBContext, key string, min float64, max float64, opt ZRangeOption) (*Message, error)
	ZRem(ctx *DBContext, key string, members []string) (*Message, error)
	ZScore(ctx *DBContext, key string, member string) (*Message, error)
	ZIncBy(ctx *DBContext, key string, inc float64, member string) (*Message, error)
}

// CommandHandler represents a command hander interface for user commands.
type CommandHandler interface {
	GenericCommandHandler
	StringCommandHandler
	HashCommandHandler
	ListCommandHandler
	SetCommandHandler
	ZSetCommandHandler
}

// SystemCommandHandler represents a hander interface for system commands.
type SystemCommandHandler interface {
	Ping(ctx *DBContext, arg string) (*Message, error)
	Echo(ctx *DBContext, arg string) (*Message, error)
	Select(ctx *DBContext, index int) (*Message, error)
	Quit(ctx *DBContext) (*Message, error)
}
