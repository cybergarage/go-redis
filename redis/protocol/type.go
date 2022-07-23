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

package protocol

// MessageType represents a message type of Redis serialization protocol.
type MessageType int

const (
	String MessageType = iota
	Error
	Integer
	Bulk
	Array
)

const (
	stringByte  = byte('+')
	errorByte   = byte('-')
	integerByte = byte(':')
	bulkByte    = byte('$')
	arrayByte   = byte('*')
)

var messageTypes = map[byte]MessageType{
	stringByte:  String,
	errorByte:   Error,
	integerByte: Integer,
	bulkByte:    Bulk,
	arrayByte:   Array,
}

func parseMessageType(b byte) (MessageType, bool) {
	t, ok := messageTypes[b]
	return t, ok
}
