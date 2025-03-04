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
	"fmt"
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

func nextSetExArguments(cmd string, args Arguments) (string, int, string, error) {
	key, err := args.NextString()
	if err != nil {
		return "", 0, "", newMissingArgumentError(cmd, "key", err)
	}
	seconds, err := args.NextInteger()
	if err != nil {
		return "", 0, "", newMissingArgumentError(cmd, "seconds", err)
	}
	if seconds < 1 {
		return "", 0, "", newInvalidArgumentError(cmd, "seconds", fmt.Errorf(errorShouldBeGreaterThanInt, "argument", 0))
	}
	val, err := args.NextString()
	if err != nil {
		return "", 0, "", newMissingArgumentError(cmd, "value", err)
	}
	return key, seconds, val, err
}

func nextMGetArguments(cmd string, args Arguments) ([]string, error) {
	return nextStringArrayArguments(cmd, "keys", args)
}

func nextMSetArguments(cmd string, args Arguments) (map[string]string, error) {
	return nextStringMapArguments(cmd, args)
}

func nextSetOptionArguments(cmd string, args Arguments) (SetOption, error) {
	opt := newDefaultSetOption()
	for {
		argStr, err := args.NextString()
		if err != nil {
			if errors.Is(err, proto.ErrEOM) {
				break
			} else {
				return opt, err
			}
		}
		argStr = strings.ToUpper(argStr)
		switch argStr {
		case "NX":
			if opt.NX || opt.XX {
				return opt, newInvalidArgumentError(cmd, argStr, fmt.Errorf(errorUseOnlyOnce, "NX|XX"))
			}
			opt.NX = true
		case "XX":
			if opt.NX || opt.XX {
				return opt, newInvalidArgumentError(cmd, argStr, fmt.Errorf(errorUseOnlyOnce, "NX|XX"))
			}
			opt.XX = true
		case "EX", "PX", "EXAT", "PXAT":
			if opt.EX > 0 || opt.PX > 0 || !opt.EXAT.IsZero() || !opt.PXAT.IsZero() {
				return opt, newInvalidArgumentError(cmd, argStr, fmt.Errorf(errorUseOnlyOnce, "EX|PX|EXAT|PXAT"))
			}
			argInt, err := args.NextInteger()
			if err != nil {
				if errors.Is(err, proto.ErrEOM) {
					return opt, newMissingArgumentError(cmd, argStr, err)
				} else {
					return opt, err
				}
			}
			if argInt < 1 {
				return opt, newInvalidArgumentError(cmd, argStr, fmt.Errorf(errorShouldBeGreaterThanInt, "expire", 0))
			}
			switch argStr {
			case "EX":
				opt.EX = time.Duration(argInt) * time.Second
			case "PX":
				opt.PX = time.Duration(argInt) * time.Millisecond
			case "EXAT":
				opt.EXAT = time.Unix(int64(argInt), 0)
			case "PXAT":
				opt.PXAT = time.UnixMilli(int64(argInt))
			}
		case "KEEPTTL":
			if opt.KEEPTTL {
				return opt, newInvalidArgumentError(cmd, argStr, fmt.Errorf(errorUseOnlyOnce, ""))
			}
			opt.KEEPTTL = true
		case "GET":
			if opt.GET {
				return opt, newInvalidArgumentError(cmd, argStr, fmt.Errorf(errorUseOnlyOnce, ""))
			}
			opt.GET = true
		default:
			return opt, newUnkownArgumentError(cmd, argStr)
		}
	}
	return opt, nil
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
		return 0, false, newInvalidArgumentError(cmd, name, err)
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
				return opt, err
			}
			opt.Count, err = nextIntegerArgument(cmd, "count", args)
			if err != nil {
				return opt, err
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
		Type:         DefaultScanType,
	}
	var err error
	param, err := args.NextString()
	for err == nil {
		switch strings.ToUpper(param) {
		case "MATCH":
			var pattern string
			pattern, err = nextStringArgument(cmd, "pattern", args)
			if err != nil {
				return opt, err
			}
			opt.MatchPattern, err = regexp.Compile(pattern)
			if err != nil {
				return opt, newInvalidArgumentError(cmd, "pattern", err)
			}
		case "COUNT":
			opt.Count, err = nextIntegerArgument(cmd, "count", args)
			if err != nil {
				return opt, err
			}
		case "Type":
			var scanType string
			scanType, err = nextStringArgument(cmd, "type", args)
			if err != nil {
				return opt, err
			}
			opt.Type, err = newScanTypeFromString(scanType)
			if err != nil {
				return opt, newInvalidArgumentError(cmd, "type", err)
			}
		}
		param, err = args.NextString()
	}
	if !errors.Is(err, proto.ErrEOM) {
		return opt, newMissingArgumentError(cmd, "", err)
	}
	return opt, nil
}
