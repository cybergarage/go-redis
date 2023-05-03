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

	"github.com/cybergarage/go-redis/redis/proto"
)

type Arguments = *proto.Array
type Executor func(*Conn, string, Arguments) (*Message, error)
type Executors map[string]Executor

// executeCommand handles a client command message.
func (server *Server) executeCommand(conn *Conn, cmd string, args Arguments) (*Message, error) {
	if server.userCommandHandler == nil {
		return NewErrorMessage(fmt.Errorf(errorNotSupportedCommand, cmd)), nil
	}

	upperCmd := strings.ToUpper(cmd)
	cmdExecutor, ok := server.commandExecutors[upperCmd]
	if !ok {
		return NewErrorMessage(fmt.Errorf(errorNotSupportedCommand, cmd)), nil
	}

	s := conn.SpanContext.Span().StartSpan(upperCmd)
	defer s.Span().Finish()

	return cmdExecutor(conn, cmd, args)
}
