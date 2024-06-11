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
	"fmt"
	"testing"
)

func testParserSingleMessages(t *testing.T, msgString string, compare func(*Message, any) (any, bool), expected any, expectedError error) {
	t.Helper()

	parser := NewParserWithBytes([]byte(msgString))
	msg, err := parser.Next()
	if err != nil {
		if expectedError == nil {
			t.Errorf("%s %s", msgString, err)
		} else if err.Error() != expectedError.Error() {
			t.Errorf("Unexpected error message: %s, expecting: %s, msg: %s", err, expectedError, msgString)
		}
		return
	}

	if actual, ok := compare(msg, expected); !ok {
		t.Errorf("%s != %s", actual, expected)
		return
	}

	msgBytes, err := msg.RESPBytes()
	if err != nil {
		t.Errorf("%s %s", msgString, err)
		return
	}

	if !bytes.Equal([]byte(msgString), msgBytes) {
		msgString += string(cr)
		msgString += string(lf)
		if !bytes.Equal([]byte(msgString), msgBytes) {
			t.Errorf("%s \n!=\n%s", string(msgBytes), msgString)
		}
	}

	_, err = parser.Next()
	if err != nil {
		t.Errorf("%s %s", msgString, err)
		return
	}
}

func TestParserStringMessages(t *testing.T) {
	// RESP protocol spec examples.
	respExamples := []struct {
		message  string
		expected string
	}{
		{
			message:  "+OK\r\n",
			expected: "OK",
		},
	}

	compare := func(msg *Message, exp any) (any, bool) {
		expected, ok := exp.(string)
		if !ok {
			return nil, false
		}
		if !msg.IsString() {
			return nil, false
		}
		actual, err := msg.String()
		if err != nil {
			return nil, false
		}
		if actual != expected {
			return actual, false
		}
		return actual, true
	}

	for _, respExample := range respExamples {
		testParserSingleMessages(t, respExample.message, compare, respExample.expected, nil)
	}
}

func TestParserErrorMessages(t *testing.T) {
	// RESP protocol spec examples.
	respExamples := []struct {
		message  string
		expected string
	}{
		{
			message:  "-Error message\r\n",
			expected: "Error message",
		},
		{
			message:  "-ERR unknown command 'helloworld'",
			expected: "ERR unknown command 'helloworld'",
		},
		{
			message:  "-WRONGTYPE Operation against a key holding the wrong kind of value",
			expected: "WRONGTYPE Operation against a key holding the wrong kind of value",
		},
	}

	compare := func(msg *Message, exp any) (any, bool) {
		expected, ok := exp.(string)
		if !ok {
			return nil, false
		}
		if !msg.IsError() {
			return nil, false
		}
		actual, err := msg.Error()
		if err != nil {
			return nil, false
		}
		if actual.Error() != expected {
			return actual, false
		}
		return actual, true
	}

	for _, respExample := range respExamples {
		testParserSingleMessages(t, respExample.message, compare, respExample.expected, nil)
	}
}

func TestParserIntegerMessages(t *testing.T) {
	// RESP protocol spec examples.
	respExamples := []struct {
		message  string
		expected int
	}{
		{
			message:  ":0\r\n",
			expected: 0,
		},
		{
			message:  ":1000\r\n",
			expected: 1000,
		},
	}

	compare := func(msg *Message, exp any) (any, bool) {
		expected, ok := exp.(int)
		if !ok {
			return nil, false
		}
		if !msg.IsInteger() {
			return nil, false
		}
		actual, err := msg.Integer()
		if err != nil {
			return nil, false
		}
		if actual != expected {
			return actual, false
		}
		return actual, true
	}

	for _, respExample := range respExamples {
		testParserSingleMessages(t, respExample.message, compare, respExample.expected, nil)
	}
}

