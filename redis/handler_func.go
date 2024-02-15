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
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cybergarage/go-redis/redis/glob"
	"github.com/cybergarage/go-redis/redis/proto"
)

// General argument fuctions

func nextIntegerArgument(cmd string, name string, args Arguments) (int, error) {
	val, err := args.NextInteger()
	if err != nil {
		return 0, newMissingArgumentError(cmd, name, err)
	}
	return val, nil
}

func nextStringArgument(cmd string, name string, args Arguments) (string, error) {
	str, err := args.NextString()
	if err != nil {
		return "", newMissingArgumentError(cmd, name, err)
	}
	return str, nil
}

func nextFloatArgument(cmd string, name string, args Arguments) (float64, error) {
	str, err := args.NextString()
	if err != nil {
		return 0, newMissingArgumentError(cmd, name, err)
	}
	score, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, newMissingArgumentError(cmd, name, err)
	}
	return score, nil
}

func nextStringArrayArguments(cmd string, name string, args Arguments) ([]string, error) {
	var str string
	var err error
	strs := []string{}
	str, err = args.NextString()
	for err == nil {
		strs = append(strs, str)
		str, err = args.NextString()
	}
	if !errors.Is(err, proto.ErrEOM) {
		return nil, newMissingArgumentError(cmd, name, err)
	}
	return strs, nil
}

func nextStringMapArguments(cmd string, args Arguments) (map[string]string, error) {
	var key, val string
	var err error
	dir := map[string]string{}
	key, err = args.NextString()
	for err == nil {
		val, err = args.NextString()
		if err != nil {
			newMissingArgumentError(cmd, key, err)
		}
		dir[key] = val
		key, err = args.NextString()
	}
	if !errors.Is(err, proto.ErrEOM) {
		return nil, err
	}
	return dir, nil
}

// Key argument fuctions

func nextKeyArgument(cmd string, args Arguments) (string, error) {
	return nextStringArgument(cmd, "key", args)
}

func nextKeysArguments(cmd string, args Arguments) ([]string, error) {
	return nextStringArrayArguments(cmd, "keys", args)
}

// String argument functions

func nextSetArguments(cmd string, args Arguments) (string, string, error) {
	key, err := args.NextString()
	if err != nil {
		return "", "", newMissingArgumentError(cmd, "key", err)
	}
	val, err := args.NextString()
	if err != nil {
		return "", "", newMissingArgumentError(cmd, "value", err)
	}
	return key, val, err
}

func nextMGetArguments(cmd string, args Arguments) ([]string, error) {
	return nextStringArrayArguments(cmd, "keys", args)
}

func nextMSetArguments(cmd string, args Arguments) (map[string]string, error) {
	return nextStringMapArguments(cmd, args)
}

// Hash argument fuctions

func nextHashArgument(cmd string, args Arguments) (string, error) {
	return nextStringArgument(cmd, "hash", args)
}

// List argument fuctions

func nextPushArguments(cmd string, args Arguments) (string, []string, error) {
	key, err := nextKeyArgument(cmd, args)
	if err != nil {
		return "", nil, err
	}
	elems, err := nextStringArrayArguments(cmd, "elements", args)
	return key, elems, err
}

func nextPopArguments(cmd string, args Arguments) (string, int, error) {
	key, err := nextKeyArgument(cmd, args)
	if err != nil {
		return "", 0, err
	}
	cnt, err := nextIntegerArgument(cmd, "count", args)
	if err != nil {
		if !errors.Is(err, proto.ErrEOM) {
			return "", 0, err
		}
		cnt = 1
	}
	return key, cnt, nil
}

// ZSet fuctions

func nextScoreArgument(cmd string, name string, args Arguments) (float64, error) {
	return nextFloatArgument(cmd, name, args)
}

func nextRangeIndexArgument(cmd string, name string, args Arguments) (int, error) {
	return nextIntegerArgument(cmd, name, args)
}

