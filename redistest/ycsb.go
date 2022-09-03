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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	ycsbTestRowCount = 10
)

const (
	ycsbPathEnv         = "YCSB_ROOT"
	ycsbWorkloadEnv     = "YCSB_WORKLOAD"
	ycsbDefaultWorkload = "workloada"
)

func ExecYCSBWorkload(t *testing.T, workload string) {
	t.Helper()

	ycsbPath, ok := os.LookupEnv(ycsbPathEnv)
	if !ok {
		t.Skipf("%s is not specified", ycsbPathEnv)
	}

	t.Logf("%s = %s", ycsbPathEnv, ycsbPath)

	ycsbCmd := filepath.Join(ycsbPath, "bin/ycsb.sh")
	_, err := os.Stat(ycsbCmd)
	if err != nil {
		t.Error(err)
		return
	}

	workloadDir := filepath.Join(ycsbPath, "workloads")
	workloadFile := filepath.Join(workloadDir, workload)

	ycsbArgs := []string{
		ycsbCmd,
		"",
		"redis",
		"-s",
		"-P",
		workloadFile,
		"-p",
		"redis.host=localhost",
		"-p",
		"redis.port=6379",
	}

	ycsbWorkloadCmds := []string{
		"load",
		"run",
	}

	for _, ycsbWorkloadCmd := range ycsbWorkloadCmds {
		ycsbArgs[1] = ycsbWorkloadCmd
		cmdStr := strings.Join(ycsbArgs, " ")
		t.Run(ycsbWorkloadCmd, func(t *testing.T) {
			t.Logf("%v", cmdStr)
			out, err := exec.Command(ycsbCmd, ycsbArgs[1:]...).CombinedOutput()
			if err != nil {
				return
			}
			outStr := string(out)
			if strings.Contains(outStr, "FAILED") {
				t.Errorf("%s", outStr)
				return
			}
			t.Logf("%s", outStr)
		})
	}
}

func ExecYCSB(t *testing.T) {
	t.Helper()
	workload, ok := os.LookupEnv(ycsbWorkloadEnv)
	t.Logf("%s = %s", ycsbWorkloadEnv, workload)
	if !ok {
		workload = ycsbDefaultWorkload
	}
	ExecYCSBWorkload(t, workload)
}
