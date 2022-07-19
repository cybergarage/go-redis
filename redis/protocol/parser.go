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

import (
	"fmt"
)

// Paser represents a Redis serialization protocol (RESP) parser.
type Parser struct {
	readBuffer    []byte
	readIndex     int
	readBufferLen int
}

// NewParser returns a new parser instance.
func NewParser() *Parser {
	Parser := &Parser{
		readBuffer: nil,
		readIndex:  0,
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
	return nil
}

// Next returns a next message.
func (parser *Parser) Next() (*Message, error) {
	return nil, nil
}
