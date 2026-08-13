// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package agents

import (
	"testing"
)

func TestCheckEnv(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want bool
	}{
		// Positive tests
		{
			name: "empty env is okay",
			env:  map[string]string{},
		},
		{
			name: "unrelated env is okay",
			env:  map[string]string{"USER": "joe"},
		},
		{
			name: "env with unknown value found",
			env:  map[string]string{"INVOKER_INFO_NAME": "hal_9000"},
		},
		// Negative tests
		{
			name: "env with any value found",
			env:  map[string]string{"ANTIGRAVITY_AGENT": "1"},
			want: true,
		},
		{
			name: "env with specific value found",
			env:  map[string]string{"INVOKER_INFO_NAME": "gemini_cli"},
			want: true,
		},
	}
	oldGetenv := getenv
	t.Cleanup(func() { getenv = oldGetenv })
	for _, tc := range tests {
		getenv = func(k string) string {
			if v, ok := tc.env[k]; ok {
				return v
			} else {
				return ""
			}
		}
		t.Run(tc.name, func(t *testing.T) {
			if got := checkEnvForAgent(); got != tc.want {
				t.Errorf("checkEnvForAgent() returned %t, wanted %t", got, tc.want)
			}
		})
	}

}

func TestAgentName(t *testing.T) {
	oldGetenv := getenv
	t.Cleanup(func() { getenv = oldGetenv })
	getenv = func(k string) string {
		if k == "INVOKER_INFO_NAME" {
			return "antigravity"
		}
		return ""
	}
	if got := AgentName(); got != "antigravity" {
		t.Errorf("AgentName() = %q, want %q", got, "antigravity")
	}
}
