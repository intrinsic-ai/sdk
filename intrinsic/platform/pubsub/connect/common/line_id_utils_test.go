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

package lineidutils

import (
	"fmt"
	"strings"
	"testing"
)

func TestValidation(t *testing.T) {
	tests := []struct {
		name                string
		lineID              string
		expectError         bool
		expectErrorContains string
	}{
		{
			name:        "Valid line id is accepted",
			lineID:      "intrinsic-vmp-123-456",
			expectError: false,
		},
		{
			name:                "Excessively long id is rejected",
			lineID:              strings.Repeat("x", 300),
			expectError:         true,
			expectErrorContains: fmt.Sprintf("more than %d characters long", maxLineIDLength),
		},
		{
			name:                "Id with invalid characters is rejected",
			lineID:              "BIG-company-node-123",
			expectError:         true,
			expectErrorContains: "contains invalid characters",
		},
		{
			name:                "Empty line id is rejected",
			lineID:              "",
			expectError:         true,
			expectErrorContains: "line id is empty",
		},
	}

	for _, tt := range tests {
		err := Validate(tt.lineID)
		if err != nil {
			if tt.expectError {
				if !strings.Contains(err.Error(), tt.expectErrorContains) {
					t.Errorf(
						"Validation of %q failed with %v, want error containing %q",
						tt.lineID,
						err,
						tt.expectErrorContains)
				}
			} else {
				t.Errorf("Validation of %q failed unexpectedly: %v", tt.lineID, err)
			}
		} else if tt.expectError {
			t.Errorf(
				"Validation of %q succeeded, want error containing %q",
				tt.lineID,
				tt.expectErrorContains)
		}
	}
}

func TestConversion(t *testing.T) {
	tests := []struct {
		name                string
		lineName            string
		expectedLineID      string
		expectError         bool
		expectErrorContains string
	}{
		{
			name:           "Valid name is converted to line id",
			lineName:       "factoryLines/intrinsic-vmp-123",
			expectedLineID: "intrinsic-vmp-123",
		},
		{
			name:                "Invalid name is rejected",
			lineName:            "clusters/vmp-123",
			expectError:         true,
			expectErrorContains: "failed to convert name \"clusters/vmp-123\" to line id",
		},
		{
			name:                "Empty name is rejected",
			lineName:            "",
			expectError:         true,
			expectErrorContains: "failed to convert name \"\" to line id",
		},
	}

	for _, tt := range tests {
		lineID, err := ConvertNameToID(tt.lineName)
		if err != nil {
			if tt.expectError {
				if !strings.Contains(err.Error(), tt.expectErrorContains) {
					t.Errorf(
						"Conversion of %q failed with %v, want error containing %q",
						tt.lineName,
						err,
						tt.expectErrorContains)
				}
			} else {
				t.Errorf("Conversion of %q failed unexpectedly: %v", tt.lineName, err)
			}
		} else if tt.expectError {
			t.Errorf(
				"Conversion of %q succeeded, want error containing %q",
				tt.lineName,
				tt.expectErrorContains)
		} else if lineID != tt.expectedLineID {
			t.Errorf(
				"Line name %q got converted to %q, want %q",
				tt.lineName,
				lineID,
				tt.expectedLineID)
		}
	}
}
