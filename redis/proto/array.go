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
	arraySize, err := strconv.Atoi(string(numBytes))
	if err != nil {
		return nil, err
	}
	if arraySize < 0 {
		return NewArray(), nil
	}

	// Gets all array messages
	msgs := make([]*Message, arraySize)
	for n := 0; n < arraySize; n++ {
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

// Append adds a message into the array.
func (array *Array) Append(msg *Message) {
	array.msgs = append(array.msgs, msg)
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

// NextMessage returns the next message if any, otherwise it returns error.
func (array *Array) NextMessage() (*Message, error) {
	msg, err := array.Next()
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, ErrEOM
	}
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

// NextBytes returns the next byte message.
func (array *Array) NextBytes() ([]byte, error) {
	msg, err := array.NextMessage()
	if err != nil {
		return nil, err
	}
	return msg.Bytes()
}

// NextString returns the message string if the message type is string, otherwise it returns an error.
func (array *Array) NextString() (string, error) {
	msg, err := array.NextMessage()
	if err != nil {
		return "", err
	}
	return msg.String()
}

// NextError returns the message error if the message type is error, otherwise it returns an error.
func (array *Array) NextError() (error, error) {
	msg, err := array.NextMessage()
	if err != nil {
		return nil, err
	}
	return msg.Error()
}

// NextInteger returns the message integer if the message type is integer, otherwise it returns an error.
func (array *Array) NextInteger() (int, error) {
	msg, err := array.NextMessage()
	if err != nil {
		return 0, err
	}
	return msg.Integer()
}

// NextArray returns the message array if the message type is array, otherwise it returns an error.
func (array *Array) NextArray() (*Array, error) {
	msg, err := array.NextMessage()
	if err != nil {
		return nil, err
	}
	return msg.Array()
}

// ReverseBy returns the reversed array with the specified step.
func (array *Array) ReverseBy(step int) *Array {
	ra := NewArray()
	l := len(array.msgs)
	for i := 0; i < l; i += step {
		for j := 0; j < step; j++ {
			idx := (l - i - 1) - (step - 1) + j
			ra.msgs = append(ra.msgs, array.msgs[idx])
		}
	}
	return ra
}

// Reverse returns the reversed array.
func (array *Array) Reverse() *Array {
	return array.ReverseBy(1)
}

// RESPBytes returns the RESP byte representation.
func (array *Array) RESPBytes() ([]byte, error) {
	var respBytes bytes.Buffer

	respBytes.WriteByte(arrayMessageByte)

	arraySize := array.Size()
	respBytes.WriteString(strconv.Itoa(arraySize))
	respBytes.WriteRune(cr)
	respBytes.WriteRune(lf)

	for n := 0; n < arraySize; n++ {
		bytes, err := array.msgs[n].RESPBytes()
		if err != nil {
			return respBytes.Bytes(), err
		}
		respBytes.Write(bytes)
	}

	return respBytes.Bytes(), nil
}
