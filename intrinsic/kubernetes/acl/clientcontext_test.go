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

package clientcontext

import (
	"context"
	"errors"
	"testing"

	"intrinsic/kubernetes/acl/cookies"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestToContextFromIncoming(t *testing.T) {
	tests := []struct {
		desc        string
		incoming    context.Context
		wantMd      metadata.MD
		wantChanged bool
		wantError   bool
	}{
		{
			desc: "empty incoming context",
			incoming: metadata.NewOutgoingContext(
				t.Context(),
				metadata.Pairs(
					"other-key", "other-value",
				),
			),
			wantMd: metadata.MD{
				"other-key": []string{"other-value"},
			},
			wantChanged: false,
		},
		{
			desc: "happy case",
			incoming: metadata.NewIncomingContext(
				metadata.NewOutgoingContext(
					t.Context(),
					metadata.Pairs(
						cookies.CookieHeaderName, "something=somethingelse; something2=somethingelse2",
						"other-key", "other-value",
					),
				),
				metadata.Pairs(
					cookies.CookieHeaderName, "something=somethingelse; something2=somethingelse2",
					AuthHeaderName, "some-token",
					ApikeyTokenHeaderName, "something2",
					"non identity relevant header", "irrelevant-value",
				),
			),
			wantMd: metadata.MD{
				"cookie":              []string{"something=somethingelse; something2=somethingelse2"},
				AuthHeaderName:        []string{"some-token"},
				ApikeyTokenHeaderName: []string{"something2"},
				"other-key":           []string{"other-value"},
			},
			wantChanged: true,
		},
		{
			desc: "duplicate incoming org headers are ignored",
			incoming: metadata.NewIncomingContext(
				t.Context(),
				metadata.Pairs(
					OrgIDHeader, "org1",
					OrgIDHeader, "org1",
				),
			),
			wantMd: metadata.MD{
				OrgIDHeader: []string{"org1"},
			},
			wantChanged: true,
		},
	}

	for _, test := range tests {
		t.Run("ToContextFromIncomingChecked "+test.desc, func(t *testing.T) {
			ctx, changed, err := ToContextFromIncomingChecked(test.incoming)
			if err != nil {
				if test.wantError {
					return
				}
				t.Errorf("ToContextFromIncomingChecked(..) returned error %v, want nil", err)
			}
			if test.wantError {
				t.Errorf("ToContextFromIncomingChecked(..) did not return error, want error")
			}
			gotMd, _ := metadata.FromOutgoingContext(ctx)
			if diff := cmp.Diff(test.wantMd, gotMd, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("metadata.FromOutgoingContext(..) returned diff (-want +got):\n%s", diff)
			}
			if changed != test.wantChanged {
				t.Errorf("ToContextFromIncomingChecked(..) returned changed=%v, want %v", changed, test.wantChanged)
			}
		})

		t.Run("ToContextFromIncoming "+test.desc, func(t *testing.T) {
			ctx, err := ToContextFromIncoming(test.incoming)
			if err != nil {
				if test.wantError {
					return
				}
				t.Errorf("ToContextFromIncoming(..) returned error %v, want nil", err)
			}
			if test.wantError {
				t.Errorf("ToContextFromIncoming(..) did not return error, want error")
			}
			gotMd, _ := metadata.FromOutgoingContext(ctx)
			if diff := cmp.Diff(test.wantMd, gotMd, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("metadata.FromOutgoingContext(..) returned diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestErrGRPC(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantGRPC codes.Code
	}{
		{
			name:     "nil",
			err:      nil,
			wantGRPC: codes.OK,
		},
		{
			name:     "unauthenticated",
			err:      ErrUnauthenticated,
			wantGRPC: codes.Unauthenticated,
		},
		{
			name:     "missing org",
			err:      ErrMissingOrgID,
			wantGRPC: codes.InvalidArgument,
		},
		{
			name:     "invalid request",
			err:      ErrInvalidRequest,
			wantGRPC: codes.InvalidArgument,
		},
		{
			name:     "unknown",
			err:      errors.New("other"),
			wantGRPC: codes.Unknown,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := status.Code(ErrGRPC(tc.err))
			if c != tc.wantGRPC {
				t.Errorf("status.Code(err) = %v, want %v", c, tc.wantGRPC)
			}
		})
	}
}
