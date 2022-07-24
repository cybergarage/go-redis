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
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

// Paser represents a Redis serialization protocol (RESP) parser.
type Parser struct {
	reader io.Reader
}

// NewParserWithReader returns a new parser for the specified reader.
func NewParserWithReader(msgReader io.Reader) *Parser {
	Parser := &Parser{
		reader: msgReader,
	}
	return Parser
}

// NewParserWithBytes returns a new parser for the specified bytes.
func NewParserWithBytes(msgBytes []byte) *Parser {
	return NewParserWithReader(bytes.NewBuffer(msgBytes))
}

// nextLineBytes gets a next line bytes.
func (parser *Parser) nextLineBytes() ([]byte, error) {
	var readBytes bytes.Buffer
	readByte := make([]byte, 1)

	// Gets a message bytes.
	n, err := parser.reader.Read(readByte)
	for n == 1 && err == nil && readByte[0] != cr {
		readBytes.WriteByte(readByte[0])
		n, err = parser.reader.Read(readByte)
	}
	if errors.Is(err, io.EOF) {
		return readBytes.Bytes(), nil
	}

	// Skips a next line field.
	parser.reader.Read(readByte)

	return readBytes.Bytes(), nil
}

// nextBulkMessage gets a next line bytes.
func (parser *Parser) nextBulkMessage() (*Message, error) {
	numBytes, err := parser.nextLineBytes()
	if err != nil {
		return nil, err
	}
	num, err := strconv.Atoi(string(numBytes))
	if err != nil {
		return nil, err
	}

	msg, err := newMessageWithTypeByte(arrayMessageByte)
	if err != nil {
		return nil, err
	}
	if num < 0 {
		return msg, nil
	}

	msg.bytes, err = parser.nextLineBytes()
	if err != nil {
		return nil, err
	}
	if len(msg.bytes) != num {
		return msg, fmt.Errorf(errorInvalidBulkStringLength, len(msg.bytes), num)
	}
	return msg, nil
}

// nextFirstArrayMessage gets a first message in the next array.
func (parser *Parser) nextFirstArrayMessage() (*Message, error) {
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
	return parser.Next()
}

// Next returns a next message.
func (parser *Parser) Next() (*Message, error) {
	// Parses a first type byte.
	typeByte := make([]byte, 1)
	_, err := parser.reader.Read(typeByte)
	if errors.Is(err, io.EOF) {
		return nil, nil
	}

	// Returns a next array if the message type is array.
	if typeByte[0] == arrayMessageByte {
		return parser.nextFirstArrayMessage()
	}

	// Returns a next bulk strings if the message type is bulk string.
	if typeByte[0] == bulkMessageByte {
		return parser.nextBulkMessage()
	}

	// Returns a next line bytes
	msg, err := newMessageWithTypeByte(typeByte[0])
	if err != nil {
		return nil, err
	}
	msg.bytes, err = parser.nextLineBytes()
	if err != nil {
		return nil, err
	}
	return msg, nil
}
