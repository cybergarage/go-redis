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

import "strconv"

// Array represents a array message.
type Array struct {
	parser *Parser
	index  int
	size   int
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
	array := &Array{
		parser: parser,
		index:  0,
		size:   num,
	}
	return array, nil
}

// Size returns the array size.
func (array *Array) Size() int {
	return array.size
}

// Next returns a next message.
func (array *Array) Next() (*Message, error) {
	if array.size <= array.index {
		return nil, nil
	}
	msg, err := array.parser.Next()
	if err != nil {
		return nil, err
	}
	array.index++
	return msg, nil
}
