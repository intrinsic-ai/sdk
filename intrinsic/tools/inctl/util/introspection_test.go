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

package introspection

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCloudTraceURL(t *testing.T) {
	tests := []struct {
		desc    string
		project string
		trace   string
		want    string
	}{
		{
			desc: "empty values blindly generate invalid url",
			want: "https://console.cloud.google.com/traces/explorer;traceId=?project=",
		},
		{
			desc:    "provided values generate valid url",
			project: "the-shire",
			trace:   "1420",
			want:    "https://console.cloud.google.com/traces/explorer;traceId=1420?project=the-shire",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := CloudTraceURL(tc.trace, tc.project)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("CloudTraceURL(%q, %q) returned unexpected, diff (-want +got):\n%s", tc.trace, tc.project, diff)
			}
		})
	}
}
