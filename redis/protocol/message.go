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
	"errors"
	"fmt"
	"strconv"
)

// Message represents a message of Redis serialization protocol.
type Message struct {
	Type  MessageType
	bytes []byte
}

// newMessageWithType returns a new message instance with the specified type.
func newMessageWithType(t MessageType) *Message {
	msg := &Message{
		Type:  t,
		bytes: nil,
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

// Bytes returns the message raw bytes.
func (msg *Message) Bytes() ([]byte, error) {
	return msg.bytes, nil
}

// String returns the message string if the message type is string, otherwise it returns an error.
func (msg *Message) String() (string, error) {
	switch msg.Type {
	case String:
		return string(msg.bytes), nil
	case Array, Bulk, Error, Integer:
		return "", fmt.Errorf(errorInvalidMessageType, msg.Type)
	}
	return "", fmt.Errorf(errorInvalidMessageType, msg.Type)
}

// Error returns the message error if the message type is error, otherwise it returns an error.
func (msg *Message) Error() (error, error) {
	switch msg.Type {
	case Error:
		return errors.New(string(msg.bytes)), nil
	case String, Array, Bulk, Integer:
		return nil, fmt.Errorf(errorInvalidMessageType, msg.Type)
	}
	return nil, fmt.Errorf(errorInvalidMessageType, msg.Type)
}

// Integer returns the message integer if the message type is integer, otherwise it returns an error.
func (msg *Message) Integer() (int, error) {
	switch msg.Type {
	case Integer:
		return strconv.Atoi(string(msg.bytes))
	case Array, String, Bulk, Error:
		return 0, fmt.Errorf(errorInvalidMessageType, msg.Type)
	}
	return 0, fmt.Errorf(errorInvalidMessageType, msg.Type)
}
