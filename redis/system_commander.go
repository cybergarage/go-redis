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

func (server *Server) Ping(ctx *DBContext, arg string) (*Message, error) {
	if len(arg) == 0 {
		return NewStringMessage("PONG"), nil
	}
	return NewBulkMessage(arg), nil
}

func (server *Server) Echo(ctx *DBContext, arg string) (*Message, error) {
	return NewBulkMessage(arg), nil
}

func (server *Server) Select(ctx *DBContext, index int) (*Message, error) {
	ctx.id = index
	return NewOKMessage(), nil
}

func (server *Server) Quit(ctx *DBContext) (*Message, error) {
	return NewOKMessage(), errQuit
}

func (server *Server) ConfigSet(ctx *DBContext, params map[string]string) (*Message, error) {
	for key, param := range params {
		server.SetConfig(key, param)
	}
	return NewOKMessage(), nil
}

func (server *Server) ConfigGet(ctx *DBContext, keys []string) (*Message, error) {
	msg := NewArrayMessage()
	for _, key := range keys {
		msg.Append(NewBulkMessage(key))
		param, ok := server.ConfigParameter(key)
		if ok {
			msg.Append(NewBulkMessage(param))
		} else {
			msg.Append(NewBulkMessage(""))
		}
	}
	return msg, nil
}