func nextRangeScoreIndexArgument(cmd string, name string, args Arguments) (float64, bool, error) {
	str, err := args.NextString()
	if err != nil || len(str) == 0 {
		return 0, false, newMissingArgumentError(cmd, name, err)
	}
	offset := 0
	exclusive := false
	if str[0] == '(' {
		offset = 1
		exclusive = true
	}
	rng, err := strconv.ParseFloat(str[offset:], 64)
	if err != nil {
		return 0, false, newMissingArgumentError(cmd, name, err)
	}
	return rng, exclusive, nil
}

func nextRangeOptionArguments(cmd string, args Arguments) (ZRangeOption, error) {
	opt := ZRangeOption{
		BYSCORE:      false,
		BYLEX:        false,
		REV:          false,
		WITHSCORES:   false,
		MINEXCLUSIVE: false,
		MAXEXCLUSIVE: false,
		Offset:       0,
		Count:        -1,
	}

	param, err := args.NextString()
	for err == nil {
		switch strings.ToUpper(param) {
		case "BYSCORE":
			opt.BYSCORE = true
		case "BYLEX":
			opt.BYLEX = true
		case "REV":
			opt.REV = true
		case "LIMIT":
			opt.Offset, err = nextIntegerArgument(cmd, "offset", args)
			if err != nil {
				return opt, newMissingArgumentError(cmd, "offset", err)
			}
			opt.Count, err = nextIntegerArgument(cmd, "count", args)
			if err != nil {
				return opt, newMissingArgumentError(cmd, "count", err)
			}
		case "WITHSCORES":
			opt.WITHSCORES = true
		}
		param, err = args.NextString()
	}
	if !errors.Is(err, proto.ErrEOM) {
		return opt, newMissingArgumentError(cmd, "", err)
	}

	return opt, nil
}

// Expire argument fuctions

func nextExpireArgument(cmd string, ttl time.Time, args Arguments) (ExpireOption, error) {
	opt := ExpireOption{
		Time: ttl,
		NX:   false,
		XX:   false,
		GT:   false,
		LT:   false,
	}
	var err error
	arg, err := args.NextString()
	if err == nil {
		switch strings.ToUpper(arg) {
		case "NX":
			opt.NX = true
		case "XX":
			opt.XX = true
		case "GT":
			opt.GT = true
		case "LT":
			opt.LT = true
		default:
			return opt, newUnkownArgumentError(cmd, arg)
		}
	}
	if !errors.Is(err, proto.ErrEOM) {
		return opt, err
	}
	return opt, nil
}

// Scan argument fuctions

func nextScanArgument(cmd string, args Arguments) (ScanOption, error) {
	opt := ScanOption{
		MatchPattern: glob.MustCompile(DefaultScanPattern),
		Count:        DefaultScanCount,
		Type:         KeyScan,
	}
	var err error
	param, err := args.NextString()
	for err == nil {
		switch strings.ToUpper(param) {
		case "MATCH":
			var pattern string
			pattern, err = nextStringArgument(cmd, "pattern", args)
			if err != nil {
				return opt, newMissingArgumentError(cmd, "pattern", err)
			}
			opt.MatchPattern, err = regexp.Compile(pattern)
			if err != nil {
				return opt, newMissingArgumentError(cmd, "pattern", err)
			}
		case "COUNT":
			opt.Count, err = nextIntegerArgument(cmd, "count", args)
			if err != nil {
				return opt, newMissingArgumentError(cmd, "count", err)
			}
		case "Type":
			var scanType string
			scanType, err = nextStringArgument(cmd, "type", args)
			if err != nil {
				return opt, newMissingArgumentError(cmd, "type", err)
			}
			opt.Type, err = newScanTypeFromString(scanType)
			if err != nil {
				return opt, newMissingArgumentError(cmd, "type", err)
			}
		}
		param, err = args.NextString()
	}
	if !errors.Is(err, proto.ErrEOM) {
		return opt, newMissingArgumentError(cmd, "", err)
	}
	return opt, nil
}
