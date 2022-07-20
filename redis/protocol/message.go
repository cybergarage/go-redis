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

package protocol

import (
	"fmt"
)

// Message represents a message of Redis serialization protocol.
type Message struct {
	Type  MessageType
	Bytes []byte
}

// newMessageWithType returns a new message instance with the specified type.
func newMessageWithType(t MessageType) *Message {
	msg := &Message{
		Type:  t,
		Bytes: nil,
	}
	return msg
}

// newMessageWithTypeByte returns a new message instance with the specified type byte.
func newMessageWithTypeByte(b byte) (*Message, error) {
	t, ok := parseMessageType(b)
	if !ok {
		return nil, fmt.Errorf(errorUnknownMessageType, b)
	}
	return newMessageWithType(t), nil
}

// String returns the message string value if the message type is string, otherwise it returns an error.
func (msg *Message) String() (string, error) {
	switch msg.Type {
	case SimpleString:
		return string(msg.Bytes), nil
	case Array, BulkString, Error, Integer:
		return "", fmt.Errorf(errorInvalidMessageType, msg.Type)
	}
	return "", fmt.Errorf(errorInvalidMessageType, msg.Type)
}
