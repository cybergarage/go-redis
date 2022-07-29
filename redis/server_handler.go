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
	"fmt"
	"strings"

	"github.com/cybergarage/go-redis/redis/proto"
)

type args = *proto.Array

// handleCommand handles a client command message.
func (server *Server) handleCommand(cmd string, cmdArgs args) (*Message, error) {
	args, err := cmdArgs.NextMessages()
	if err != nil {
		return nil, err
	}

	var resMsg *Message

	switch strings.ToUpper(cmd) {
	case "PING": // 1.0.0
		arg := ""
		if 0 < len(args) {
			arg, err = args[0].String()
			if err != nil {
				return nil, err
			}
		}
		return server.systemCmdHandler.Ping(arg)
	default:
		resMsg = NewErrorMessage(fmt.Errorf(errorNotSupportedCommand, cmd))
	}

	if server.CommandHandler == nil {
		return NewErrorMessage(fmt.Errorf(errorNotSupportedCommand, cmd)), nil
	}

	switch strings.ToUpper(cmd) {
	case "SET": // 1.0.0
	case "GET": // 1.0.0
	default:
		resMsg = NewErrorMessage(fmt.Errorf(errorNotSupportedCommand, cmd))
	}

	return resMsg, err
}
