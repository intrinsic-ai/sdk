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

package device

import (
	"testing"
)

func TestValidHostname(t *testing.T) {
	testCases := []struct {
		name     string
		hostname string
		valid    bool
		index    int
	}{
		{
			name:     "empty",
			hostname: "",
			valid:    false,
			index:    0,
		},
		{
			name:     "single-char",
			hostname: "h",
			valid:    true,
		},
		{
			name:     "simple",
			hostname: "host",
			valid:    true,
		},
		{
			name:     "complex",
			hostname: "host-13245-123413240-998",
			valid:    true,
		},
		{
			name:     "generated",
			hostname: "node-7da4bec3-3aa2-4859-97dd-b07dcce37c19",
			valid:    true,
		},
		{
			name:     "with-alpha",
			hostname: "host123",
			valid:    true,
		},
		{
			name:     "with-dash",
			hostname: "host-123",
			valid:    true,
		},
		{
			name:     "no-starting-dash",
			hostname: "-host-123",
			valid:    false,
			index:    0,
		},
		{
			name:     "no-ending-dash",
			hostname: "host-123-",
			valid:    false,
			index:    9,
		},
		{
			name:     "no-capital",
			hostname: "HOST-123",
			valid:    false,
		},
		{
			name:     "no-underscore",
			hostname: "host_123",
			valid:    false,
			index:    5,
		},
		{
			name:     "no-overlength",
			hostname: "host-123456790-123456790-123456790-123456790-123456790-123456790",
			valid:    false,
			index:    43,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			indexGot, validGot := validHostname(tc.hostname)
			if got, want := validGot, tc.valid; got != want {
				t.Errorf("validHostname(%q) = %t, want %t", tc.hostname, got, want)
			}

			if got, want := indexGot, tc.index; got != want {
				t.Errorf("validHostname(%q).index = %v, want %v", tc.hostname, got, want)
			}
		})
	}
}
