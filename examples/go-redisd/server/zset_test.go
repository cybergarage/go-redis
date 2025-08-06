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

package server

import (
	"fmt"
	"reflect"
	"testing"
)

func TestZSet(t *testing.T) {
	zset := NewZSet()

	addCases := []struct {
		score    float64
		data     string
		expected []string
	}{
		{6, "six", []string{"six"}},
		{1, "one", []string{"one", "six"}},
		{2, "two", []string{"one", "two", "six"}},
		{3, "three", []string{"one", "two", "three", "six"}},
		{5, "five", []string{"one", "two", "three", "five", "six"}},
		{4, "four", []string{"one", "two", "three", "four", "five", "six"}},
	}

	zaopt := ZAddOption{
		XX:   false,
		NX:   false,
		LT:   false,
		GT:   false,
		CH:   false,
		INCR: false,
	}

	zropt := ZRangeOption{
		BYSCORE:      false,
		BYLEX:        false,
		REV:          false,
		WITHSCORES:   false,
		MINEXCLUSIVE: false,
		MAXEXCLUSIVE: false,
		Offset:       0,
		Count:        -1,
	}

	for _, r := range addCases {
		t.Run("Add", func(t *testing.T) {
			t.Run(fmt.Sprintf("%s(%f)", r.data, r.score), func(t *testing.T) {
				m := &ZSetMember{
					Score:  r.score,
					Member: r.data,
				}
				zset.Add([]*ZSetMember{m}, zaopt)
				mems := zset.Range(0, -1, zropt)

				memdata := []string{}
				for _, mem := range mems {
					memdata = append(memdata, mem.Member)
				}

				if !reflect.DeepEqual(memdata, r.expected) {
					t.Errorf("%s != %s", memdata, r.expected)
					return
				}
			})
		})
	}

	remCases := []struct {
		data     string
		expected []string
	}{
		{"six", []string{"one", "two", "three", "four", "five"}},
		{"one", []string{"two", "three", "four", "five"}},
		{"two", []string{"three", "four", "five"}},
		{"three", []string{"four", "five"}},
		{"five", []string{"four"}},
		{"four", []string{}},
	}

	for _, r := range remCases {
		t.Run("Rem", func(t *testing.T) {
			t.Run(r.data, func(t *testing.T) {
				zset.Rem([]string{r.data})
				mems := zset.Range(0, -1, zropt)

				memdata := []string{}
				for _, mem := range mems {
					memdata = append(memdata, mem.Member)
				}

				if !reflect.DeepEqual(memdata, r.expected) {
					t.Errorf("%s != %s", memdata, r.expected)
					return
				}
			})
		})
	}
}
