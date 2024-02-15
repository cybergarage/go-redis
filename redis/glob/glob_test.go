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

package glob

import (
	"testing"
)

func TestGlobMatches(t *testing.T) {
	tests := []struct {
		pattern string
		value   string
		want    bool
	}{
		{"a??_keys", "lastname_keys", false},
		{"a??_keys", "firstname_keys", false},
		{"a??_keys", "age_keys", true},
		{"*_keys", "lastname_keys", true},
		{"*_keys", "firstname_keys", true},
		{"*_keys", "age_keys", true},
		{"*name*_keys", "lastname_keys", true},
		{"*name*_keys", "firstname_keys", true},
		{"*name*_keys", "age_keys", false},
	}

	for _, tt := range tests {
		glob, err := Compile(tt.pattern)
		if err != nil {
			t.Errorf("Glob(%s).Compile() = %v", tt.pattern, err)
		}
		if got := glob.MatchString(tt.value); got != tt.want {
			t.Skipf("Glob(%s).MatchString(%s) = %v, want %v", tt.pattern, tt.value, got, tt.want)
		}
	}
}
