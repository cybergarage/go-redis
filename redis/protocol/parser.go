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
	"strconv"
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

// nextLineBytes gets a next line bytes.
func (parser *Parser) nextLineBytes() ([]byte, error) {
	// Gets a message bytes.
	startIndex := parser.readIndex
	for (parser.readIndex < parser.readBufferLen) && (parser.readBuffer[parser.readIndex] != cr) {
		parser.readIndex++
	}

	// a next carriage return filed is not found, and all bytes have been read.
	if parser.readBufferLen <= parser.readIndex {
		return parser.readBuffer[startIndex:parser.readBufferLen], nil
	}

	lineBytes := parser.readBuffer[startIndex:parser.readIndex]

	// Skips a next line field.
	parser.readIndex++
	if parser.readBufferLen <= parser.readIndex { // a next line filed is not found.
		return lineBytes, nil
	}
	parser.readIndex++

	return lineBytes, nil
}

// nextBulkStrings gets a next line bytes.
func (parser *Parser) nextBulkStrings() ([]byte, error) {
	numBytes, err := parser.nextLineBytes()
	if err != nil {
		return nil, err
	}
	num, err := strconv.Atoi(string(numBytes))
	if err != nil {
		return nil, err
	}
	if num < 0 {
		return nil, nil
	}
	bulkBytes, err := parser.nextLineBytes()
	if err != nil {
		return nil, err
	}
	if len(bulkBytes) != num {
		return bulkBytes, fmt.Errorf(errorInvalidBulkStringLength, len(bulkBytes), num)
	}
	return bulkBytes, nil
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
	parser.readIndex++

	// Gets a next bulk strings
	if typeByte == bulkByte {
		msg.Bytes, err = parser.nextBulkStrings()
		if err != nil {
			return nil, err
		}
		return msg, nil
	}

	// Gets a next line bytes
	msg.Bytes, err = parser.nextLineBytes()
	if err != nil {
		return nil, err
	}
	return msg, nil
}