func TestParserBulkStringrMessages(t *testing.T) {
	// RESP protocol spec examples.
	respExamples := []struct {
		message       string
		expected      []byte
		expectedError error
	}{
		{
			message:       "$5\r\nhello\r\n",
			expected:      []byte("hello"),
			expectedError: nil,
		},
		{
			message:       "$6\r\nbina\ry\r\n",
			expected:      []byte("bina\ry"),
			expectedError: nil,
		},
		{
			message:       "$200\r\nlength err\r\n",
			expected:      []byte("noway"),
			expectedError: fmt.Errorf(errorInvalidBulkStringLength, 12, 200),
		},
		{
			message:       "$15\r\ndelimiter errorXY",
			expected:      []byte("noway"),
			expectedError: fmt.Errorf(errorInvalidBulkStringDelim, "XY"),
		},
		{
			message:       "$0\r\n\r\n",
			expected:      []byte(""),
			expectedError: nil,
		},
		{
			message:       "$-1\r\n",
			expected:      nil,
			expectedError: nil,
		},
	}

	compare := func(msg *Message, exp any) (any, bool) {
		expected, ok := exp.([]byte)
		if !ok {
			return nil, false
		}
		if !msg.IsBulk() {
			return nil, false
		}
		actual, err := msg.Bytes()
		if err != nil {
			return nil, false
		}
		if !bytes.Equal(actual, expected) {
			return actual, false
		}
		return actual, true
	}

	for _, respExample := range respExamples {
		testParserSingleMessages(t, respExample.message, compare, respExample.expected, respExample.expectedError)
	}
}

func TestParserArrayMessages(t *testing.T) {
	// RESP protocol spec examples.
	respExamples := []struct {
		message  string
		expected [][]byte
	}{
		{
			message:  "*0\r\n",
			expected: [][]byte{},
		},
		{
			message:  "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			expected: [][]byte{[]byte("hello"), []byte("world")},
		},
		{
			message:  "*3\r\n:1\r\n:2\r\n:3\r\n",
			expected: [][]byte{[]byte("1"), []byte("2"), []byte("3")},
		},
		{
			message:  "*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n",
			expected: [][]byte{[]byte("1"), []byte("2"), []byte("3"), []byte("4"), []byte("hello")},
		},
		// Null Arrays
		{
			message:  "*-1\r\n",
			expected: [][]byte{},
		},
		// Nested arrays
		{
			message:  "*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n",
			expected: [][]byte{[]byte("1"), []byte("2"), []byte("3"), []byte("Hello"), []byte("World")},
		},
		// Null elements in Arrays
		{
			message:  "*3\r\n$5\r\nhello\r\n$-1\r\n$5\r\nworld\r\n",
			expected: [][]byte{[]byte("hello"), nil, []byte("world")},
		},
	}

	compare := func(msg *Message, exp any) (any, bool) {
		expected, ok := exp.([]byte)
		if !ok {
			return nil, false
		}
		actual, err := msg.Bytes()
		if err != nil {
			return nil, false
		}
		if !bytes.Equal(actual, expected) {
			return actual, false
		}
		return actual, true
	}

	for _, respExample := range respExamples {
		msgStr := respExample.message
		parser := NewParserWithBytes([]byte(msgStr))

		msg, err := parser.Next()
		if err != nil {
			t.Errorf("%s %s", msgStr, err)
			continue
		}

		msgIndex := 0
		for msg != nil {
			array, err := msg.Array()
			if err != nil {
				t.Errorf("%s %s", msgStr, err)
				continue
			}

			arrayMsg, err := array.Next()
			if err != nil {
				t.Errorf("%s %s", msgStr, err)
				continue
			}

			if arrayMsg == nil {
				break
			}

			// Nested array ?
			if arrayMsg.IsArray() {
				array, err = arrayMsg.Array()
				if err != nil {
					t.Errorf("%s %s", msgStr, err)
					continue
				}
				arrayMsg, err = array.Next()
				if err != nil {
					t.Errorf("%s %s", msgStr, err)
					continue
				}
			}

			for arrayMsg != nil {
				expectedBytes := respExample.expected[msgIndex]
				if actual, ok := compare(arrayMsg, expectedBytes); !ok {
					t.Errorf("%s %s != %s", msgStr, actual, expectedBytes)
					return
				}
				msgIndex++
				arrayMsg, err = array.Next()
				if err != nil {
					t.Errorf("%s %s", msgStr, err)
					continue
				}
			}

			msg, err = parser.Next()
			if err != nil {
				t.Errorf("%s %s", msgStr, err)
				continue
			}
		}
	}
}
