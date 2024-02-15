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

package redis

import (
	"time"

	"github.com/cybergarage/go-redis/redis/glob"
)

type ExpireOption struct {
	Time time.Time
	NX   bool
	XX   bool
	GT   bool
	LT   bool
}

type SetOption struct {
	EX      time.Duration
	PX      time.Duration
	EXAT    time.Time
	PXAT    time.Time
	NX      bool
	XX      bool
	KEEPTTL bool
	GET     bool
}

type HSetOption struct {
	NX bool
}

type RenameOption struct {
	NX bool
}

type PushOption struct {
	X bool
}

type ZAddOption struct {
	XX   bool
	NX   bool
	LT   bool
	GT   bool
	CH   bool
	INCR bool
}

type ZRangeOption struct {
	BYSCORE      bool
	BYLEX        bool
	REV          bool
	WITHSCORES   bool
	MINEXCLUSIVE bool
	MAXEXCLUSIVE bool
	Offset       int
	Count        int
}

type ScanType int

const (
	KeyScan ScanType = iota
	SetScan
	HashScan
	SortedSetScan
)

type ScanOption struct {
	MatchPattern *glob.Glob
	Count        int
	Type         ScanType
}

func newScanTypeFromString(str string) (ScanType, error) {
	if len(str) == 0 {
		return KeyScan, nil
	}
	switch str {
	case "SSCAN":
		return SetScan, nil
	case "HSCAN":
		return HashScan, nil
	case "ZSCAN":
		return SortedSetScan, nil
	}
	return 0, NewErrNotSupported(str)
}

func newDefaultSetOption() SetOption {
	now := time.Now()
	return SetOption{
		NX:      false,
		XX:      false,
		EX:      0,
		PX:      0,
		EXAT:    now,
		PXAT:    now,
		KEEPTTL: false,
		GET:     false,
	}
}
