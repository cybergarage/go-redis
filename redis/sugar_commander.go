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

import (
	"strconv"
)

// nolint: gocyclo, maintidx, nilerr
func (server *Server) registerSugarExecutors() {
	// common internal sugar functions

	incdecExecutor := func(conn *Conn, cmd string, key string, val int) (*Message, error) {
		getRet, err := server.userCommandHandler.Get(conn, key)
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
		_, err = server.userCommandHandler.Set(conn, key, strconv.Itoa(newVal), opt)
		if err != nil {
			return nil, err
		}
		return NewIntegerMessage(newVal), nil
	}

	// Registers sugar string commands.

	server.RegisterExexutor("APPEND", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, appendVal, err := nextSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		getRet, err := server.userCommandHandler.Get(conn, key)
		if err != nil {
			return nil, err
		}
		newVal := appendVal
		if getVal, err := getRet.String(); err == nil {
			newVal = getVal + appendVal
		}
		opt := newDefaultSetOption()
		_, err = server.userCommandHandler.Set(conn, key, newVal, opt)
		if err != nil {
			return nil, err
		}
		return NewIntegerMessage(len(newVal)), nil
	})

	server.RegisterExexutor("DECR", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return incdecExecutor(conn, cmd, key, -1)
	})

	server.RegisterExexutor("DECRBY", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		inc, err := nextIntegerArgument(cmd, "decrement", args)
		if err != nil {
			return nil, err
		}
		return incdecExecutor(conn, cmd, key, -inc)
	})

	server.RegisterExexutor("GETRANGE", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
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
		getRet, err := server.userCommandHandler.Get(conn, key)
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
	})

	server.RegisterExexutor("INCR", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return incdecExecutor(conn, cmd, key, 1)
	})

	server.RegisterExexutor("INCRBY", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		inc, err := nextIntegerArgument(cmd, "increment", args)
		if err != nil {
			return nil, err
		}
		return incdecExecutor(conn, cmd, key, inc)
	})

	server.RegisterExexutor("STRLEN", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		getRet, err := server.executeCommand(conn, "GET", args)
		if err != nil {
			return NewIntegerMessage(0), nil
		}
		getVal, err := getRet.String()
		if err != nil {
			return NewIntegerMessage(0), nil
		}
		return NewIntegerMessage(len(getVal)), nil
	})

	server.RegisterExexutor("SUBSTR", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		return server.executeCommand(conn, "GETRANGE", args)
	})

	// Registers sugar hash commands.

	server.RegisterExexutor("HEXISTS", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		getRet, err := server.executeCommand(conn, "HGET", args)
		if err != nil {
			return NewIntegerMessage(0), nil
		}
		_, err = getRet.String()
		if err != nil {
			return NewIntegerMessage(0), nil
		}
		return NewIntegerMessage(1), nil
	})

	server.RegisterExexutor("HKEYS", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		getAllRet, err := server.executeCommand(conn, "HGETALL", args)
		if err != nil {
			return nil, err
		}
		arrayMsg, err := getAllRet.Array()
		if err != nil {
			return NewArrayMessage(), nil
		}
		retMsg := NewArrayMessage()
		nextMsg, err := arrayMsg.Next()
		for nextMsg != nil {
			// Appends a next key
			key, err := nextMsg.String()
			if err != nil {
				break
			}
			retMsg.Append(NewBulkMessage(key))
			// Skips a next value string
			_, err = arrayMsg.Next()
			if err != nil {
				break
			}
			// Reads a next key string
			nextMsg, err = arrayMsg.Next()
			if err != nil {
				break
			}
		}
		if err != nil {
			return nil, err
		}
		return retMsg, nil
	})

	server.RegisterExexutor("HLEN", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		retMsg, err := server.executeCommand(conn, "HKEYS", args)
		if err != nil {
			return nil, err
		}
		arrayMsg, err := retMsg.Array()
		if err != nil {
			return nil, err
		}
		return NewIntegerMessage(arrayMsg.Size()), nil
	})

	server.RegisterExexutor("HSTRLEN", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		retMsg, err := server.executeCommand(conn, "HGET", args)
		if err != nil {
			return nil, err
		}
		if retMsg.IsNil() {
			return NewIntegerMessage(0), nil
		}
		retStr, err := retMsg.String()
		if err != nil {
			return nil, err
		}
		return NewIntegerMessage(len(retStr)), nil
	})

	server.RegisterExexutor("HVALS", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		getAllRet, err := server.executeCommand(conn, "HGETALL", args)
		if err != nil {
			return nil, err
		}
		arrayMsg, err := getAllRet.Array()
		if err != nil {
			return NewArrayMessage(), nil
		}
		retMsg := NewArrayMessage()
		nextMsg, err := arrayMsg.Next()
		for nextMsg != nil {
			// Skips a next key, and adds a next value string
			nextMsg, err = arrayMsg.Next()
			if nextMsg == nil || err != nil {
				break
			}
			val, err := nextMsg.String()
			if err != nil {
				break
			}
			retMsg.Append(NewBulkMessage(val))
			// Reads a next key string
			nextMsg, err = arrayMsg.Next()
			if err != nil {
				break
			}
		}
		if err != nil {
			return nil, err
		}
		return retMsg, nil
	})

	// Registers sugar set commands.

	server.RegisterExexutor("SCARD", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		retMsg, err := server.userCommandHandler.SMembers(conn, key)
		if err != nil {
			return nil, err
		}

		arrayMsg, err := retMsg.Array()
		if err != nil {
			return NewIntegerMessage(0), nil
		}

		memberCount := 0
		nextMsg, _ := arrayMsg.Next()
		for nextMsg != nil {
			memberCount++
			nextMsg, _ = arrayMsg.Next()
		}

		return NewIntegerMessage(memberCount), nil
	})

	server.RegisterExexutor("SISMEMBER", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		member, err := nextStringArgument(cmd, "member", args)
		if err != nil {
			return nil, err
		}

		retMsg, err := server.userCommandHandler.SMembers(conn, key)
		if err != nil {
			return nil, err
		}

		arrayMsg, err := retMsg.Array()
		if err != nil {
			return NewIntegerMessage(0), nil
		}

		nextMsg, _ := arrayMsg.Next()
		for nextMsg != nil {
			val, err := nextMsg.String()
			if err != nil {
				break
			}
			if val == member {
				return NewIntegerMessage(1), nil
			}
			nextMsg, _ = arrayMsg.Next()
		}

		return NewIntegerMessage(0), nil
	})

	// Registers sugar set commands.

	server.RegisterExexutor("ZCARD", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}

		opt := ZRangeOption{
			BYSCORE:      false,
			BYLEX:        false,
			REV:          false,
			WITHSCORES:   false,
			MINEXCLUSIVE: false,
			MAXEXCLUSIVE: false,
			Offset:       0,
			Count:        -1,
		}

		retMsg, err := server.userCommandHandler.ZRange(conn, key, 0, -1, opt)
		if err != nil {
			return nil, err
		}

		arrayMsg, err := retMsg.Array()
		if err != nil {
			return NewIntegerMessage(0), nil
		}

		memberCount := 0
		nextMsg, _ := arrayMsg.Next()
		for nextMsg != nil {
			memberCount++
			nextMsg, _ = arrayMsg.Next()
		}

		return NewIntegerMessage(memberCount), nil
	})
}
