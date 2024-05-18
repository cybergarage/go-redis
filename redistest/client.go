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

package redistest

import (
	"fmt"
	"time"

	goredis "github.com/go-redis/redis"
)

// Client represents a client for the Redis server.
type Client struct {
	*goredis.Client
}

// ClientOptions represents a client options for the Redis server.
type ClientOptions = goredis.Options

// NewClientOptions returns a default client options.
// nolint: exhaustivestruct
func NewClientOptions() ClientOptions {
	return goredis.Options{ // nolint:exhaustruct
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
		DB:           1,
	}
}

// NewClient returns a client instance.
func NewClient() *Client {
	// nolint: exhaustivestruct
	client := &Client{
		Client: nil,
	}
	return client
}

// Open opens a connection with the specified host.
func (client *Client) Open(host string) error {
	opts := NewClientOptions()
	return client.OpenWith(host, DefaultPort, &opts)
}

func (client *Client) OpenWith(host string, port int, opts *ClientOptions) error {
	opts.Addr = fmt.Sprintf("%s:%d", host, port)
	client.Client = goredis.NewClient(opts)
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
