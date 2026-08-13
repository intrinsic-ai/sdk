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

package kvstore

import (
	"testing"
)

func TestMakeKey(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		expected string
	}{
		{
			name:     "some parts without slashes and some without",
			parts:    []string{"/foo", "bar", "baz/"},
			expected: "foo/bar/baz",
		},
		{
			name:     "no slashes present",
			parts:    []string{"foo", "bar", "baz"},
			expected: "foo/bar/baz",
		},
		{
			name:     "multiple slashes",
			parts:    []string{"///foo", "bar///", "///baz///"},
			expected: "foo/bar/baz",
		},
		{
			name:     "single slashes on both sides",
			parts:    []string{"/foo/", "/bar/", "/baz/"},
			expected: "foo/bar/baz",
		},
		{
			name:     "empty part in middle",
			parts:    []string{"foo", "", "bar"},
			expected: "foo/bar",
		},
		{
			name:     "slashes only part in middle",
			parts:    []string{"foo", "///", "bar"},
			expected: "foo/bar",
		},
		{
			name:     "only slashes across all parts",
			parts:    []string{"///", "///", "///"},
			expected: "",
		},
		{
			name:     "no parts",
			parts:    []string{},
			expected: "",
		},
		{
			name:     "internal slash preserved",
			parts:    []string{"foo/bar", "baz"},
			expected: "foo/bar/baz",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := MakeKey(tc.parts...)
			if actual != tc.expected {
				t.Errorf("MakeKey() = %v, expected %v", actual, tc.expected)
			}
		})
	}
}
