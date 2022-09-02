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

package server

import (
	"fmt"
	"reflect"
	"testing"
)

func TestZSet(t *testing.T) {
	testCases := []struct {
		score    string
		data     string
		expected []string
	}{
		{"6", "six", []string{"six"}},
		{"1", "one", []string{"one", "six"}},
		{"2", "two", []string{"one", "two", "six"}},
		{"3", "three", []string{"one", "two", "three", "six"}},
		{"5", "five", []string{"one", "two", "three", "five", "six"}},
		{"4", "four", []string{"one", "two", "three", "four", "five", "six"}},
	}

	zset := NewZSet()
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%s(%s)", testCase.data, testCase.score), func(t *testing.T) {
			m := &ZSetMember{
				Score: testCase.score,
				Data:  testCase.data,
			}
			zset.Add(m)
			mems := zset.Range(0, -1)
			if !reflect.DeepEqual(mems, testCase.expected) {
				t.Errorf("%s != %s", mems, testCase.expected)
				return
			}
		})
	}
}
