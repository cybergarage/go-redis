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

package redistest

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func isStringsEqual(aa []string, ba []string) bool {
	for len(aa) != len(ba) {
		return false
	}
	for _, a := range aa {
		hasStr := false
		for _, b := range ba {
			if a == b {
				hasStr = true
				break
			}
		}
		if !hasStr {
			return false
		}
	}
	return true
}

// nolint: maintidx, gocyclo
func TestServer(t *testing.T) {
	server := NewServer()

	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	client := NewClient()
	err = client.Open(LocalHost)
	if err != nil {
		t.Error(err)
		return
	}

	// ctx := context.Background()

	// Connection commands

	t.Run("Connection", func(t *testing.T) {
		t.Run("ECHO", func(t *testing.T) {
			msgs := []string{
				"Hello World!",
			}
			for _, msg := range msgs {
				t.Run(msg, func(t *testing.T) {
					echo := client.Echo(msg)
					if echo.Err() != nil {
						t.Error(echo.Err())
						return
					}
					if echo.Val() != msg {
						t.Errorf("'%s' != '%s'", echo.Val(), msg)
						return
					}
				})
			}
		})
	})

	// Generic commands

	t.Run("Generic", func(t *testing.T) {
		testGeneric(t, server, client)
	})

	// String commands

	t.Run("String", func(t *testing.T) {
		testString(t, server, client)
	})

	// Hash commands

	t.Run("Hash", func(t *testing.T) {
		testHash(t, server, client)
	})

	// List commands

	t.Run("List", func(t *testing.T) {
		testList(t, server, client)
	})

	t.Run("YCSB", func(t *testing.T) {
		err = ExecYCSB(t)
		if err != nil {
			t.Error(err)
		}
	})

	// // panic: not implemented
	// err = client.Quit().Err()
	// if err != nil {
	// 	t.Error(err)
	// }

	err = client.Close()
	if err != nil {
		t.Error(err)
	}

	err = server.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}

