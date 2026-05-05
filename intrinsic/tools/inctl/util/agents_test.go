// Copyright 2023 Intrinsic Innovation LLC

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
	defer func() { getenv = oldGetenv }()
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
