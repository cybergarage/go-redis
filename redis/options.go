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

import "time"

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

type GetOption struct {
}

type MSetOption struct {
	NX bool
}

type MGetOption struct {
}

type HSetOption struct {
}

type HGetOption struct {
}

type HMSetOption struct {
}

type HMGetOption struct {
}
