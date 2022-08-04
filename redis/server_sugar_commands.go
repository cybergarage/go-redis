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

// nolint: gocyclo, maintidx
func (server *Server) registerSugarExecutors() {
	// Registers string commands.

	server.RegisterExexutor("APPEND", func(ctx *DBContext, cmd string, args Arguments) (*Message, error) {
		key, appendVal, err := parseSetArgs(cmd, args)
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
}