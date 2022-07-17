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
	"bytes"
	"fmt"
)

// Paser represents a Redis serialization protocol (RESP) parser.
type Parser struct {
}

// NewParser returns a new parser instance.
func NewParser() *Parser {
	Parser := &Parser{}
	return Parser
}

// Parse parses a serialized request binary from the client.
func (parser *Parser) Paerse(protoBytes []byte) ([]Message, error) {
	if len(protoBytes) == 0 {
		return nil, fmt.Errorf(errorEmptyMessage, len(protoBytes))
	}
	allMsgBytes := bytes.Split(protoBytes, []byte(crlr))
	if msgCount := len(allMsgBytes); msgCount == 0 {
		return nil, fmt.Errorf(errorEmptyMessage, msgCount)
	}
	var msgs []Message
	for _, msgBytes := range allMsgBytes {
		if len(msgBytes) == 0 {
			continue
		}
		_, ok := parseMessageType(msgBytes[0])
		if !ok {
			return nil, fmt.Errorf(errorUnknownMessageType, rune(msgBytes[0]))
		}
	}
	return msgs, nil
}
