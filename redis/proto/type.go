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

// MessageType represents a message type of Redis serialization protocol.
type MessageType int

const (
	StringMessage MessageType = iota
	ErrorMessage
	IntegerMessage
	BulkMessage
	ArrayMessage
)

const (
	stringMessageByte  = byte('+')
	errorMessageByte   = byte('-')
	integerMessageByte = byte(':')
	bulkMessageByte    = byte('$')
	arrayMessageByte   = byte('*')
)

var messageTypes = map[byte]MessageType{
	stringMessageByte:  StringMessage,
	errorMessageByte:   ErrorMessage,
	integerMessageByte: IntegerMessage,
	bulkMessageByte:    BulkMessage,
	arrayMessageByte:   ArrayMessage,
}

var messageTypeBytes = map[MessageType]byte{
	StringMessage:  stringMessageByte,
	ErrorMessage:   errorMessageByte,
	IntegerMessage: integerMessageByte,
	BulkMessage:    bulkMessageByte,
	ArrayMessage:   arrayMessageByte,
}

func parseMessageType(b byte) (MessageType, bool) {
	t, ok := messageTypes[b]
	return t, ok
}

func messageTypeToByte(t MessageType) (byte, bool) {
	b, ok := messageTypeBytes[t]
	return b, ok
}
