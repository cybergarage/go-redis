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

package server

import (
	"github.com/cybergarage/go-redis/redis"
)

type List []string

func (server *Server) LPush(ctx *redis.DBContext, key string, elements []string, opt redis.PushOption) (*redis.Message, error) {
	return nil, nil
}

func (server *Server) RPush(ctx *redis.DBContext, key string, elements []string, opt redis.PushOption) (*redis.Message, error) {
	return nil, nil
}

func (server *Server) LPop(ctx *redis.DBContext, key string, count int) (*redis.Message, error) {
	return nil, nil
}

func (server *Server) RPop(ctx *redis.DBContext, key string, count int) (*redis.Message, error) {
	return nil, nil
}

func (server *Server) LRange(ctx *redis.DBContext, key string, start int, stop int) (*redis.Message, error) {
	return nil, nil
}

func (server *Server) LIndex(ctx *redis.DBContext, key string, index int) (*redis.Message, error) {
	return nil, nil
}

func (server *Server) LLen(ctx *redis.DBContext, key string, index int) (*redis.Message, error) {
	return nil, nil
}
