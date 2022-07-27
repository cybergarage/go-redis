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

package redistest

import (
	"fmt"

	goredis "github.com/go-redis/redis"
)

// Client represents a client for the Redis server.
type Client struct {
	*goredis.Client
}

// NewClient returns a client instance.
func NewClient() *Client {
	client := &Client{
		Client: nil,
	}
	return client
}

// Open opens a connection with the specified host.
func (client *Client) Open(host string) error {
	client.Client = goredis.NewClient(&goredis.Options{
		Addr: fmt.Sprintf("%s:%d", host, DefaultPort),
	})
	status := client.Ping()
	if err := status.Err(); err != nil {
		return err
	}
	return nil
}

// Close closes the current connection with the specified host.
func (client *Client) Close() error {
	if client.Client == nil {
		return nil
	}
	return client.Client.Close()
}
