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

import "github.com/cybergarage/go-redis/redis"

////////////////////////////////////////////////////////////
// ZSet
////////////////////////////////////////////////////////////

type ZSet struct {
	members []*ZSetMember
}

type ZSetMember = redis.ZSetMember
type ZRangeOption = redis.ZRangeOption
type ZAddOption = redis.ZAddOption

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

func (zset *ZSet) Add(nms []*ZSetMember, opt ZAddOption) int {
	addedMemberCount := 0
	for _, nm := range nms {
		isAdded := false
		for n, tm := range zset.members {
			if nm.Score < tm.Score {
				zset.members = append(zset.members[:n+1], zset.members[n:]...)
				zset.members[n] = nm
				isAdded = true
				addedMemberCount++
				break
			}
		}
		if !isAdded {
			zset.members = append(zset.members, nm)
			addedMemberCount++
		}
	}
	return addedMemberCount
}

func (zset *ZSet) Range(start int, stop int, opt ZRangeOption) []*ZSetMember {
	if start < 0 {
		start = len(zset.members) + start
	}
	if stop < 0 {
		stop = len(zset.members) + stop
	}
	mems := []*ZSetMember{}
	for n := start; n <= stop; n++ {
		if (n < 0) || ((len(zset.members) - 1) < n) {
			continue
		}
		mems = append(mems, zset.members[n])
	}
	return mems
}

func (zset *ZSet) Rem(members []string) int {
	removedMemberCount := 0
	for _, rm := range members {
		for n, m := range zset.members {
			if m.Data == rm {
				zset.members = append(zset.members[:n], zset.members[n+1:]...)
				removedMemberCount++
				break
			}
		}
	}
	return removedMemberCount
}
