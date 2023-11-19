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
)

func (server *Server) Auth(conn *Conn, username string, password string) (*Message, error) {
	required, configPassword := server.ConfigRequirePass()
	if required && password != configPassword {
		return nil, errors.New("invalid password")
	}
	conn.SetAuthrized(true)
	return NewOKMessage(), nil
}
