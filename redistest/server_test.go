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

	// conn := context.Background()

	t.Run("Command", func(t *testing.T) {
		CommandTest(t, client)
	})

	// redis-benchmark
	params := []string{"-t get,set -n 10000"}
	t.Run("redis-benchmark", func(t *testing.T) {
		for _, param := range params {
			t.Run(param, func(t *testing.T) {
				ExecRedisBenchmark(t, param)
			})
		}
	})

	// YCSB

	workloads := []string{"workloada", "workloadb"}
	t.Run("YCSB", func(t *testing.T) {
		for _, workload := range workloads {
			t.Run(workload, func(t *testing.T) {
				ExecYCSBWorkload(t, workload)
			})
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
