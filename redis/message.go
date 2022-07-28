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
	"github.com/cybergarage/go-redis/redis/proto"
)

// Message represents a message of Redis serialization protocol.
type Message = proto.Message

// NewStringMessage creates a string message.
func NewStringMessage(msg string) *Message {
	return proto.NewMessageWithType(proto.StringMessage).SetBytes([]byte(msg))
}

// NewBulkMessage creates a bulk string message.
func NewBulkMessage(msg string) *Message {
	return proto.NewMessageWithType(proto.BulkMessage).SetBytes([]byte(msg))
}
