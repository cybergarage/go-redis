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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cybergarage/go-redis/redis/proto"
)

type cmdArgs = *proto.Array

// handleCommand handles a client command message.
// nolint: gocyclo, maintidx
func (server *Server) handleCommand(ctx *DBContext, cmd string, args cmdArgs) (*Message, error) {
	var resMsg *Message
	var err error
	now := time.Now()

	// Handles system commands.

	switch strings.ToUpper(cmd) {
	case "PING": // 1.0.0
		arg := ""
		if msg, _ := args.Next(); msg != nil {
			arg, err = msg.String()
			if err != nil {
				return nil, err
			}
		}
		return server.systemCmdHandler.Ping(ctx, arg)
	case "ECHO": // 1.0.0
		msg, err := args.NextString()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "msg", err)
		}
		return server.systemCmdHandler.Echo(ctx, msg)
	case "SELECT": // 1.0.0
		id, err := args.NextInteger()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "id", err)
		}
		return server.systemCmdHandler.Select(ctx, id)
	case "QUIT": // 1.0.0
		return server.systemCmdHandler.Quit(ctx)
	}

	if server.userCommandHandler == nil {
		return NewErrorMessage(fmt.Errorf(errorNotSupportedCommand, cmd)), nil
	}

	// Parser commands for user commands.

	parseHashArg := func(args cmdArgs) (string, error) {
		hash, err := args.NextString()
		if err != nil {
			return "", newMissingArgumentError(cmd, "hash", err)
		}
		return hash, nil
	}

	parseKeyArg := func(args cmdArgs) (string, error) {
		key, err := args.NextString()
		if err != nil {
			return "", newMissingArgumentError(cmd, "key", err)
		}
		return key, nil
	}

	parseSetArgs := func(args cmdArgs) (string, string, error) {
		key, err := args.NextString()
		if err != nil {
			return "", "", newMissingArgumentError(cmd, "key", err)
		}
		val, err := args.NextString()
		if err != nil {
			return "", "", newMissingArgumentError(cmd, "value", err)
		}
		return key, val, err
	}

	parseMSetArgs := func(args cmdArgs) (map[string]string, error) {
		var key, val string
		var err error
		dir := map[string]string{}
		key, err = args.NextString()
		for err == nil {
			val, err = args.NextString()
			if err != nil {
				newMissingArgumentError(cmd, key, err)
			}
			dir[key] = val
			key, err = args.NextString()
		}
		if !errors.Is(err, proto.ErrEOM) {
			return nil, err
		}
		return dir, nil
	}

	parseMGetArgs := func(args cmdArgs) ([]string, error) {
		var key string
		var err error
		keys := []string{}
		key, err = args.NextString()
		for err == nil {
			keys = append(keys, key)
			key, err = args.NextString()
		}
		if !errors.Is(err, proto.ErrEOM) {
			return nil, err
		}
		return keys, nil
	}

	// Handles user commands.

	switch strings.ToUpper(cmd) {
	case "SET": // 1.0.0
		opt := SetOption{
			NX:      false,
			XX:      false,
			EX:      0,
			PX:      0,
			EXAT:    now,
			PXAT:    now,
			KEEPTTL: false,
			GET:     false,
		}
		key, val, err := parseSetArgs(args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Set(ctx, key, val, opt)
	case "SETNX": // 1.0.0
		opt := SetOption{
			NX:      true,
			XX:      false,
			EX:      0,
			PX:      0,
			EXAT:    now,
			PXAT:    now,
			KEEPTTL: false,
			GET:     false,
		}
		key, val, err := parseSetArgs(args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Set(ctx, key, val, opt)
	case "GET": // 1.0.0
		opt := GetOption{}
		key, err := parseKeyArg(args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Get(ctx, key, opt)
	case "GETSET": // 1.0.0
		opt := SetOption{
			NX:      false,
			XX:      false,
			EX:      0,
			PX:      0,
			EXAT:    now,
			PXAT:    now,
			KEEPTTL: false,
			GET:     true,
		}
		key, val, err := parseSetArgs(args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.Set(ctx, key, val, opt)
	case "HSET": // 2.0.0
		opt := HSetOption{}
		hash, err := parseHashArg(args)
		if err != nil {
			return nil, err
		}
		key, val, err := parseSetArgs(args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HSet(ctx, hash, key, val, opt)
	case "HGET": // 2.0.0
		opt := HGetOption{}
		hash, err := parseHashArg(args)
		if err != nil {
			return nil, err
		}
		key, err := parseKeyArg(args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HGet(ctx, hash, key, opt)
	case "HMSET": // 2.0.0
		opt := HMSetOption{}
		hash, err := parseHashArg(args)
		if err != nil {
			return nil, err
		}
		dir, err := parseMSetArgs(args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HMSet(ctx, hash, dir, opt)
	case "HMGET": // 2.0.0
		opt := HMGetOption{}
		hash, err := parseHashArg(args)
		if err != nil {
			return nil, err
		}
		keys, err := parseMGetArgs(args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.HMGet(ctx, hash, keys, opt)
	case "MSET": // 1.0.1
		opt := MSetOption{
			NX: false,
		}
		dir, err := parseMSetArgs(args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.MSet(ctx, dir, opt)
	case "MSETNX": // 1.0.1
		opt := MSetOption{
			NX: true,
		}
		dir, err := parseMSetArgs(args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.MSet(ctx, dir, opt)
	case "MGET": // 1.0.1
		opt := MGetOption{}
		keys, err := parseMGetArgs(args)
		if err != nil {
			return nil, err
		}
		return server.userCommandHandler.MGet(ctx, keys, opt)
	default:
		resMsg = NewErrorMessage(fmt.Errorf(errorNotSupportedCommand, cmd))
	}

	return resMsg, err
}
