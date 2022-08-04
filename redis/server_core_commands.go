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

// nolint: gocyclo, maintidx
func (server *Server) registerCoreExecutors() {
	// Registers connection management commands.

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

	// Registers string commands.

	server.RegisterExexutor("SET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		opt := newDefaultSetOption()
		key, val, err := nextSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Set(ctx, key, val, opt)
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

	server.RegisterExexutor("GET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		opt := GetOption{}
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Get(ctx, key, opt)
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

	// Registers hash commands.

	server.RegisterExexutor("HSET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		opt := HSetOption{}
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

	server.RegisterExexutor("HGET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		opt := HGetOption{}
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HGet(ctx, hash, key, opt)
	})

	server.RegisterExexutor("HGETALL", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HGetAll(ctx, hash)
	})

	server.RegisterExexutor("HMSET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		opt := HMSetOption{}
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		dir, err := nextMSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HMSet(ctx, hash, dir, opt)
	})

	server.RegisterExexutor("HMGET", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		opt := HMGetOption{}
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		keys, err := nextMGetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HMGet(ctx, hash, keys, opt)
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
		opt := MGetOption{}
		keys, err := nextMGetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.MGet(ctx, keys, opt)
	})
}
