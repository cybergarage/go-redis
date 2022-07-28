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

package proto

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

// Message represents a message of Redis serialization protocol.
type Message struct {
	Type  MessageType
	bytes []byte
	array *Array
}

// NewMessageWithType returns a new message instance with the specified type.
func NewMessageWithType(t MessageType) *Message {
	msg := &Message{
		Type:  t,
		bytes: nil,
		array: nil,
	}
	return msg
}

// newMessageWithTypeByte returns a new message instance with the specified type byte.
func newMessageWithTypeByte(b byte) (*Message, error) {
	t, ok := parseMessageType(b)
	if !ok {
		return nil, fmt.Errorf(errorUnknownMessageType, b)
	}
	return NewMessageWithType(t), nil
}

// SetBytes sets a byte array to the message.
func (msg *Message) SetBytes(bytes []byte) *Message {
	msg.bytes = bytes
	return msg
}

// SetArray sets an array to the message.
func (msg *Message) SetArray(array *Array) *Message {
	msg.array = array
	return msg
}

// IsType returns true if the message type is the specified type, otherwise false.
func (msg *Message) IsType(t MessageType) bool {
	return msg.Type == t
}

// IsString returns true if the message type is string, otherwise false.
func (msg *Message) IsString() bool {
	return msg.IsType(StringMessage)
}

// IsError returns true if the message type is error, otherwise false.
func (msg *Message) IsError() bool {
	return msg.IsType(ErrorMessage)
}

// IsInteger returns true if the message type is integer, otherwise false.
func (msg *Message) IsInteger() bool {
	return msg.IsType(IntegerMessage)
}

// IsBulk returns true if the message type is bulk, otherwise false.
func (msg *Message) IsBulk() bool {
	return msg.IsType(BulkMessage)
}

// IsArray returns true if the message type is array, otherwise false.
func (msg *Message) IsArray() bool {
	return msg.IsType(ArrayMessage)
}

// Bytes returns the message raw bytes.
func (msg *Message) Bytes() ([]byte, error) {
	return msg.bytes, nil
}

// String returns the message string if the message type is string, otherwise it returns an error.
func (msg *Message) String() (string, error) {
	switch msg.Type {
	case StringMessage:
		return string(msg.bytes), nil
	case ArrayMessage, BulkMessage, ErrorMessage, IntegerMessage:
		return "", fmt.Errorf(errorInvalidMessageType, msg.Type)
	}
	return "", fmt.Errorf(errorInvalidMessageType, msg.Type)
}

// Error returns the message error if the message type is error, otherwise it returns an error.
func (msg *Message) Error() (error, error) {
	switch msg.Type {
	case ErrorMessage:
		return errors.New(string(msg.bytes)), nil
	case StringMessage, ArrayMessage, BulkMessage, IntegerMessage:
		return nil, fmt.Errorf(errorInvalidMessageType, msg.Type)
	}
	return nil, fmt.Errorf(errorInvalidMessageType, msg.Type)
}

// Integer returns the message integer if the message type is integer, otherwise it returns an error.
func (msg *Message) Integer() (int, error) {
	switch msg.Type {
	case IntegerMessage:
		return strconv.Atoi(string(msg.bytes))
	case ArrayMessage, StringMessage, BulkMessage, ErrorMessage:
		return 0, fmt.Errorf(errorInvalidMessageType, msg.Type)
	}
	return 0, fmt.Errorf(errorInvalidMessageType, msg.Type)
}

// Array returns the message array if the message type is array, otherwise it returns an error.
func (msg *Message) Array() (*Array, error) {
	switch msg.Type {
	case ArrayMessage:
		return msg.array, nil
	case IntegerMessage, StringMessage, BulkMessage, ErrorMessage:
		return nil, fmt.Errorf(errorInvalidMessageType, msg.Type)
	}
	return nil, fmt.Errorf(errorInvalidMessageType, msg.Type)
}

// RESPBytes returns the RESP byte representation.
func (msg *Message) RESPBytes() ([]byte, error) {
	var respBytes bytes.Buffer

	switch msg.Type {
	case StringMessage, ErrorMessage, IntegerMessage:
		b, ok := messageTypeToByte(msg.Type)
		if !ok {
			return nil, fmt.Errorf(errorUnknownMessageType, msg.Type)
		}
		respBytes.WriteByte(b)
		respBytes.Write(msg.bytes)
		respBytes.WriteRune(cr)
		respBytes.WriteRune(lf)
	case BulkMessage:
		respBytes.WriteByte(bulkMessageByte)
		switch {
		case msg.bytes == nil:
			respBytes.WriteString("-1")
		case len(msg.bytes) == 0:
			respBytes.WriteString("0")
			respBytes.WriteRune(cr)
			respBytes.WriteRune(lf)
		default:
			respBytes.WriteString(strconv.Itoa(len(msg.bytes)))
			respBytes.WriteRune(cr)
			respBytes.WriteRune(lf)
			respBytes.Write(msg.bytes)
		}
		respBytes.WriteRune(cr)
		respBytes.WriteRune(lf)
	case ArrayMessage:
		array, err := msg.Array()
		if err != nil {
			return nil, err
		}
		bytes, err := array.RESPBytes()
		if err != nil {
			return nil, err
		}
		respBytes.Write(bytes)
	}

	return respBytes.Bytes(), nil
}
