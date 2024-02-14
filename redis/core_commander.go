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

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/cybergarage/go-redis/redis/proto"
)

// nolint: gocyclo, maintidx
func (server *Server) registerCoreExecutors() {
	// Connection management commands.

	server.RegisterExexutor("AUTH", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		var user, passwd string
		var err error
		passwd, err = args.NextString()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "password", err)
		}
		if msg, _ := args.Next(); msg != nil {
			token, err := msg.String()
			if err != nil {
				return nil, err
			}
			user = passwd
			passwd = token
		}

		if server.authCommandHandler == nil {
			return NewErrorNotSupportedMessage("AUTH"), nil
		}

		return server.authCommandHandler.Auth(conn, user, passwd)
	})

	server.RegisterExexutor("PING", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		arg := ""
		var err error
		if msg, _ := args.Next(); msg != nil {
			arg, err = msg.String()
			if err != nil {
				return nil, err
			}
		}
		return server.systemCommandHandler.Ping(conn, arg)
	})

	server.RegisterExexutor("ECHO", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		msg, err := args.NextString()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "msg", err)
		}
		return server.systemCommandHandler.Echo(conn, msg)
	})

	server.RegisterExexutor("SELECT", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		id, err := args.NextInteger()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "id", err)
		}
		msg, err := server.systemCommandHandler.Select(conn, id)
		if err == nil {
			conn.SetDatabase(id)
		}
		return msg, err
	})

	server.RegisterExexutor("QUIT", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		return server.systemCommandHandler.Quit(conn)
	})

	// Server management commands.

	server.RegisterExexutor("CONFIG", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		opt := ""
		var err error
		if msg, _ := args.Next(); msg != nil {
			opt, err = msg.String()
			if err != nil {
				return nil, err
			}
		}

		switch strings.ToUpper(opt) {
		case "SET":
			params, err := nextStringMapArguments(cmd, args)
			if err != nil {
				return nil, err
			}
			return server.systemCommandHandler.ConfigSet(conn, params)
		case "GET":
			params, err := nextStringArrayArguments(cmd, "params", args)
			if err != nil {
				return nil, err
			}
			return server.systemCommandHandler.ConfigGet(conn, params)
		}

		return nil, errors.New(opt)
	})

	// Generic commands.

	server.RegisterExexutor("DEL", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		keys, err := nextKeysArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Del(conn, keys)
	})

	server.RegisterExexutor("EXPIRE", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
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
		return server.userCommandHandler.Expire(conn, key, opt)
	})

	server.RegisterExexutor("EXPIREAT", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
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
		return server.userCommandHandler.Expire(conn, key, opt)
	})

	server.RegisterExexutor("EXISTS", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		keys, err := nextKeysArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Exists(conn, keys)
	})

	server.RegisterExexutor("KEYS", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		pattern, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Keys(conn, pattern)
	})

	server.RegisterExexutor("TYPE", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Type(conn, key)
	})

	server.RegisterExexutor("RENAME", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		newkey, err := nextStringArgument(cmd, "newkey", args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Rename(conn, key, newkey, RenameOption{NX: false})
	})

	server.RegisterExexutor("RENAMENX", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		newkey, err := nextStringArgument(cmd, "newkey", args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Rename(conn, key, newkey, RenameOption{NX: true})
	})

	server.RegisterExexutor("TTL", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.TTL(conn, key)
	})

	server.RegisterExexutor("SCAN", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		cursor, err := nextIntegerArgument(cmd, "cursor", args)
		if err != nil {
			return nil, err
		}
		opt := ScanOption{
			MatchPattern: nil,
			Count:        0,
			Type:         Scan,
		}
		return server.userCommandHandler.Scan(conn, cursor, opt)
	})

	// String commands.

	server.RegisterExexutor("GET", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Get(conn, key)
	})

	server.RegisterExexutor("SET", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		opt := newDefaultSetOption()
		key, val, err := nextSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Set(conn, key, val, opt)
	})

	server.RegisterExexutor("GETSET", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		opt := newDefaultSetOption()
		opt.GET = true
		key, val, err := nextSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Set(conn, key, val, opt)
	})

	server.RegisterExexutor("MSET", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		opt := newDefaultSetOption()
		dict, err := nextMSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		for key, val := range dict {
			if _, err := server.userCommandHandler.Set(conn, key, val, opt); err != nil {
				return nil, err
			}
		}
		return NewOKMessage(), nil
	})

	server.RegisterExexutor("MSETNX", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		dict, err := nextMSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		for key := range dict {
			res, err := server.userCommandHandler.Get(conn, key)
			if err != nil {
				return nil, err
			}
			if !res.IsNil() {
				return NewIntegerMessage(0), nil
			}
		}
		opt := newDefaultSetOption()
		opt.NX = true
		for key, val := range dict {
			if _, err := server.userCommandHandler.Set(conn, key, val, opt); err != nil {
				return nil, err
			}
		}
		return NewIntegerMessage(1), nil
	})

	server.RegisterExexutor("MGET", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		keys, err := nextMGetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		arrayMsg := NewArrayMessage()
		array, _ := arrayMsg.Array()
		for _, key := range keys {
			msg, err := server.userCommandHandler.Get(conn, key)
			if err != nil {
				return nil, err
			}
			array.Append(msg)
		}
		return arrayMsg, nil
	})

	server.RegisterExexutor("SETNX", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		opt := newDefaultSetOption()
		opt.NX = true
		key, val, err := nextSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Set(conn, key, val, opt)
	})

	// Hash commands.

	server.RegisterExexutor("HDEL", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		keys, err := nextKeysArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HDel(conn, hash, keys)
	})

	server.RegisterExexutor("HGET", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HGet(conn, hash, key)
	})

	server.RegisterExexutor("HGETALL", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HGetAll(conn, hash)
	})

	server.RegisterExexutor("HSET", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
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
		return server.userCommandHandler.HSet(conn, hash, key, val, opt)
	})

	server.RegisterExexutor("HSETNX", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
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
		return server.userCommandHandler.HSet(conn, hash, key, val, opt)
	})

	server.RegisterExexutor("HMSET", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		dict, err := nextMSetArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		hsetOpt := HSetOption{
			NX: false,
		}
		for field, val := range dict {
			if _, err := server.userCommandHandler.HSet(conn, hash, field, val, hsetOpt); err != nil {
				return nil, err
			}
		}
		return NewOKMessage(), nil
	})

	server.RegisterExexutor("HMGET", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		hash, err := nextHashArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		fields, err := nextKeysArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		arrayMsg := NewArrayMessage()
		array, _ := arrayMsg.Array()
		for _, field := range fields {
			msg, err := server.userCommandHandler.HGet(conn, hash, field)
			if err != nil {
				return nil, err
			}
			array.Append(msg)
		}
		return arrayMsg, nil
	})

	// List commands.

	server.RegisterExexutor("LINDEX", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		idx, err := nextIntegerArgument(cmd, "index", args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.LIndex(conn, key, idx)
	})

	server.RegisterExexutor("LLEN", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.LLen(conn, key)
	})

	server.RegisterExexutor("LPOP", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, cnt, err := nextPopArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.LPop(conn, key, cnt)
	})

	server.RegisterExexutor("LPUSH", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, elems, err := nextPushArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		opt := PushOption{X: false}
		return server.userCommandHandler.LPush(conn, key, elems, opt)
	})

	server.RegisterExexutor("LPUSHX", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, elems, err := nextPushArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		opt := PushOption{X: true}
		return server.userCommandHandler.LPush(conn, key, elems, opt)
	})

	server.RegisterExexutor("LRANGE", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
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
		return server.userCommandHandler.LRange(conn, key, start, end)
	})

	server.RegisterExexutor("RPOP", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, cnt, err := nextPopArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.RPop(conn, key, cnt)
	})

	server.RegisterExexutor("RPUSH", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, elems, err := nextPushArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		opt := PushOption{X: false}
		return server.userCommandHandler.RPush(conn, key, elems, opt)
	})

	server.RegisterExexutor("RPUSHX", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, elems, err := nextPushArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		opt := PushOption{X: true}
		return server.userCommandHandler.RPush(conn, key, elems, opt)
	})

	// Set commands.

	server.RegisterExexutor("SADD", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		members, err := nextStringArrayArguments(cmd, "member", args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.SAdd(conn, key, members)
	})

	server.RegisterExexutor("SMEMBERS", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.SMembers(conn, key)
	})

	server.RegisterExexutor("SREM", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		members, err := nextStringArrayArguments(cmd, "member", args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.SRem(conn, key, members)
	})

	// ZSet commands.

	server.RegisterExexutor("ZADD", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}

		opt := ZAddOption{
			XX:   false,
			NX:   false,
			LT:   false,
			GT:   false,
			CH:   false,
			INCR: false,
		}

		var score float64
		param, err := args.NextString()
		for err == nil {
			isOption := true
			switch strings.ToUpper(param) {
			case "NX":
				opt.NX = true
			case "XX":
				opt.XX = true
			case "GT":
				opt.GT = true
			case "LT":
				opt.LT = true
			case "CH":
				opt.CH = true
			case "INCR":
				opt.INCR = true
			default:
				score, err = strconv.ParseFloat(param, 64)
				isOption = false
			}
			if !isOption {
				break
			}
		}
		if err != nil {
			return nil, newMissingArgumentError(cmd, "score", err)
		}

		members := []*ZSetMember{}
		member, err := args.NextString()
		if err != nil {
			err = newMissingArgumentError(cmd, "member", err)
		}
		for err == nil {
			members = append(members, &ZSetMember{Score: score, Member: member})
			score, err = nextScoreArgument(cmd, "score", args)
			if err != nil {
				break
			}
			member, err = nextStringArgument(cmd, "member", args)
			if err != nil {
				break
			}
		}
		if !errors.Is(err, proto.ErrEOM) {
			return nil, err
		}

		return server.userCommandHandler.ZAdd(conn, key, members, opt)
	})

	server.RegisterExexutor("ZINCRBY", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		inc, err := nextScoreArgument(cmd, "increment", args)
		if err != nil {
			return nil, err
		}
		member, err := nextStringArgument(cmd, "member", args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.ZIncBy(conn, key, inc, member)
	})

	server.RegisterExexutor("ZRANGE", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}

		start, startEx, err := nextRangeScoreIndexArgument(cmd, "start", args)
		if err != nil {
			return nil, err
		}

		stop, stopEx, err := nextRangeScoreIndexArgument(cmd, "stop", args)
		if err != nil {
			return nil, err
		}

		opt, err := nextRangeOptionArguments(cmd, args)
		if err != nil {
			return nil, err
		}

		if opt.BYSCORE {
			opt.MINEXCLUSIVE = startEx
			opt.MAXEXCLUSIVE = stopEx
			return server.userCommandHandler.ZRangeByScore(conn, key, start, stop, opt)
		}

		return server.userCommandHandler.ZRange(conn, key, int(start), int(stop), opt)
	})

	server.RegisterExexutor("ZREVRANGE", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}

		start, err := nextRangeIndexArgument(cmd, "start", args)
		if err != nil {
			return nil, err
		}

		stop, err := nextRangeIndexArgument(cmd, "stop", args)
		if err != nil {
			return nil, err
		}

		opt, err := nextRangeOptionArguments(cmd, args)
		if err != nil {
			return nil, err
		}

		msg, err := server.userCommandHandler.ZRange(conn, key, start, stop, opt)
		if err != nil {
			return msg, err
		}

		array, err := msg.Array()
		if err != nil {
			return msg, err
		}

		if opt.WITHSCORES {
			return NewArrayMessageWithArray(array.ReverseBy(2)), nil
		}
		return NewArrayMessageWithArray(array.Reverse()), nil
	})

	server.RegisterExexutor("ZRANGEBYSCORE", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}

		min, minEx, err := nextRangeScoreIndexArgument(cmd, "min", args)
		if err != nil {
			return nil, err
		}

		max, maxEx, err := nextRangeScoreIndexArgument(cmd, "max", args)
		if err != nil {
			return nil, err
		}

		opt, err := nextRangeOptionArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		opt.MINEXCLUSIVE = minEx
		opt.MAXEXCLUSIVE = maxEx

		return server.userCommandHandler.ZRangeByScore(conn, key, min, max, opt)
	})

	server.RegisterExexutor("ZREVRANGEBYSCORE", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}

		max, maxEx, err := nextRangeScoreIndexArgument(cmd, "max", args)
		if err != nil {
			return nil, err
		}

		min, minEx, err := nextRangeScoreIndexArgument(cmd, "min", args)
		if err != nil {
			return nil, err
		}

		opt, err := nextRangeOptionArguments(cmd, args)
		if err != nil {
			return nil, err
		}
		opt.MINEXCLUSIVE = minEx
		opt.MAXEXCLUSIVE = maxEx

		msg, err := server.userCommandHandler.ZRangeByScore(conn, key, min, max, opt)
		if err != nil {
			return msg, err
		}

		array, err := msg.Array()
		if err != nil {
			return msg, err
		}

		if opt.WITHSCORES {
			return NewArrayMessageWithArray(array.ReverseBy(2)), nil
		}
		return NewArrayMessageWithArray(array.Reverse()), nil
	})

	server.RegisterExexutor("ZREM", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		members, err := nextStringArrayArguments(cmd, "member", args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.ZRem(conn, key, members)
	})

	server.RegisterExexutor("ZSCORE", func(conn *Conn, cmd string, args Arguments) (*Message, error) {
		key, err := nextKeyArgument(cmd, args)
		if err != nil {
			return nil, err
		}
		member, err := nextStringArgument(cmd, "member", args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.ZScore(conn, key, member)
	})
}
