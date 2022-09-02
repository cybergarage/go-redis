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

////////////////////////////////////////////////////////////
// ZSet
////////////////////////////////////////////////////////////

type ZSet struct {
	members []*ZSetMember
}

type ZSetMember struct {
	Score string
	Data  string
}

func NewZSet() *ZSet {
	return &ZSet{
		members: []*ZSetMember{},
	}
}

func NewZSetMember(score string, data string) *ZSetMember {
	return &ZSetMember{
		Score: score,
		Data:  data,
	}
}

func (zset *ZSet) Add(nm *ZSetMember) {
	for n, tm := range zset.members {
		if nm.Score < tm.Score {
			zset.members = append(zset.members[:n+1], zset.members[n:]...)
			zset.members[n] = nm
			return
		}
	}
	zset.members = append(zset.members, nm)
}

func (zset *ZSet) Range(start int, stop int) []string {
	if start < 0 {
		start = len(zset.members) + start
	}
	if stop < 0 {
		stop = len(zset.members) + stop
	}
	mems := []string{}
	for n := start; n <= stop; n++ {
		if (n < 0) || ((len(zset.members) - 1) < n) {
			continue
		}
		mems = append(mems, zset.members[n].Data)
	}
	return mems
}

////////////////////////////////////////////////////////////
// ZSet command handler
////////////////////////////////////////////////////////////
