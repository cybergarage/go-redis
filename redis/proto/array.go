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

package proto

import (
	"bytes"
	"strconv"
)

// Array represents a array message.
type Array struct {
	index int
	msgs  []*Message
}

// NewArray returns a new array message.
func NewArray() *Array {
	array := &Array{
		index: 0,
		msgs:  make([]*Message, 0),
	}
	return array
}

// newArrayWithParser returns a new array message.
func newArrayWithParser(parser *Parser) (*Array, error) {
	numBytes, err := parser.nextLineBytes()
	if err != nil {
		return nil, err
	}
	num, err := strconv.Atoi(string(numBytes))
	if err != nil {
		return nil, err
	}
	if num < 0 {
		return NewArray(), nil
	}

	// Gets all array messages
	msgs := make([]*Message, num)
	for n := 0; n < num; n++ {
		msg, err := parser.Next()
		if err != nil {
			return nil, err
		}
		msgs[n] = msg
	}
	array := &Array{
		index: 0,
		msgs:  msgs,
	}

	return array, nil
}

// Size returns the array size.
func (array *Array) Size() int {
	return len(array.msgs)
}

// Next returns a next message.
func (array *Array) Next() (*Message, error) {
	if array.Size() <= array.index {
		return nil, nil
	}
	msg := array.msgs[array.index]
	array.index++
	return msg, nil
}

// NextMessages returns all unread messages.
func (array *Array) NextMessages() ([]*Message, error) {
	unreadMsgCnt := array.Size() - array.index
	if unreadMsgCnt <= 0 {
		return []*Message{}, nil
	}
	unreadMsgs := make([]*Message, unreadMsgCnt)
	for n := 0; n < unreadMsgCnt; n++ {
		msg, err := array.Next()
		if err != nil {
			return nil, err
		}
		unreadMsgs[n] = msg
	}
	return unreadMsgs, nil
}

// RESPBytes returns the RESP byte representation.
func (array *Array) RESPBytes() []byte {
	var respBytes bytes.Buffer
	return respBytes.Bytes()
}
