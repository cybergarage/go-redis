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

import "time"

// nolint: gocyclo, maintidx
func (server *Server) registerCoreExecutors() {
	// Connection management commands.

	server.RegisterExexutor("PING", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		arg := ""
		var err error
		if msg, _ := args.Next(); msg != nil {
			arg, err = msg.String()
			if err != nil {
				return nil, err
			}
		}
		return server.systemCommandHandler.Ping(ctx, arg)
	})

	server.RegisterExexutor("ECHO", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		msg, err := args.NextString()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "msg", err)
		}
		return server.systemCommandHandler.Echo(ctx, msg)
	})

	server.RegisterExexutor("SELECT", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		id, err := args.NextInteger()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "id", err)
		}
		return server.systemCommandHandler.Select(ctx, id)
	})

	server.RegisterExexutor("QUIT", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		return server.systemCommandHandler.Quit(ctx)
	})

	// Generic commands.

	server.RegisterExexutor("DEL", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		keys, err := nextKeysArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Del(ctx, keys)
	})

	server.RegisterExexutor("EXPIRE", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		ttl, err := nextIntegerArgument(cmd, "ttl", args)
		if err != nil {
			return nil, err
		}
		ttlTime := time.Now().Add(time.Duration(ttl) * time.Second)
		opt, err := nextExpireArgument(cmd, ttlTime, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Expire(ctx, key, opt)
	})

	server.RegisterExexutor("EXPIREAT", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		ttl, err := nextIntegerArgument(cmd, "ttl", args)
		if err != nil {
			return nil, err
		}
		ttlTime := time.Unix(int64(ttl), 0)
		opt, err := nextExpireArgument(cmd, ttlTime, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Expire(ctx, key, opt)
	})

	server.RegisterExexutor("EXISTS", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		keys, err := nextKeysArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Exists(ctx, keys)
	})

	server.RegisterExexutor("KEYS", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		pattern, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Keys(ctx, pattern)
	})

	server.RegisterExexutor("TYPE", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Type(ctx, key)
	})

	server.RegisterExexutor("RENAME", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		newkey, err := nextStringArgument(cmd, "newkey", args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Rename(ctx, key, newkey, RenameOption{NX: false})
	})

	server.RegisterExexutor("RENAMENX", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		newkey, err := nextStringArgument(cmd, "newkey", args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Rename(ctx, key, newkey, RenameOption{NX: true})
	})

	server.RegisterExexutor("TTL", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.TTL(ctx, key)
	})

	// String commands.

	server.RegisterExexutor("GET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Get(ctx, key)
	})

	server.RegisterExexutor("SET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		opt := newDefaultSetOption()
		key, val, err := nextSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Set(ctx, key, val, opt)
	})

	server.RegisterExexutor("GETSET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		opt := newDefaultSetOption()
		opt.GET = true
		key, val, err := nextSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Set(ctx, key, val, opt)
	})

	server.RegisterExexutor("MSET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		opt := MSetOption{
			NX: false,
		}
		dir, err := nextMSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.MSet(ctx, dir, opt)
	})

	server.RegisterExexutor("MSETNX", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		opt := MSetOption{
			NX: true,
		}
		dir, err := nextMSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.MSet(ctx, dir, opt)
	})

	server.RegisterExexutor("MGET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		keys, err := nextKeysArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.MGet(ctx, keys)
	})

	server.RegisterExexutor("SETNX", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		opt := newDefaultSetOption()
		opt.NX = true
		key, val, err := nextSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Set(ctx, key, val, opt)
	})

	// Hash commands.

	server.RegisterExexutor("HDEL", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		keys, err := nextKeysArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HDel(ctx, hash, keys)
	})

	server.RegisterExexutor("HGET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HGet(ctx, hash, key)
	})

	server.RegisterExexutor("HGETALL", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HGetAll(ctx, hash)
	})

	server.RegisterExexutor("HSET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		opt := HSetOption{
			NX: false,
		}
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		key, val, err := nextSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HSet(ctx, hash, key, val, opt)
	})

	server.RegisterExexutor("HSETNX", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		opt := HSetOption{
			NX: true,
		}
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		key, val, err := nextSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HSet(ctx, hash, key, val, opt)
	})

	server.RegisterExexutor("HMSET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		dir, err := nextMSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HMSet(ctx, hash, dir)
	})

	server.RegisterExexutor("HMGET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		keys, err := nextKeysArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HMGet(ctx, hash, keys)
	})

	// List commands.

	server.RegisterExexutor("LINDEX", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		idx, err := nextIntegerArgument(cmd, "index", args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.LIndex(ctx, key, idx)
	})

	server.RegisterExexutor("LLEN", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.LLen(ctx, key)
	})

	server.RegisterExexutor("LPOP", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, cnt, err := nextPopArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.LPop(ctx, key, cnt)
	})

	server.RegisterExexutor("LPUSH", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, elems, err := nextPushArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		opt := PushOption{X: false}
		return server.userCommandHandler.LPush(ctx, key, elems, opt)
	})

	server.RegisterExexutor("LPUSHX", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, elems, err := nextPushArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		opt := PushOption{X: true}
		return server.userCommandHandler.LPush(ctx, key, elems, opt)
	})

	server.RegisterExexutor("LRANGE", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		start, err := nextIntegerArgument(cmd, "start", args)
		if err != nil {
			return nil, err
		}
		end, err := nextIntegerArgument(cmd, "end", args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.LRange(ctx, key, start, end)
	})

	server.RegisterExexutor("RPOP", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, cnt, err := nextPopArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.RPop(ctx, key, cnt)
	})

	server.RegisterExexutor("RPUSH", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, elems, err := nextPushArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		opt := PushOption{X: false}
		return server.userCommandHandler.RPush(ctx, key, elems, opt)
	})

	server.RegisterExexutor("RPUSHX", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, elems, err := nextPushArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		opt := PushOption{X: true}
		return server.userCommandHandler.RPush(ctx, key, elems, opt)
	})

	// Set commands.

	server.RegisterExexutor("SADD", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		members, err := nextStringsArguments(cmd, "member", args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.SAdd(ctx, key, members)
	})

	server.RegisterExexutor("SMEMBERS", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.SMembers(ctx, key)
	})

	server.RegisterExexutor("SREM", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		members, err := nextStringsArguments(cmd, "member", args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.SRem(ctx, key, members)
	})
}
