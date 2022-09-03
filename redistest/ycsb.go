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
	ycsbPathEnv         = "GO_REDIS_YCSB_ROOT_PATH"
	ycsbWorkloadEnv     = "GO_REDIS_YCSB_WORKLOAD"
	ycsbDefaultWorkload = "workloada"
)

func ExecYCSBWorkload(t *testing.T, defaultWorkload string) error {
	t.Helper()
	outputYcsbParams := func(t *testing.T, ycsbEnvs []string, ycsbParams []string) {
		t.Helper()
		for n, ycsbEnv := range ycsbEnvs {
			t.Logf("%s = %s", ycsbEnv, ycsbParams[n])
		}
	}

	ycsbEnvs := []string{
		ycsbPathEnv,
		ycsbWorkloadEnv,
	}

	ycsbParams := []string{
		"",
		defaultWorkload,
	}

	for n, ycsbEnv := range ycsbEnvs {
		if v, ok := os.LookupEnv(ycsbEnv); ok {
			ycsbParams[n] = v
		}
		if len(ycsbParams[n]) == 0 {
			outputYcsbParams(t, ycsbEnvs, ycsbParams)
			t.Skipf("%s is not specified", ycsbEnv)
			return nil
		}
	}

	outputYcsbParams(t, ycsbEnvs, ycsbParams)

	ycsbPath := ycsbParams[0]
	ycsbCmd := filepath.Join(ycsbPath, "bin/ycsb.sh")
	_, err := os.Stat(ycsbCmd)
	if err != nil {
		return err
	}

	workload := ycsbParams[1]
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
		t.Logf("%v", cmdStr)
		out, err := exec.Command(ycsbCmd, ycsbArgs[1:]...).CombinedOutput()
		if err != nil {
			return err
		}
		outStr := string(out)
		if strings.Contains(outStr, "FAILED") {
			t.Errorf("%s", outStr)
			continue
		}
		t.Logf("%s", outStr)
	}

	return nil
}

func ExecYCSB(t *testing.T) error {
	t.Helper()
	return ExecYCSBWorkload(t, ycsbDefaultWorkload)
}
