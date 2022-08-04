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

package redistest

import (
	"strconv"
	"strings"
	"testing"
)

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

	// String commands

	t.Run("String", func(t *testing.T) {
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
					}
					if res != r.expected {
						t.Errorf("%d != %d", res, r.expected)
					}
				})
			}
		})

		t.Run("DECR", func(t *testing.T) {
			key := "mykey_decr"
			startVal := 10
			err = client.Set(key, strconv.Itoa(startVal), 0).Err()
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
					}
					if res != r.expected {
						t.Errorf("%d != %d", res, r.expected)
					}
				})
			}
		})

		t.Run("DECRBY", func(t *testing.T) {
			key := "mykey_decrby"
			startVal := 10
			err = client.Set(key, strconv.Itoa(startVal), 0).Err()
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
					}
					if res != r.expected {
						t.Errorf("%d != %d", res, r.expected)
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
					err = client.Set(r.key, r.val, 0).Err()
					if err != nil {
						t.Error(err)
					}
					res, err := client.Get(r.key).Result()
					if err != nil {
						t.Error(err)
					}
					if res != r.expected {
						t.Errorf("%s != %s", res, r.expected)
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
					}
					if res != r.expected {
						t.Errorf("%t != %t", res, r.expected)
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
	})

	// Hash commands

	t.Run("Hash", func(t *testing.T) {
		t.Run("HSET", func(t *testing.T) {
			records := []struct {
				hash     string
				key      string
				val      string
				expected string
			}{
				{"key_hset", "key1", "Hello", "Hello"},
			}

			for _, r := range records {
				t.Run(r.hash+":"+r.key+":"+r.val, func(t *testing.T) {
					err := client.HSet(r.hash, r.key, r.val).Err()
					if err != nil {
						t.Error(err)
					}
					res, err := client.HGet(r.hash, r.key).Result()
					if err != nil {
						t.Error(err)
					}
					if res != r.expected {
						t.Errorf("%s != %s", res, r.expected)
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
					}
					res, err := client.HGetAll(r.hash).Result()
					if err != nil {
						t.Error(err)
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
