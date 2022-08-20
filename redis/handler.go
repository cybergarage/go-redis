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
	Get(ctx *DBContext, key string, opt GetOption) (*Message, error)
	MSet(ctx *DBContext, dict map[string]string, opt MSetOption) (*Message, error)
	MGet(ctx *DBContext, keys []string, opt MGetOption) (*Message, error)
}

// HashCommandHandler represents a core command hander interface for hash commands.
type HashCommandHandler interface {
	HSet(ctx *DBContext, key string, field string, val string, opt HSetOption) (*Message, error)
	HGet(ctx *DBContext, key string, field string, opt HGetOption) (*Message, error)
	HGetAll(ctx *DBContext, key string) (*Message, error)
	HMSet(ctx *DBContext, key string, dict map[string]string, opt HMSetOption) (*Message, error)
	HMGet(ctx *DBContext, key string, fields []string, opt HMGetOption) (*Message, error)
}

// CommandHandler represents a command hander interface for user commands.
type CommandHandler interface {
	GenericCommandHandler
	StringCommandHandler
	HashCommandHandler
}

// SystemCommandHandler represents a hander interface for system commands.
type SystemCommandHandler interface {
	Ping(ctx *DBContext, arg string) (*Message, error)
	Echo(ctx *DBContext, arg string) (*Message, error)
	Select(ctx *DBContext, index int) (*Message, error)
	Quit(ctx *DBContext) (*Message, error)
}
