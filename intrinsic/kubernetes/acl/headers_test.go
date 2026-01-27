// Copyright 2023 Intrinsic Innovation LLC

package headers

import (
	"net/http"
	"testing"

	"intrinsic/kubernetes/acl/org"

	"github.com/google/go-cmp/cmp"
)

func TestAddOrgToRequest(t *testing.T) {
	tests := []struct {
		name     string
		existing http.Header
		orgID    string
		expected http.Header
	}{
		{
			name:     "add org header",
			existing: http.Header{},
			orgID:    "my-org",
			expected: http.Header{
				http.CanonicalHeaderKey(org.OrgIDHeader): []string{"my-org"},
			},
		},
		{
			name: "overwrite existing org header",
			existing: http.Header{
				http.CanonicalHeaderKey(org.OrgIDHeader): []string{"old-org"},
			},
			orgID: "new-org",
			expected: http.Header{
				http.CanonicalHeaderKey(org.OrgIDHeader): []string{"new-org"},
			},
		},
		{
			name: "preserve other headers",
			existing: http.Header{
				"X-Existing": []string{"ExistingValue"},
			},
			orgID: "my-org",
			expected: http.Header{
				"X-Existing":                             []string{"ExistingValue"},
				http.CanonicalHeaderKey(org.OrgIDHeader): []string{"my-org"},
			},
		},
		{
			name:     "nil request",
			existing: nil, // Special case handled in test
			orgID:    "my-org",
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.name != "nil request" {
				req = &http.Request{Header: tc.existing.Clone()}
			}

			AddOrgToRequest(req, tc.orgID)

			if req == nil {
				if tc.expected != nil {
					t.Errorf("AddOrgToRequest() expected nil request, got %v", req)
				}
				return
			}

			if diff := cmp.Diff(tc.expected, req.Header); diff != "" {
				t.Errorf("AddOrgToRequest() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
