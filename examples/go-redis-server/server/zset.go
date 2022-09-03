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
	"github.com/cybergarage/go-redis/redis"
)

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

func NewZSetMember(score float64, data string) *ZSetMember {
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

	offset := opt.Offset
	if offset < 0 {
		offset = 0
	}
	count := opt.Count
	if count < 0 {
		count = len(mems)
	}

	return mems[offset:count]
}

func (zset *ZSet) RangeByScore(min float64, max float64, opt ZRangeOption) []*ZSetMember {
	mems := []*ZSetMember{}
	for _, mem := range zset.members {
		if (mem.Score < min && !opt.MINEXCLUSIVE) || (mem.Score <= min && opt.MINEXCLUSIVE) {
			continue
		}
		if (max < mem.Score && !opt.MAXEXCLUSIVE) || (max <= mem.Score && opt.MAXEXCLUSIVE) {
			continue
		}
		mems = append(mems, mem)
	}

	offset := opt.Offset
	if offset < 0 {
		offset = 0
	}
	count := opt.Count
	if count < 0 {
		count = len(mems)
	}

	return mems[offset:count]
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

func (zset *ZSet) Score(member string) (float64, bool) {
	for _, m := range zset.members {
		if m.Data == member {
			return m.Score, true
		}
	}
	return 0, false
}

func (zset *ZSet) IncBy(inc float64, member string) float64 {
	var m *ZSetMember
	var n int
	for n, m = range zset.members {
		if m.Data == member {
			zset.members = append(zset.members[:n], zset.members[n+1:]...)
			m.Score += inc
			break
		}
	}
	if m == nil {
		m = &ZSetMember{
			Score: inc,
			Data:  member,
		}
	}
	opt := ZAddOption{
		XX:   false,
		NX:   false,
		LT:   false,
		GT:   false,
		CH:   false,
		INCR: false,
	}
	zset.Add([]*ZSetMember{m}, opt)
	return m.Score
}

////////////////////////////////////////////////////////////
// ZSet command handler
////////////////////////////////////////////////////////////

func (server *Server) ZAdd(ctx *redis.DBContext, key string, members []*redis.ZSetMember, opt redis.ZAddOption) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}
	_, zset, err := db.GetZSetRecord(key)
	if err != nil {
		return nil, err
	}
	return redis.NewIntegerMessage(zset.Add(members, opt)), nil
}

func (server *Server) ZRange(ctx *redis.DBContext, key string, start int, stop int, opt redis.ZRangeOption) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}
	_, zset, err := db.GetZSetRecord(key)
	if err != nil {
		return nil, err
	}
	mems := zset.Range(start, stop, opt)
	arrayMsg := redis.NewArrayMessage()
	array, _ := arrayMsg.Array()
	for _, mem := range mems {
		array.Append(redis.NewBulkMessage(mem.Data))
		if opt.WITHSCORES {
			array.Append(redis.NewFloatMessage(mem.Score))
		}
	}
	return arrayMsg, nil
}

func (server *Server) ZRangeByScore(ctx *redis.DBContext, key string, start float64, stop float64, opt redis.ZRangeOption) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}
	_, zset, err := db.GetZSetRecord(key)
	if err != nil {
		return nil, err
	}
	mems := zset.RangeByScore(start, stop, opt)
	arrayMsg := redis.NewArrayMessage()
	array, _ := arrayMsg.Array()
	for _, mem := range mems {
		array.Append(redis.NewBulkMessage(mem.Data))
		if opt.WITHSCORES {
			array.Append(redis.NewFloatMessage(mem.Score))
		}
	}
	return arrayMsg, nil
}

func (server *Server) ZRem(ctx *redis.DBContext, key string, members []string) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}
	_, zset, err := db.GetZSetRecord(key)
	if err != nil {
		return nil, err
	}
	return redis.NewIntegerMessage(zset.Rem(members)), nil
}

func (server *Server) ZScore(ctx *redis.DBContext, key string, member string) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}
	_, zset, err := db.GetZSetRecord(key)
	if err != nil {
		return redis.NewNilMessage(), nil
	}
	score, ok := zset.Score(member)
	if !ok {
		return redis.NewNilMessage(), nil
	}
	return redis.NewFloatMessage(score), nil
}

func (server *Server) ZIncBy(ctx *redis.DBContext, key string, inc float64, member string) (*redis.Message, error) {
	db, err := server.GetDatabase(ctx.ID())
	if err != nil {
		return nil, err
	}
	_, zset, err := db.GetZSetRecord(key)
	if err != nil {
		return nil, err
	}
	return redis.NewFloatMessage(zset.IncBy(inc, member)), nil
}
