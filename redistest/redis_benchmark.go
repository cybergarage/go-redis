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
	"os/exec"
	"strings"
	"testing"
)

const (
	redisBenchmarkCmd = "redis-benchmark"
)

func ExecRedisBenchmark(t *testing.T, param string) {
	t.Helper()

	args := strings.Split(param, " ")
	out, err := exec.Command(redisBenchmarkCmd, args...).CombinedOutput()
	if err != nil {
		return
	}
	outStr := string(out)
	if strings.Contains(outStr, "FAILED") {
		t.Errorf("%s", outStr)
		return
	}
	t.Logf("%s", outStr)
}
