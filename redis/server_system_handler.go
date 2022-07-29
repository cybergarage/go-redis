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

func (server *Server) Ping(arg string) (*Message, error) {
	if len(arg) == 0 {
		return NewStringMessage("PONG"), nil
	}
	return NewBulkMessage(arg), nil
}

func (server *Server) Echo(arg string) (*Message, error) {
	return NewBulkMessage(arg), nil
}

func (server *Server) Quit() (*Message, error) {
	return NewStringMessage("OK"), errQuit
}
