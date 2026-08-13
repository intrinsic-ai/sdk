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

package promptutil

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestPromptYesNo(t *testing.T) {
	tests := []struct {
		name                 string
		prompt               string
		defaultBehavior      DefaultBehavior
		invalidInputBehavior InvalidInputBehavior
		input                string
		wantResult           bool
		wantErr              bool
	}{
		{
			name:                 "default yes when empty input and defaultYes true",
			prompt:               "Continue?",
			defaultBehavior:      DefaultYes,
			invalidInputBehavior: ReemitPromptOnInvalidInput,
			input:                "\n",
			wantResult:           true,
		},
		{
			name:                 "default no when empty input and defaultYes false",
			prompt:               "Continue?",
			defaultBehavior:      DefaultNo,
			invalidInputBehavior: ReemitPromptOnInvalidInput,
			input:                "\n",
			wantResult:           false,
		},
		{
			name:                 "explicit yes returns true regardless of default",
			prompt:               "Continue?",
			defaultBehavior:      DefaultNo,
			invalidInputBehavior: ReemitPromptOnInvalidInput,
			input:                "y\n",
			wantResult:           true,
		},
		{
			name:                 "capital Y returns true",
			prompt:               "Continue?",
			defaultBehavior:      DefaultNo,
			invalidInputBehavior: ReemitPromptOnInvalidInput,
			input:                "Y\n",
			wantResult:           true,
		},
		{
			name:                 "explicit no returns false regardless of default",
			prompt:               "Continue?",
			defaultBehavior:      DefaultYes,
			invalidInputBehavior: ReemitPromptOnInvalidInput,
			input:                "n\n",
			wantResult:           false,
		},
		{
			name:                 "random character triggers reprompt, then explicit yes returns true",
			prompt:               "Continue?",
			defaultBehavior:      DefaultNo,
			invalidInputBehavior: ReemitPromptOnInvalidInput,
			input:                "x\ny\n",
			wantResult:           true,
		},
		{
			name:                 "random character triggers error when configured",
			prompt:               "Continue?",
			defaultBehavior:      DefaultNo,
			invalidInputBehavior: EmitErrorOnInvalidInput,
			input:                "x\n",
			wantErr:              true,
		},
		{
			name:                 "no default triggers error on empty input",
			prompt:               "Continue?",
			defaultBehavior:      NoDefault,
			invalidInputBehavior: EmitErrorOnInvalidInput,
			input:                "\n",
			wantErr:              true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := &cobra.Command{}

			var inBuf bytes.Buffer
			inBuf.WriteString(tc.input)
			cmd.SetIn(&inBuf)

			var outBuf bytes.Buffer
			cmd.SetOut(&outBuf)

			got, err := PromptYesNo(cmd, tc.prompt, tc.defaultBehavior, tc.invalidInputBehavior)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "invalid input")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantResult, got)
				assert.Contains(t, outBuf.String(), tc.prompt)
			}
		})
	}
}
