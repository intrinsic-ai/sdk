// Copyright 2023 Intrinsic Innovation LLC

package jwt

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"intrinsic/kubernetes/acl/testing/jwttesting"
)

// Non-standard token with weird payload.
const tokenMixedPayload = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJoZWxsbyI6IndvcmxkIiwiZmllbGQiOjEyMywiZm9yZXN0Ijp0cnVlfQ.ZtB0fzZKCwmC93BB7sS-b2BTw6H53krjnOMc94Y1Kto"

var (
	tokenEmail  = jwttesting.MustMintToken(jwttesting.WithEmail("doe@example.com"), jwttesting.WithEmailVerified(true), jwttesting.WithAudience("intrinsic-portal-prod"))
	tokenUID    = jwttesting.MustMintToken(jwttesting.WithUID("doe@example.com"), jwttesting.WithAudience("intrinsic-portal-prod"))
	tokenExpiry = jwttesting.MustMintToken(jwttesting.WithExpiresAt(time.Unix(420, 0)))
)

func TestUnmarshalUnsafe(t *testing.T) {
	tests := []struct {
		desc string
		jwtk string
		want *Data
	}{
		{
			desc: "happy case",
			jwtk: jwttesting.MustMintToken(
				jwttesting.WithAudience("testaud"),
				jwttesting.WithAuthorized(true),
				jwttesting.WithEmail("test@domain"),
				jwttesting.WithEmailVerified(true),
				jwttesting.WithExpiresAt(time.Unix(123, 0)),
				// some extra data we do not care about right now
				jwttesting.WithIssuedAt(time.UnixMilli(0)),
			),
			want: &Data{
				Aud:           "testaud",
				Authorized:    true,
				Email:         "test@domain",
				EmailVerified: true,
				ExpiresAt:     123,
			},
		},
		{
			desc: "empty",
			jwtk: jwttesting.MustMintToken(),
			want: &Data{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := UnmarshalUnsafe(tc.jwtk)
			if err != nil {
				t.Fatalf("ParseUnsafe(%v) returned an unexpected error: %v", tc.jwtk, err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ParseUnsafe(%v) returned an unexpected diff (-want +got): %v", tc.jwtk, diff)
			}
		})
	}
}

func TestPayloadUnsafe(t *testing.T) {
	tests := []struct {
		desc        string
		token       string
		wantErr     error
		wantPayload map[string]any
	}{
		{
			desc:  "happy path 1",
			token: tokenMixedPayload,
			wantPayload: map[string]any{
				"hello":  "world",
				"field":  123.0,
				"forest": true,
			},
		},
		{
			desc:  "happy path 2",
			token: tokenEmail,
			wantPayload: map[string]any{
				"email":          "doe@example.com",
				"email_verified": true,
				"aud":            "intrinsic-portal-prod",
				"uid":            "",
			},
		},
		{
			desc:  "happy path 3",
			token: tokenUID,
			wantPayload: map[string]any{
				"email":          "",
				"email_verified": false,
				"aud":            "intrinsic-portal-prod",
				"uid":            "doe@example.com",
			},
		},
		{
			desc:  "happy path 4",
			token: tokenExpiry,
			wantPayload: map[string]any{
				"email":          "",
				"email_verified": false,
				"exp":            float64(420),
				"uid":            "",
			},
		},
		{
			desc:    "no three parts",
			token:   "foobarbaz",
			wantErr: cmpopts.AnyError,
		},
		{
			desc:    "second part not base64",
			token:   "foo.!!!.baz",
			wantErr: cmpopts.AnyError,
		},
		{
			desc:    "second part not a valid JSON",
			token:   "foo.SSdtIG5vdCBhIEpTT04g8J-Ymw.baz",
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			gotPayload, err := PayloadUnsafe(tc.token)
			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("PayloadUnsafe(%q): unexpected error, diff (-want,+got):\n%s", tc.token, diff)
			}
			if diff := cmp.Diff(tc.wantPayload, gotPayload); diff != "" {
				t.Errorf("PayloadUnsafe(%q): unexpected payload, diff (-want,+got):\n%s", tc.token, diff)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		desc      string
		token     string
		wantErr   error
		wantEmail string
	}{
		{
			desc:      "happy path with email",
			token:     tokenEmail,
			wantEmail: "doe@example.com",
		},
		{
			desc:      "happy path with UID",
			token:     tokenUID,
			wantEmail: "doe@example.com",
		},
		{
			desc:    "no three parts",
			token:   "foobarbaz",
			wantErr: cmpopts.AnyError,
		},
		{
			desc:    "second part not base64",
			token:   "foo.!!!.baz",
			wantErr: cmpopts.AnyError,
		},
		{
			desc:    "second part not a valid JSON",
			token:   "foo.SSdtIG5vdCBhIEpTT04g8J-Ymw.baz",
			wantErr: cmpopts.AnyError,
		},
		{
			desc:    "valid payload, but no email or uid field",
			token:   tokenMixedPayload,
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			gotEmail, err := Email(tc.token)
			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Email(%q): unexpected error, diff (-want,+got):\n%s", tc.token, diff)
			}
			if gotEmail != tc.wantEmail {
				t.Errorf("Email(%q)=%q, want %q", tc.token, gotEmail, tc.wantEmail)
			}
		})
	}
}