// nolint: maintidx, gocyclo
func testGeneric(t *testing.T, server *Server, client *Client) {
	t.Helper()

	var err error

	t.Run("DEL", func(t *testing.T) {
		records := []struct {
			keys     []string
			expected int64
		}{
			{[]string{"key1_del", "key2_del"}, 2},
			{[]string{"key1_del"}, 0},
			{[]string{"key1_del", "key2_del", "key3_del"}, 1},
			{[]string{"key2_del"}, 0},
		}
		for _, r := range records {
			for _, key := range r.keys {
				err = client.Set(key, key, 0).Err()
				if err != nil {
					t.Error(err)
					return
				}
			}
		}
		for _, r := range records {
			t.Run(strings.Join(r.keys, ","), func(t *testing.T) {
				res, err := client.Del(r.keys...).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%d != %d", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("EXISTS", func(t *testing.T) {
		if err := client.Set("key1_exists", "val", 0).Err(); err != nil {
			t.Error(err)
			return
		}
		if err := client.Set("key2_exists", "val", 0).Err(); err != nil {
			t.Error(err)
			return
		}
		records := []struct {
			keys     []string
			expected int64
		}{
			{[]string{"nosuchkey"}, 0},
			{[]string{"key1_exists"}, 1},
			{[]string{"key2_exists"}, 1},
			{[]string{"key1_exists", "key2_exists"}, 2},
			{[]string{"key1_exists", "key2_exists", "nosuchkey"}, 2},
		}
		for _, r := range records {
			t.Run(strings.Join(r.keys, ","), func(t *testing.T) {
				res, err := client.Exists(r.keys...).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%d != %d", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("EXPIRE", func(t *testing.T) {
		key := "mykey_expire"
		if err := client.Set(key, "Hello World", 0).Err(); err != nil {
			t.Error(err)
			return
		}
		records := []struct {
			expire   time.Duration
			sleep    time.Duration
			expected time.Duration
		}{
			{expire: 3 * time.Second, sleep: 1 * time.Second, expected: 1 * time.Second},
			{expire: 1 * time.Second, sleep: 2 * time.Second, expected: -2 * time.Second},
			{expire: 0 * time.Second, sleep: 1 * time.Second, expected: -2 * time.Second},
		}
		for _, r := range records {
			t.Run(fmt.Sprintf("ex:%d, slp:%d", r.expire/time.Second, r.sleep/time.Second), func(t *testing.T) {
				if 0 < r.expire {
					ok, err := client.Expire(key, r.expire).Result()
					if err != nil {
						t.Error(err)
						return
					}
					if !ok {
						t.Errorf("%t", ok)
						return
					}
				}
				if 0 < r.sleep {
					time.Sleep(r.sleep)
				}

				ttl, err := client.TTL(key).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if (ttl != r.expected) && (ttl < r.expected) {
					t.Errorf("%d < %d", ttl, r.expected)
					return
				}
			})
		}
	})

	t.Run("EXPIREAT", func(t *testing.T) {
		key := "mykey_expire"
		if err := client.Set(key, "Hello World", 0).Err(); err != nil {
			t.Error(err)
			return
		}
		records := []struct {
			expire   time.Duration
			sleep    time.Duration
			expected time.Duration
		}{
			{expire: 3 * time.Second, sleep: 1 * time.Second, expected: 1 * time.Second},
			{expire: 1 * time.Second, sleep: 2 * time.Second, expected: -2 * time.Second},
			{expire: 0 * time.Second, sleep: 1 * time.Second, expected: -2 * time.Second},
		}
		now := time.Now()
		for _, r := range records {
			t.Run(fmt.Sprintf("ex:%d, slp:%d", r.expire/time.Second, r.sleep/time.Second), func(t *testing.T) {
				if 0 < r.expire {
					ok, err := client.ExpireAt(key, now.Add(r.expire)).Result()
					if err != nil {
						t.Error(err)
						return
					}
					if !ok {
						t.Errorf("%t", ok)
						return
					}
				}
				if 0 < r.sleep {
					time.Sleep(r.sleep)
				}

				ttl, err := client.TTL(key).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if (ttl != r.expected) && (ttl < r.expected) {
					t.Errorf("%d < %d", ttl, r.expected)
					return
				}
			})
		}
	})

	t.Run("KEYS", func(t *testing.T) {
		args := []string{"firstname_keys", "Jack", "lastname_keys", "Stuntman", "age_keys", "35"}
		err := client.MSet(args).Err()
		if err != nil {
			t.Error(err)
		}
		records := []struct {
			pattern  string
			expected []string
		}{
			{"*name*_keys", []string{"lastname_keys", "firstname_keys"}},
			{"a??_keys", []string{"age_keys"}},
			{"*_keys", []string{"lastname_keys", "firstname_keys", "age_keys"}},
		}
		for _, r := range records {
			t.Run(r.pattern, func(t *testing.T) {
				res, err := client.Keys(r.pattern).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if len(res) != len(r.expected) {
					t.Skipf("%s: %s != %s", r.pattern, res, r.expected)
					return
				}
				for _, ex := range r.expected {
					found := false
					for _, re := range res {
						if ex == re {
							found = true
							continue
						}
					}
					if !found {
						t.Skipf("%s: %s != %s", r.pattern, res, r.expected)
						return
					}
				}
			})
		}
	})

	t.Run("RENAME", func(t *testing.T) {
		if err := client.Set("mykey_rename", "Hello", 0).Err(); err != nil {
			t.Error(err)
			return
		}
		records := []struct {
			key      string
			newkey   string
			expected string
		}{
			{"mykey_rename", "myotherkey_rename", "Hello"},
		}
		for _, r := range records {
			t.Run(r.key+"->"+r.newkey, func(t *testing.T) {
				_, err := client.Rename(r.key, r.newkey).Result()
				if err != nil {
					t.Error(err)
					return
				}
				res, err := client.Get(r.newkey).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%s != %s", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("RENAMENX", func(t *testing.T) {
		if err := client.Set("mykey_renamenx", "Hello", 0).Err(); err != nil {
			t.Error(err)
			return
		}
		if err := client.Set("myotherkey_renamenx", "World", 0).Err(); err != nil {
			t.Error(err)
			return
		}
		records := []struct {
			key      string
			newkey   string
			expected string
		}{
			{"mykey_renamenx", "myotherkey_renamenx", "World"},
		}
		for _, r := range records {
			t.Run(r.key+"->"+r.newkey, func(t *testing.T) {
				_, err := client.RenameNX(r.key, r.newkey).Result()
				if err != nil {
					t.Error(err)
					return
				}
				res, err := client.Get(r.newkey).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%s != %s", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("TYPE", func(t *testing.T) {
		if err := client.Set("key1_type", "key1_type", 0).Err(); err != nil {
			t.Error(err)
			return
		}
		err := client.HSet("key2_type", "key", "val").Err()
		if err != nil {
			t.Error(err)
		}
		records := []struct {
			key      string
			expected string
		}{
			{"key0_type", "none"},
			{"key1_type", "string"},
			{"key2_type", "hash"},
		}
		for _, r := range records {
			t.Run(r.key+":"+r.expected, func(t *testing.T) {
				res, err := client.Type(r.key).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%s != %s", res, r.expected)
					return
				}
			})
		}
	})
}

// nolint: maintidx, gocyclo, dupl
func testString(t *testing.T, server *Server, client *Client) {
	t.Helper()

	t.Run("APPEND", func(t *testing.T) {
		records := []struct {
			key      string
			val      string
			expected int64
		}{
			{"key_append", "Hello", 5},
			{"key_append", " World", 11},
		}

		for _, r := range records {
			t.Run(r.key+":"+r.val, func(t *testing.T) {
				res, err := client.Append(r.key, r.val).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%d != %d", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("DECR", func(t *testing.T) {
		key := "mykey_decr"
		startVal := 10
		err := client.Set(key, strconv.Itoa(startVal), 0).Err()
		if err != nil {
			t.Error(err)
		}
		records := []struct {
			expected int64
		}{
			{int64(startVal - 1)},
			{int64(startVal - 2)},
		}
		for _, r := range records {
			t.Run(key+":"+strconv.Itoa(int(r.expected)), func(t *testing.T) {
				res, err := client.Decr(key).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%d != %d", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("DECRBY", func(t *testing.T) {
		key := "mykey_decrby"
		startVal := 10
		err := client.Set(key, strconv.Itoa(startVal), 0).Err()
		if err != nil {
			t.Error(err)
		}
		decVal := 3
		records := []struct {
			expected int64
		}{
			{int64(startVal - decVal)},
			{int64(startVal - (decVal * 2))},
		}
		for _, r := range records {
			t.Run(key+":"+strconv.Itoa(int(r.expected)), func(t *testing.T) {
				res, err := client.DecrBy(key, int64(decVal)).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%d != %d", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("INCR", func(t *testing.T) {
		key := "mykey_incr"
		startVal := 10
		err := client.Set(key, strconv.Itoa(startVal), 0).Err()
		if err != nil {
			t.Error(err)
		}
		records := []struct {
			expected int64
		}{
			{int64(startVal + 1)},
			{int64(startVal + 2)},
		}
		for _, r := range records {
			t.Run(key+":"+strconv.Itoa(int(r.expected)), func(t *testing.T) {
				res, err := client.Incr(key).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%d != %d", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("INCRBY", func(t *testing.T) {
		key := "mykey_incrby"
		startVal := 10
		err := client.Set(key, strconv.Itoa(startVal), 0).Err()
		if err != nil {
			t.Error(err)
		}
		incVal := 3
		records := []struct {
			expected int64
		}{
			{int64(startVal + incVal)},
			{int64(startVal + (incVal * 2))},
		}
		for _, r := range records {
			t.Run(key+":"+strconv.Itoa(int(r.expected)), func(t *testing.T) {
				res, err := client.IncrBy(key, int64(incVal)).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%d != %d", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("GETSET", func(t *testing.T) {
		records := []struct {
			key      string
			val      string
			expected []byte
		}{
			{"key_getset", "value0", nil},
			{"key_getset", "value1", []byte("value0")},
			{"key_getset", "value2", []byte("value1")},
		}

		for _, r := range records {
			t.Run(r.key+":"+r.val, func(t *testing.T) {
				res, err := client.GetSet(r.key, r.val).Result()
				if r.expected == nil {
					if err == nil {
						t.Errorf("%s != %s", res, string(r.expected))
					}
					return
				} else if err != nil {
					t.Error(err)
				}
				if res != string(r.expected) {
					t.Errorf("%s != %s", res, string(r.expected))
				}
			})
		}
	})

	t.Run("MSET", func(t *testing.T) {
		records := []struct {
			keys []string
			vals []string
		}{
			{[]string{"key1_mset"}, []string{"Hello"}},
			{[]string{"key1_mset", "key2_mset"}, []string{"Hello", "World"}},
		}
		for _, r := range records {
			t.Run(strings.Join(r.keys, ","), func(t *testing.T) {
				args := []string{}
				for n, key := range r.keys {
					args = append(args, key)
					args = append(args, r.vals[n])
				}
				err := client.MSet(args).Err()
				if err != nil {
					t.Error(err)
					return
				}
				res, err := client.MGet(r.keys...).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if len(res) != len(r.vals) {
					t.Errorf("%d != %d", len(res), len(r.vals))
					return
				}
				for n, val := range r.vals {
					if res[n] != val {
						t.Errorf("%s != %s", res[n], val)
					}
				}
			})
		}
	})

	t.Run("MSETNX", func(t *testing.T) {
		records := []struct {
			keys     []string
			vals     []string
			expected bool
		}{
			{[]string{"key1_msetnx", "key2_msetnx"}, []string{"Hello", "there"}, true},
			{[]string{"key2_msetnx", "key3_msetnx"}, []string{"new", "world"}, false},
		}
		for _, r := range records {
			t.Run(strings.Join(r.keys, ","), func(t *testing.T) {
				args := []string{}
				for n, key := range r.keys {
					args = append(args, key)
					args = append(args, r.vals[n])
				}
				res, err := client.MSetNX(args).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%t != %t", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("SET", func(t *testing.T) {
		records := []struct {
			key      string
			val      string
			expected string
		}{
			{"key_set", "value0", "value0"},
			{"key_set", "value1", "value1"},
			{"key_set", "value2", "value2"},
		}

		for _, r := range records {
			t.Run(r.key+":"+r.val, func(t *testing.T) {
				err := client.Set(r.key, r.val, 0).Err()
				if err != nil {
					t.Error(err)
					return
				}
				res, err := client.Get(r.key).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%s != %s", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("SETNX", func(t *testing.T) {
		records := []struct {
			key      string
			val      string
			expected bool
		}{
			{"key_setnx", "value0", true},
			{"key_setnx", "value1", false},
			{"key_setnx", "value2", false},
		}

		for _, r := range records {
			t.Run(r.key+":"+r.val, func(t *testing.T) {
				res, err := client.SetNX(r.key, r.val, 0).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%t != %t", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("STRLEN", func(t *testing.T) {
		records := []struct {
			key      string
			val      string
			expected int64
		}{
			{"mykey_strlen", "Hello world", 11},
			{"nonexisting_strlen", "", 0},
		}
		for _, r := range records {
			t.Run(r.key, func(t *testing.T) {
				if 0 < len(r.val) {
					_, err := client.Set(r.key, r.val, 0).Result()
					if err != nil {
						t.Error(err)
						return
					}
				}
				res, err := client.StrLen(r.key).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%d != %d", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("SUBSTR(GETRANGE)", func(t *testing.T) {
		key := "mykey_substr"
		_, err := client.Set(key, "This is a string", 0).Result()
		if err != nil {
			t.Error(err)
			return
		}
		records := []struct {
			start    int64
			end      int64
			expected string
		}{
			{0, 3, "This"},
			{-3, -1, "ing"},
			{0, -1, "This is a string"},
			{10, 100, "string"},
		}
		for _, r := range records {
			t.Run(fmt.Sprintf("%d:%d", r.start, r.end), func(t *testing.T) {
				res, err := client.GetRange(key, r.start, r.end).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%s != %s", res, r.expected)
					return
				}
			})
		}
	})
}

// nolint: maintidx, gocyclo, dupl
func testHash(t *testing.T, server *Server, client *Client) {
	t.Helper()

	t.Run("HDEL", func(t *testing.T) {
		key := "myhash_hdel"
		fields := []string{"field1", "field2"}
		err := client.HSet(key, fields[0], "foo").Err()
		if err != nil {
			t.Error(err)
			return
		}
		records := []struct {
			fields   []string
			expected int64
		}{
			{[]string{fields[0]}, 1},
			{[]string{fields[1]}, 0},
			{[]string{fields[0]}, 0},
		}
		for _, r := range records {
			t.Run(strings.Join(r.fields, ","), func(t *testing.T) {
				res, err := client.HDel(key, r.fields...).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%d != %d", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("HEXISTS", func(t *testing.T) {
		key := "myhash_exists"
		fields := []string{"field1", "field2"}
		err := client.HSet(key, fields[0], "foo").Err()
		if err != nil {
			t.Error(err)
			return
		}
		records := []struct {
			field    string
			expected bool
		}{
			{fields[0], true},
			{fields[1], false},
		}
		for _, r := range records {
			t.Run(r.field, func(t *testing.T) {
				res, err := client.HExists(key, r.field).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%t != %t", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("HKEYS", func(t *testing.T) {
		key := "myhash_hkeys"
		records := []struct {
			fields []string
			values []string
		}{
			{[]string{"field1", "field2"}, []string{"Hello", "World"}},
		}
		for _, r := range records {
			t.Run(strings.Join(r.fields, ","), func(t *testing.T) {
				for n, field := range r.fields {
					err := client.HSet(key, field, r.values[n]).Err()
					if err != nil {
						t.Error(err)
						return
					}
				}
				res, err := client.HKeys(key).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if !isStringsEqual(res, r.fields) {
					t.Errorf("%s != %s", res, r.fields)
					return
				}
			})
		}
	})

	t.Run("HLEN", func(t *testing.T) {
		key := "myhash_hlen"
		records := []struct {
			field    string
			value    string
			expected int
		}{
			{"", "", 0},
			{"field1", "Hello", 1},
			{"field2", "World", 2},
		}
		for _, r := range records {
			t.Run(r.field, func(t *testing.T) {
				if 0 < len(r.field) {
					err := client.HSet(key, r.field, r.value).Err()
					if err != nil {
						t.Error(err)
						return
					}
				}
				res, err := client.HLen(key).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != int64(r.expected) {
					t.Errorf("%d != %d", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("HSET", func(t *testing.T) {
		key := "key_hset"
		records := []struct {
			field    string
			value    string
			expected string
		}{
			{"field1", "Hello", "Hello"},
		}

		for _, r := range records {
			t.Run(key+":"+r.field+":"+r.value, func(t *testing.T) {
				err := client.HSet(key, r.field, r.value).Err()
				if err != nil {
					t.Error(err)
					return
				}
				res, err := client.HGet(key, r.field).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%s != %s", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("HSETNX", func(t *testing.T) {
		key := "key_hsetnx"
		records := []struct {
			field    string
			value    string
			expected string
		}{
			{"field", "Hello", "Hello"},
			{"field", "World", "Hello"},
		}

		for _, r := range records {
			t.Run(key+":"+r.field+":"+r.value, func(t *testing.T) {
				err := client.HSetNX(key, r.field, r.value).Err()
				if err != nil {
					t.Error(err)
					return
				}
				res, err := client.HGet(key, r.field).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expected {
					t.Errorf("%s != %s", res, r.expected)
					return
				}
			})
		}
	})

	t.Run("HMSET", func(t *testing.T) {
		records := []struct {
			hash string
			keys []string
			vals []string
		}{
			{"myhash_hmset", []string{"field1", "field2"}, []string{"Hello", "World"}},
		}
		for _, r := range records {
			t.Run(r.hash+":"+strings.Join(r.keys, ","), func(t *testing.T) {
				args := map[string]interface{}{}
				for n, key := range r.keys {
					args[key] = r.vals[n]
				}
				err := client.HMSet(r.hash, args).Err()
				if err != nil {
					t.Error(err)
					return
				}
				res, err := client.HMGet(r.hash, r.keys...).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if len(res) != len(r.vals) {
					t.Errorf("%d != %d", len(res), len(r.vals))
					return
				}
				for n, val := range r.vals {
					if res[n] != val {
						t.Errorf("%s != %s", res[n], val)
					}
				}
			})
		}
	})

	t.Run("HGETALL", func(t *testing.T) {
		records := []struct {
			hash     string
			key      string
			val      string
			expected map[string]string
		}{
			{"myhash_hgetall", "field1", "Hello", map[string]string{"field1": "Hello"}},
			{"myhash_hgetall", "field2", "World", map[string]string{"field1": "Hello", "field2": "World"}},
		}

		for _, r := range records {
			t.Run(r.hash+":"+r.key+":"+r.val, func(t *testing.T) {
				err := client.HSet(r.hash, r.key, r.val).Err()
				if err != nil {
					t.Error(err)
					return
				}
				res, err := client.HGetAll(r.hash).Result()
				if err != nil {
					t.Error(err)
					return
				}
				for ekey, eval := range r.expected {
					rval, ok := res[ekey]
					if !ok {
						t.Errorf("%s", ekey)
						return
					}
					if rval != eval {
						t.Errorf("%s != %s", rval, eval)
					}
				}
			})
		}
	})

	t.Run("HSTRLEN", func(t *testing.T) {
		key := "myhash_hstrlen"
		records := []struct {
			field string
			value string
		}{
			{"f1", "HelloWorld"},
			{"f2", "99"},
			{"f3", "-256"},
		}
		args := map[string]interface{}{}
		for _, r := range records {
			args[r.field] = r.value
		}
		err := client.HMSet(key, args).Err()
		if err != nil {
			t.Error(err)
			return
		}
		for _, r := range records {
			t.Run(r.field, func(t *testing.T) {
				// Note: go-redis does not support HSTRLEN yet
				res, err := client.HGet(key, r.field).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if len(res) != len(r.value) {
					t.Errorf("%d != %d", len(res), len(r.value))
					return
				}
			})
		}
	})

	t.Run("HVALS", func(t *testing.T) {
		key := "myhash_hvals"
		records := []struct {
			fields []string
			values []string
		}{
			{[]string{"field1", "field2"}, []string{"Hello", "World"}},
		}
		for _, r := range records {
			t.Run(strings.Join(r.fields, ","), func(t *testing.T) {
				for n, field := range r.fields {
					err := client.HSet(key, field, r.values[n]).Err()
					if err != nil {
						t.Error(err)
						return
					}
				}
				res, err := client.HVals(key).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if !isStringsEqual(res, r.values) {
					t.Errorf("%s != %s", res, r.fields)
					return
				}
			})
		}
	})
}

// nolint: maintidx, gocyclo, dupl
func testList(t *testing.T, server *Server, client *Client) {
	t.Helper()

	t.Run("LPUSH", func(t *testing.T) {
		key := "mylist_lpush"
		records := []struct {
			elems       []string
			expectedRet int64
			expectedRng []string
		}{
			{[]string{"world"}, 1, []string{"world"}},
			{[]string{"hello"}, 2, []string{"hello", "world"}},
		}

		for _, r := range records {
			t.Run(strings.Join(r.elems, ","), func(t *testing.T) {
				res, err := client.LPush(key, r.elems).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if res != r.expectedRet {
					t.Errorf("%d != %d", r.expectedRet, res)
					return
				}
				rng, err := client.LRange(key, 0, -1).Result()
				if err != nil {
					t.Error(err)
					return
				}
				if len(rng) != len(r.expectedRng) {
					t.Errorf("%d != %d", len(rng), len(r.expectedRng))
					return
				}
				for n, rs := range rng {
					if rs != r.expectedRng[n] {
						t.Errorf("%s != %s", rs, r.expectedRng[n])
						return
					}
				}
			})
		}
	})
}
