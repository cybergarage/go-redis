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

type ZSet []*ZSetMember

////////////////////////////////////////////////////////////
// ZSet
////////////////////////////////////////////////////////////

type ZSetMember struct {
	Score string
	Data  string
}

// nolint: staticcheck
func (zset ZSet) Add(nm *ZSetMember) {
	for n, tm := range zset {
		if nm.Score < tm.Score {
			zset = append(zset[:n+1], zset[n:]...)
			zset[n] = nm
			return
		}
	}
	zset = append(zset, nm)
}

////////////////////////////////////////////////////////////
// ZSet command handler
////////////////////////////////////////////////////////////
