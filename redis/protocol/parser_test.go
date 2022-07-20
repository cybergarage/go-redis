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
	"testing"
)

// RESP protocol spec examples.
var respExampleMessages = []string{
	// RESP Simple Strings
	"+OK\r\n",
	// RESP Errors
	"-Error message\r\n",
	"-ERR unknown command 'helloworld'",
	"-WRONGTYPE Operation against a key holding the wrong kind of value",
	// RESP Integers
	":0\r\n",
	":1000\r\n",
	// RESP Bulk Strings
	"$5\r\nhello\r\n",
	"$0\r\n\r\n",
	"$-1\r\n",
	// RESP Arrays
	"*0\r\n",
	"*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
	"*3\r\n:1\r\n:2\r\n:3\r\n",
	"*-1\r\n",
}

var testMessages = respExampleMessages

func TestParser(t *testing.T) {
	for _, testMsg := range testMessages {
		parser := NewParser()
		err := parser.Parse([]byte(testMsg))
		if err != nil {
			t.Errorf("%s %s", testMsg, err)
		}
		msg, err := parser.Next()
		for msg != nil && err != nil {
			msg, err = parser.Next()
			if err != nil {
				t.Errorf("%s %s", testMsg, err)
				continue
			}
		}
	}
}

func testParsergMessages(t *testing.T, msgString string, compare func(*Message, any) (any, bool), expected any) {
	t.Helper()

	parser := NewParser()
	err := parser.Parse([]byte(msgString))
	if err != nil {
		t.Errorf("%s %s", msgString, err)
		return
	}
	msg, err := parser.Next()
	if err != nil {
		t.Errorf("%s %s", msgString, err)
		return
	}
	if actual, ok := compare(msg, expected); !ok {
		t.Errorf("%s != %s", actual, expected)
		return
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
		testParsergMessages(t, respExample.message, compare, respExample.expected)
	}
}
