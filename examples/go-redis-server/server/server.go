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

package server

import (
	"github.com/cybergarage/go-redis/redis"
)

// Server represents an example server.
type Server struct {
	*redis.Server
	*Databases
}

// NewServer returns an example server instance.
func NewServer() *Server {
	server := &Server{
		Server:    redis.NewServer(),
		Databases: NewDatabases(),
	}
	server.SetCommandHandler(server)
	return server
}

// GetDatabase returns the database with the specified ID.
func (server *Server) GetDatabase(id int) (*Database, error) {
	db, ok := server.Databases.GetDatabase(id)
	if !ok {
		server.SetDatabase(id, db)
	}
	return db, nil
}
