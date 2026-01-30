// Copyright 2023 Intrinsic Innovation LLC

package headers

import (
	"net/http"
	"testing"

	"intrinsic/kubernetes/acl/org"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/metadata"
)

func TestAddOrgToContext(t *testing.T) {
	tests := []struct {
		name     string
		existing metadata.MD
		orgID    string
		expected metadata.MD
	}{
		{
			name:     "add org header",
			existing: metadata.MD{},
			orgID:    "my-org",
			expected: metadata.MD{
				org.OrgIDHeader: []string{"my-org"},
			},
		},
		{
			name: "add to existing headers",
			existing: metadata.MD{
				"other-header": []string{"value"},
			},
			orgID: "my-org",
			expected: metadata.MD{
				"other-header":  []string{"value"},
				org.OrgIDHeader: []string{"my-org"},
			},
		},
		{
			name: "overwrite existing org header",
			existing: metadata.MD{
				org.OrgIDHeader: []string{"old-org"},
			},
			orgID: "new-org",
			expected: metadata.MD{
				org.OrgIDHeader: []string{"new-org"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			if len(tc.existing) > 0 {
				ctx = metadata.NewOutgoingContext(ctx, tc.existing)
			}

			ctx = AddOrgToContext(ctx, tc.orgID)

			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				t.Fatal("AddOrgToContext() did not set outgoing metadata")
			}

			// We need to check if expected metadata is present.
			// Since we might be appending, we should check specifically for our header.
			if diff := cmp.Diff(tc.expected, md); diff != "" {
				t.Errorf("AddOrgToContext() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

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

func TestOrgFromRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *http.Request
		wantOrg *org.Organization
		wantErr bool
	}{
		{
			name:    "nil request",
			req:     nil,
			wantOrg: nil,
			wantErr: false,
		},
		{
			name:    "no org header",
			req:     &http.Request{Header: http.Header{}},
			wantOrg: nil,
			wantErr: false,
		},
		{
			name: "single org header",
			req: &http.Request{
				Header: http.Header{
					http.CanonicalHeaderKey(org.OrgIDHeader): []string{"my-org"},
				},
			},
			wantOrg: &org.Organization{ID: "my-org"},
			wantErr: false,
		},
		{
			name: "multiple org headers",
			req: &http.Request{
				Header: http.Header{
					http.CanonicalHeaderKey(org.OrgIDHeader): []string{"org1", "org2"},
				},
			},
			wantOrg: nil,
			wantErr: true,
		},
		{
			name: "empty org header",
			req: &http.Request{
				Header: http.Header{
					http.CanonicalHeaderKey(org.OrgIDHeader): []string{""},
				},
			},
			wantOrg: nil,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotOrg, err := OrgFromRequest(tc.req)
			if (err != nil) != tc.wantErr {
				t.Errorf("OrgFromRequest() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if diff := cmp.Diff(tc.wantOrg, gotOrg); diff != "" {
				t.Errorf("OrgFromRequest() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
