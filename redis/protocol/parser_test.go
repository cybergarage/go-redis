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

import "testing"

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
			}
		}
	}
}
