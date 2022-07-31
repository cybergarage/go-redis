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
	"fmt"
	"strings"
	"time"

	"github.com/cybergarage/go-redis/redis/proto"
)

type cmdArgs = *proto.Array

// handleCommand handles a client command message.
func (server *Server) handleCommand(ctx *DBContext, cmd string, args cmdArgs) (*Message, error) {
	var resMsg *Message
	var err error
	now := time.Now()

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
		key, err := args.NextString()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "key", err)
		}
		val, err := args.NextString()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "value", err)
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
		key, err := args.NextString()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "key", err)
		}
		val, err := args.NextString()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "value", err)
		}
		return server.userCommandHandler.Set(ctx, key, val, opt)
	case "GET": // 1.0.0
		opt := GetOption{}
		key, err := args.NextString()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "key", err)
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
		key, err := args.NextString()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "key", err)
		}
		val, err := args.NextString()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "value", err)
		}
		return server.userCommandHandler.Set(ctx, key, val, opt)
	case "HSET": // 2.0.0
		opt := HSetOption{}
		hash, err := args.NextString()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "hash", err)
		}
		key, err := args.NextString()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "key", err)
		}
		val, err := args.NextString()
		if err != nil {
			return nil, newMissingArgumentError(cmd, "value", err)
		}
		return server.userCommandHandler.HSet(ctx, hash, key, val, opt)
	default:
		resMsg = NewErrorMessage(fmt.Errorf(errorNotSupportedCommand, cmd))
	}

	return resMsg, err
}
