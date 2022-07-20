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

// Paser represents a Redis serialization protocol (RESP) parser.
type Parser struct {
	readBuffer    []byte
	readBufferLen int
	readIndex     int
}

// NewParser returns a new parser instance.
func NewParser() *Parser {
	Parser := &Parser{
		readBuffer:    nil,
		readBufferLen: 0,
		readIndex:     0,
	}
	return Parser
}

// Parse parses a serialized request binary from the client.
func (parser *Parser) Parse(protoBytes []byte) error {
	if len(protoBytes) == 0 {
		return fmt.Errorf(errorEmptyMessage, 0)
	}
	parser.readBuffer = protoBytes
	parser.readBufferLen = len(protoBytes)
	parser.readIndex = 0
	return nil
}

// Next returns a next message.
func (parser *Parser) Next() (*Message, error) {
	// Finishes when all bytes have been read.
	if parser.readBufferLen <= parser.readIndex {
		return nil, nil
	}

	// Parses a first type byte.
	typeByte := parser.readBuffer[parser.readIndex]
	msg, err := newMessageWithTypeByte(typeByte)
	if err != nil {
		return nil, err
	}

	// Gets a message bytes.
	parser.readIndex++
	startIndex := parser.readIndex
	for (parser.readIndex < parser.readBufferLen) && (parser.readBuffer[parser.readIndex] != cr) {
		parser.readIndex++
	}
	if parser.readBufferLen <= parser.readIndex {
		return nil, fmt.Errorf(errorInvalidMessage, string(parser.readBuffer))
	}
	msg.Bytes = parser.readBuffer[startIndex:parser.readIndex]

	// Skips a next line field.
	parser.readIndex++
	for (parser.readIndex < parser.readBufferLen) && (parser.readBuffer[parser.readIndex] != lf) {
		parser.readIndex++
	}
	if parser.readBufferLen <= parser.readIndex {
		return nil, fmt.Errorf(errorInvalidMessage, string(parser.readBuffer))
	}
	parser.readIndex++

	return msg, nil
}
