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

import "strconv"

// nolint: gocyclo, maintidx
func (server *Server) registerSugarExecutors() {
	// common internal sugar functions

	incdecExecutor := func(ctx *DBContext, cmd string, key string, val int) (*Message, error) {
		getOpt := GetOption{}
		getRet, err := server.userCommandHandler.Get(ctx, key, getOpt)
		if err != nil {
			return nil, err
		}
		currVal := 0
		if !getRet.IsNil() {
			retVal, err := getRet.Integer()
			if err != nil {
				return nil, err
			}
			currVal = retVal
		}
		newVal := currVal + val
		opt := newDefaultSetOption()
		_, err = server.userCommandHandler.Set(ctx, key, strconv.Itoa(newVal), opt)
		if err != nil {
			return nil, err
		}
		return NewIntegerMessage(newVal), nil
	}

	getRangeExecutor := func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		rageValidiator := func(val int, max int) int {
			if val < 0 {
				val = max + val
				if val < 0 {
					return 0
				}
			}
			if max < val {
				val = max - 1
			}
			return val
		}
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
		getOpt := GetOption{}
		getRet, err := server.userCommandHandler.Get(ctx, key, getOpt)
		if err != nil {
			return NewNilMessage(), nil
		}
		getVal, err := getRet.String()
		if err != nil {
			return NewNilMessage(), nil
		}
		start = rageValidiator(start, len(getVal))
		end = rageValidiator(end, len(getVal))
		return NewBulkMessage(getVal[start:(end + 1)]), nil
	}

	// Registers sugar string commands.

	server.RegisterExexutor("APPEND", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, appendVal, err := nextSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		getOpt := GetOption{}
		getRet, err := server.userCommandHandler.Get(ctx, key, getOpt)
		if err != nil {
			return nil, err
		}
		newVal := appendVal
		if getVal, err := getRet.String(); err == nil {
			newVal = getVal + appendVal
		}
		opt := newDefaultSetOption()
		_, err = server.userCommandHandler.Set(ctx, key, newVal, opt)
		if err != nil {
			return nil, err
		}
		return NewIntegerMessage(len(newVal)), nil
	})

	server.RegisterExexutor("DECR", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return incdecExecutor(ctx, cmd, key, -1)
	})

	server.RegisterExexutor("DECRBY", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		inc, err := nextIntegerArgument(cmd, "decrement", args)
		if err != nil {
			return nil, err
		}
		return incdecExecutor(ctx, cmd, key, -inc)
	})

	server.RegisterExexutor("GETRANGE", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		return getRangeExecutor(ctx, cmd, args)
	})

	server.RegisterExexutor("INCR", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return incdecExecutor(ctx, cmd, key, 1)
	})

	server.RegisterExexutor("INCRBY", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		inc, err := nextIntegerArgument(cmd, "increment", args)
		if err != nil {
			return nil, err
		}
		return incdecExecutor(ctx, cmd, key, inc)
	})

	server.RegisterExexutor("STRLEN", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		getOpt := GetOption{}
		getRet, err := server.userCommandHandler.Get(ctx, key, getOpt)
		if err != nil {
			return NewIntegerMessage(0), nil
		}
		getVal, err := getRet.String()
		if err != nil {
			return NewIntegerMessage(0), nil
		}
		return NewIntegerMessage(len(getVal)), nil
	})

	server.RegisterExexutor("SUBSTR", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		return getRangeExecutor(ctx, cmd, args)
	})

	// Registers sugar hash commands.

	server.RegisterExexutor("HEXISTS", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		opt := HGetOption{}
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		field, err := nextStringArgument(cmd, "field", args)
		if err != nil {
			return nil, err
		}
		getRet, err := server.userCommandHandler.HGet(ctx, key, field, opt)
		if err != nil {
			return NewIntegerMessage(0), nil
		}
		_, err = getRet.String()
		if err != nil {
			return NewIntegerMessage(0), nil
		}
		return NewIntegerMessage(1), nil
	})
}
