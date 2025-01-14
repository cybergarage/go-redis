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

	"github.com/cybergarage/go-redis/redis/auth"
)

func (server *server) Auth(conn *Conn, username string, password string) (*Message, error) {
	q, err := auth.NewQuery(
		auth.WithQueryUsername(username),
		auth.WithQueryPassword(password),
	)
	if err != nil {
		return nil, err
	}

	ok, err := server.VerifyCredential(conn, q)

	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("invalid username or password")
	}

	conn.SetUserName(username)
	conn.SetPassword(password)
	conn.SetAuthrized(true)
	return NewOKMessage(), nil
}

// SetCredential sets a credential.
func (server *server) SetCredential(cred auth.Credential) {
	server.credStore[cred.Username()] = cred
}

// LookupCredential looks up a credential.
func (server *server) LookupCredential(q auth.Query) (auth.Credential, bool, error) {
	user := q.Username()
	cred, ok := server.credStore[user]
	return cred, ok, nil
}