func TestAud(t *testing.T) {
	tests := []struct {
		desc    string
		token   string
		wantErr error
		wantAud string
	}{
		{
			desc:    "happy path with email",
			token:   tokenEmail,
			wantAud: "intrinsic-portal-prod",
		},
		{
			desc:    "happy path with UID",
			token:   tokenUID,
			wantAud: "intrinsic-portal-prod",
		},
		{
			desc:    "no three parts",
			token:   "foobarbaz",
			wantErr: cmpopts.AnyError,
		},
		{
			desc:    "second part not base64",
			token:   "foo.!!!.baz",
			wantErr: cmpopts.AnyError,
		},
		{
			desc:    "second part not a valid JSON",
			token:   "foo.SSdtIG5vdCBhIEpTT04g8J-Ymw.baz",
			wantErr: cmpopts.AnyError,
		},
		{
			desc:    "valid payload, but no aud field",
			token:   tokenMixedPayload,
			wantErr: cmpopts.AnyError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			gotAud, err := Aud(tc.token)
			if diff := cmp.Diff(tc.wantErr, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("Aud(%q): unexpected error, diff (-want,+got):\n%s", tc.token, diff)
			}
			if gotAud != tc.wantAud {
				t.Errorf("Aud(%q)=%q, want %q", tc.token, gotAud, tc.wantAud)
			}
		})
	}
}

type IsVerifiedAndAuthorizedUnsafeTest struct {
	desc    string
	token   string
	wantErr bool
}

func TestIsVerifiedAndAuthorizedUnsafe(t *testing.T) {
	tests := []IsVerifiedAndAuthorizedUnsafeTest{
		{
			desc:    "happy path",
			token:   jwttesting.MustMintToken(jwttesting.WithEmail("doe@example.com"), jwttesting.WithEmailVerified(true), jwttesting.WithAuthorized(true)),
			wantErr: false,
		},
		{
			desc:    "not verified",
			token:   jwttesting.MustMintToken(jwttesting.WithEmail("doe@example.com"), jwttesting.WithEmailVerified(false), jwttesting.WithAuthorized(true)),
			wantErr: true,
		},
		{
			desc:    "not authorized",
			token:   jwttesting.MustMintToken(jwttesting.WithEmail("doe@example.com"), jwttesting.WithEmailVerified(true), jwttesting.WithAuthorized(false)),
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			runTestIsVerifiedAndAuthorizedUnsafe(t, &tc)
		})
	}
}

func runTestIsVerifiedAndAuthorizedUnsafe(t *testing.T, tc *IsVerifiedAndAuthorizedUnsafeTest) {
	err := IsVerifiedAndAuthorizedUnsafe(tc.token)
	if (err != nil) != tc.wantErr {
		t.Fatalf("IsVerifiedAndAuthorizedUnsafe(%q) returned an unexpected error: %v, wantErr: %v", tc.token, err, tc.wantErr)
	}
}

func TestValueStringHappyCase(t *testing.T) {
	tk := jwttesting.MustMintToken(jwttesting.WithUserID("testuid"))
	d, err := UnmarshalUnsafe(tk)
	if err != nil {
		t.Fatalf("failed to unmarshal token: %v", err)
	}
	if d.UserID != "testuid" {
		t.Errorf("UID is not as expected: got %q, want %q", d.UserID, "testuid")
	}
}
