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

// CommandHandler is a hander interface for user commands.
type CommandHandler interface {
	Set(ctx *DBContext, key string, val string, opt SetOption) (*Message, error)
	Get(ctx *DBContext, key string) (*Message, error)
}

// SystemCommandHandler is a hander interface for system commands.
type SystemCommandHandler interface {
	Ping(ctx *DBContext, arg string) (*Message, error)
	Echo(ctx *DBContext, arg string) (*Message, error)
	Select(ctx *DBContext, index int) (*Message, error)
	Quit(ctx *DBContext) (*Message, error)
}
