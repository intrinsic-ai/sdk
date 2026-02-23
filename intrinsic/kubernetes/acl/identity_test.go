// Copyright 2023 Intrinsic Innovation LLC

package identity

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"intrinsic/kubernetes/acl/cookies"
	"intrinsic/kubernetes/acl/org"
	"intrinsic/kubernetes/acl/testing/jwttesting"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	aud    = "example-project"
	aud2   = "example-project"
	email  = "doe@example.com"
	email2 = "ray@example.com"
	email3 = "ray.ray@example.com"
)

var (
	token  = jwttesting.MustMintToken(jwttesting.WithEmail(email), jwttesting.WithAudience(aud))
	token2 = jwttesting.MustMintToken(jwttesting.WithEmail(email2), jwttesting.WithAudience(aud2))
	token3 = jwttesting.MustMintToken(jwttesting.WithEmail(email3), jwttesting.WithAudience(aud2))
)

func TestUserFromJWT(t *testing.T) {
	uid, err := UserFromJWT(token)
	if err != nil {
		t.Fatal(err)
	}
	if uid.EmailCanonicalized() != email {
		t.Fatalf("email, got %s, want %s", uid.EmailCanonicalized(), email)
	}
}

func TestUserFromJWTRaw(t *testing.T) {
	uid, err := UserFromJWT(token3)
	if err != nil {
		t.Fatal(err)
	}
	if uid.EmailRaw() != email3 {
		t.Fatalf("email, got %s, want %s", uid.EmailRaw(), email3)
	}
}

func TestUserFromRequest(t *testing.T) {
	t.Run("cookie", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{Name: AuthProxyCookieName, Value: token})
		u, err := UserFromRequest(r)
		if err != nil {
			t.Fatal(err)
		}
		if u.EmailCanonicalized() != email {
			t.Fatalf("email, got %s, want %s", u.EmailCanonicalized(), email)
		}
	})

	t.Run("portal-cookie", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{Name: portalCookieName, Value: token})
		u, err := UserFromRequest(r)
		if err != nil {
			t.Fatal(err)
		}
		if u.Email() != email {
			t.Fatalf("email, got %s, want %s", u.Email(), email)
		}
	})

	t.Run("header", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Add(ApikeyTokenHeaderName, token)
		u, err := UserFromRequest(r)
		if err != nil {
			t.Fatal(err)
		}
		if u.EmailCanonicalized() != email {
			t.Fatalf("email, got %s, want %s", u.EmailCanonicalized(), email)
		}
	})
}

func TestUserToRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	u := &User{jwt: token}
	UserToRequest(r, &User{jwt: "overwrite me please"})
	UserToRequest(r, u)
	if cookies.FromRequestNamed(r, []string{AuthProxyCookieName})[0].Value != token {
		t.Errorf("UserToRequest(..) did not add the user's identity to the request")
	}
	if len(r.Cookies()) != 1 {
		t.Errorf("UserToRequest(..) did not add exactly one cookie to the request")
	}
}

func TestUserFromMetadata(t *testing.T) {
	ctx := metadata.NewIncomingContext(t.Context(),
		metadata.Pairs(cookies.CookieHeaderName, AuthProxyCookieName+"="+token))

	u, err := UserFromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if u.EmailCanonicalized() != email {
		t.Fatalf("email, got %s, want %s", u.EmailCanonicalized(), email)
	}
}

func TestUserToContext(t *testing.T) {
	ctx := t.Context()
	u := &User{jwt: token}
	ctx, err := UserToContext(ctx, u)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{AuthProxyCookieName + "=" + token}
	md, _ := metadata.FromOutgoingContext(ctx)
	if diff := cmp.Diff(want, md.Get(cookies.CookieHeaderName)); diff != "" {
		t.Errorf("UserToContext(..) did not add the user's identity to the context (-want +got):\n%s", diff)
	}

	// Overwrite existing cookie.
	u.jwt = token2
	ctx, err = UserToContext(ctx, u)
	if err != nil {
		t.Fatal(err)
	}
	want = []string{AuthProxyCookieName + "=" + token2}
	md, _ = metadata.FromOutgoingContext(ctx)
	if diff := cmp.Diff(want, md.Get(cookies.CookieHeaderName)); diff != "" {
		t.Errorf("UserToContext(..) did not add the user's identity to the context (-want +got):\n%s", diff)
	}
}

func TestUserToIncomingContext(t *testing.T) {
	ctx := t.Context()
	u := &User{jwt: token}
	ctx, err := UserToIncomingContext(ctx, u)
	if err != nil {
		t.Fatal(err)
	}
	md, _ := metadata.FromIncomingContext(ctx)
	if md.Get(cookies.CookieHeaderName)[0] != AuthProxyCookieName+"="+token {
		t.Errorf("UserToIncomingContext(..) did not add the user's identity to the context")
	}
}

type RequestToContextTest struct {
	desc    string
	cookies map[string]string
	headers map[string]string
	wantMd  map[string]string
}

func TestRequestToContext(t *testing.T) {
	tests := []RequestToContextTest{
		{
			desc: "just auth-proxy cookie",
			cookies: map[string]string{
				"auth-proxy": token,
			},
			wantMd: map[string]string{
				"x-intrinsic-auth-project": aud,
				"cookie":                   "auth-proxy=" + token,
			},
		},
		{
			desc: "just portal-token cookie",
			cookies: map[string]string{
				"portal-token": token,
			},
			wantMd: map[string]string{
				"x-intrinsic-auth-project": aud,
				"cookie":                   "auth-proxy=" + token,
			},
		},
		{
			desc: "just auth-proxy cookie, custom org-id",
			cookies: map[string]string{
				"auth-proxy": token,
				"org-id":     "customorg",
			},
			wantMd: map[string]string{
				"x-intrinsic-auth-project": aud,
				"cookie":                   "auth-proxy=" + token + "; org-id=customorg",
			},
		},
		{
			desc: "auth-proxy cookie + extra cookie",
			cookies: map[string]string{
				"auth-proxy":   token,
				"extra-cookie": "random",
			},
			wantMd: map[string]string{
				"x-intrinsic-auth-project": aud,
				"cookie":                   "auth-proxy=" + token,
			},
		},
		{
			desc: "backfill org-id cookie from header",
			cookies: map[string]string{
				"auth-proxy": token,
			},
			headers: map[string]string{
				"org-id": "customorg",
			},
			wantMd: map[string]string{
				"x-intrinsic-auth-project": aud,
				"cookie":                   "auth-proxy=" + token + "; org-id=customorg",
			},
		},
		{
			desc: "x-intrinsic-org header",
			cookies: map[string]string{
				"auth-proxy": token,
			},
			headers: map[string]string{
				"x-intrinsic-org": "neworg",
			},
			wantMd: map[string]string{
				"x-intrinsic-auth-project": aud,
				"cookie":                   "auth-proxy=" + token,
				"x-intrinsic-org":          "neworg",
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			runRequestToContextTest(t, &tc)
			runRequestToIncomingContextTest(t, &tc)
		})
	}
}

func createRequest(test *RequestToContextTest) *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	for k, v := range test.cookies {
		r.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	for k, v := range test.headers {
		r.Header.Add(k, v)
	}
	return r
}

func runRequestToContextTest(t *testing.T, test *RequestToContextTest) {
	r := createRequest(test)
	ctx, err := RequestToContext(r.Context(), r)
	if err != nil {
		t.Fatal(err)
	}
	gotMd, _ := metadata.FromOutgoingContext(ctx)
	metadataTest(t, test.wantMd, gotMd)
}

func runRequestToIncomingContextTest(t *testing.T, test *RequestToContextTest) {
	r := createRequest(test)
	ctx, err := RequestToIncomingContext(r.Context(), r)
	if err != nil {
		t.Fatal(err)
	}
	gotMd, _ := metadata.FromIncomingContext(ctx)
	metadataTest(t, test.wantMd, gotMd)
}

func metadataTest(t *testing.T, wantMd map[string]string, gotMd metadata.MD) {
	if len(wantMd) != len(gotMd) {
		fmt.Printf("%+v\n", gotMd)
		t.Errorf("len(metadata): got %d (%v), want %d", len(gotMd), gotMd, len(wantMd))
	}
	for k, wantV := range wantMd {
		mdValue, found := gotMd[k]
		if !found {
			t.Errorf("missin key %s in metadata", k)
		}
		if len(mdValue) != 1 {
			t.Errorf("len(mdValue): got %d (%v), want 1", len(mdValue), mdValue)
		} else {
			if mdValue[0] != wantV {
				t.Errorf("mdValue[%s]: got %q, want %q", k, mdValue, wantV)
			}
		}
	}
}

func TestEnsureAuthProxyCookie(t *testing.T) {
	requestWithAuthProxy := httptest.NewRequest(http.MethodGet, "/", nil)
	requestWithAuthProxy.AddCookie(&http.Cookie{Name: AuthProxyCookieName, Value: token})

	requestWithOnPrem := httptest.NewRequest(http.MethodGet, "/", nil)
	requestWithOnPrem.AddCookie(&http.Cookie{Name: onpremTokenCookieName, Value: token2})

	requestWithPortal := httptest.NewRequest(http.MethodGet, "/", nil)
	requestWithPortal.AddCookie(&http.Cookie{Name: portalCookieName, Value: token2})

	requestWithBoth := httptest.NewRequest(http.MethodGet, "/", nil)
	requestWithBoth.AddCookie(&http.Cookie{Name: AuthProxyCookieName, Value: token})
	requestWithBoth.AddCookie(&http.Cookie{Name: onpremTokenCookieName, Value: token2})

	tests := []struct {
		desc      string
		request   *http.Request
		wantError bool
		wantToken string
	}{
		{
			desc:      "has_auth_proxy_cookie",
			request:   requestWithAuthProxy,
			wantError: false,
			wantToken: token,
		},
		{
			desc:      "no_cookie",
			request:   httptest.NewRequest(http.MethodGet, "/", nil),
			wantError: true,
		},
		{
			desc:      "has_onprem_cookie",
			request:   requestWithOnPrem,
			wantError: false,
			wantToken: token2,
		},
		{
			desc:      "has_portal_cookie",
			request:   requestWithPortal,
			wantError: false,
			wantToken: token2,
		},
		{
			desc:      "has_both_cookies",
			request:   requestWithBoth,
			wantError: false,
			wantToken: token,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			if err := EnsureAuthProxyCookie(tc.request); (err != nil) != tc.wantError {
				t.Errorf("EnsureAuthProxyCookie(%v) = %v, wantError=%v", tc.request, err, tc.wantError)
			}

			if !tc.wantError {
				cookie, err := tc.request.Cookie(AuthProxyCookieName)
				if err != nil {
					t.Errorf("EnsureAuthProxyCookie(%v) has no auth-proxy cookie", tc.request)
				}
				if cookie.Value != tc.wantToken {
					t.Errorf("EnsureAuthProxyCookie(%v) auth-proxy token = %q, want %q", tc.request, cookie.Value, tc.wantToken)
				}
			}
		})
	}
}

func TestIsIntrinsicUserDiscriminates(t *testing.T) {
	tests := []struct {
		desc      string
		email     string
		intrinsic bool
	}{
		{
			desc:      "intrinsic user",
			email:     "bender@google.com",
			intrinsic: true,
		},
		{
			desc:      "non-intrinsic user",
			email:     "bender@bing.com",
			intrinsic: false,
		},
		{
			desc:      "almost-intrinsic user",
			email:     "bender@ggoogle.com",
			intrinsic: false,
		},
		{
			desc:      "intrinsic domain in prefix",
			email:     "google.com@hacker.com",
			intrinsic: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			isIntrinsic := isIntrinsicUser(tc.email)
			if isIntrinsic != tc.intrinsic {
				t.Errorf("IsIntrinsicUser(%s) = %v, want %v", tc.email, isIntrinsic, tc.intrinsic)
			}
		})
	}
}

type FakeVerifier struct {
	ReturnErr error
	WantToken string
}

func (f *FakeVerifier) VerifyIDToken(ctx context.Context, t string) error {
	if t != f.WantToken {
		return fmt.Errorf("unexpected token: got %q, want %q", t, f.WantToken)
	}
	return f.ReturnErr
}

type UserFromRequestVerifiedTest struct {
	desc      string
	r         *http.Request
	verifyErr error
	wantToken string
	wantEmail string
	wantErr   bool
}

func TestUserFromRequestVerified(t *testing.T) {
	withAuthCookie := httptest.NewRequest(http.MethodGet, "/", nil)
	withAuthCookie.AddCookie(&http.Cookie{Name: AuthProxyCookieName, Value: token})
	tests := []UserFromRequestVerifiedTest{
		{
			desc:    "no auth",
			r:       httptest.NewRequest(http.MethodGet, "/", nil),
			wantErr: true,
		},
		{
			desc:      "with auth cookie",
			r:         withAuthCookie,
			wantToken: token,
			wantEmail: email,
			wantErr:   false,
		},
		{
			desc:      "with auth cookie but verify error",
			r:         withAuthCookie,
			wantToken: token,
			wantEmail: email,
			wantErr:   true,
			verifyErr: fmt.Errorf("test verify error"),
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) { runUserFromRequestVerifiedTest(t, &test) })
	}
}

func runUserFromRequestVerifiedTest(t *testing.T, test *UserFromRequestVerifiedTest) {
	fv := &FakeVerifier{ReturnErr: test.verifyErr, WantToken: test.wantToken}
	u, err := UserFromRequestVerified(test.r, fv)
	if (err != nil) != test.wantErr {
		t.Errorf("UserFromRequestVerified(..) = _, %v, wantError=%v", err, test.wantErr)
	}
	if err != nil {
		return
	}
	if u.EmailCanonicalized() != test.wantEmail {
		t.Errorf("UserFromRequestVerified(..) = %q, want %q", u.EmailCanonicalized(), test.wantEmail)
	}
}

type UserFromContextVerifiedTest struct {
	desc      string
	ctx       context.Context
	verifyErr error
	wantToken string
	wantEmail string
	wantErr   bool
}

func TestUserFromContextVerified(t *testing.T) {
	withAuthCookie := metadata.NewIncomingContext(t.Context(),
		metadata.Pairs(cookies.CookieHeaderName, AuthProxyCookieName+"="+token))
	tests := []UserFromContextVerifiedTest{
		{
			desc:    "no auth",
			ctx:     t.Context(),
			wantErr: true,
		},
		{
			desc:      "with auth cookie",
			ctx:       withAuthCookie,
			wantToken: token,
			wantEmail: email,
		},
		{
			desc:      "with auth cookie but verify error",
			ctx:       withAuthCookie,
			verifyErr: fmt.Errorf("test verify error"),
			wantToken: token,
			wantEmail: email,
			wantErr:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) { runUserFromContextVerifiedTest(t, &test) })
	}
}

func runUserFromContextVerifiedTest(t *testing.T, test *UserFromContextVerifiedTest) {
	fv := &FakeVerifier{ReturnErr: test.verifyErr, WantToken: test.wantToken}
	u, err := UserFromContextVerified(test.ctx, fv)
	if (err != nil) != test.wantErr {
		t.Errorf("FromContextVerified(..) = _, %v, wantError=%v", err, test.wantErr)
	}
	if err != nil {
		return
	}
	if u.EmailCanonicalized() != test.wantEmail {
		t.Errorf("FromContextVerified(..) = %q, want %q", u.EmailCanonicalized(), test.wantEmail)
	}
}

func TestCanonicalizeEmailInvalid(t *testing.T) {
	tests := []string{
		"",
		"john",
		"john@",
		"@gmail.com",
		"@",
		"a@b@c",
	}

	for _, tc := range tests {
		_, err := CanonicalizeEmail(tc)
		if err == nil {
			t.Errorf("CanonicalizeEmail(%q) did not return error", tc)
		}
	}
}

func TestCanonicalizeEmail(t *testing.T) {
	tests := []struct {
		email string
		want  string
	}{
		{
			email: "doe@gmail.com",
			want:  "doe@gmail.com",
		},
		{
			email: "john.doe@gmail.com",
			want:  "johndoe@gmail.com",
		},
		{
			email: ".john..doe.@gmail.com",
			want:  "johndoe@gmail.com",
		},
		{
			email: "John.Doe@gmail.com",
			want:  "johndoe@gmail.com",
		},
		{
			email: "doe+foo@gmail.com",
			want:  "doe@gmail.com",
		},
		{
			email: "doe@googlemail.com",
			want:  "doe@gmail.com",
		},
		{
			email: "!john.doe#@gmail.com",
			want:  "johndoe@gmail.com",
		},
		{
			email: "!john.doe#@yahoo.com",
			want:  "!johndoe#@yahoo.com",
		},
	}

	for _, tc := range tests {
		got, _ := CanonicalizeEmail(tc.email)
		if got != tc.want {
			t.Errorf("CanonicalizeEmail(%q) = %q, want: %q", tc.email, got, tc.want)
		}
	}
}

type ToContextFromIncomingTest struct {
	desc        string
	incoming    context.Context
	wantMd      metadata.MD
	wantChanged bool
	wantError   bool
}

func TestToContextFromIncoming(t *testing.T) {
	tests := []ToContextFromIncomingTest{
		{
			desc:     "blank context",
			incoming: t.Context(),
			wantMd:   metadata.MD{},
		},
		{
			desc: "keep metadata in outgoing if incoming has nothing to copyover",
			incoming: metadata.NewIncomingContext(
				metadata.NewOutgoingContext(t.Context(),
					metadata.Pairs(cookies.CookieHeaderName, "something"),
				),
				metadata.Pairs(),
			),
			wantMd: metadata.MD{
				"cookie": []string{"something"},
			},
		},
		{
			desc: "copy cookie from incoming to outgoing",
			incoming: metadata.NewIncomingContext(t.Context(),
				metadata.Pairs(
					cookies.CookieHeaderName, "something=somethingelse",
				),
			),
			wantMd: metadata.MD{
				"cookie": []string{"something=somethingelse"},
			},
			wantChanged: true,
		},
		{
			desc: "incoming auth-proxy and cookie with multiple values",
			incoming: metadata.NewIncomingContext(t.Context(),
				metadata.Pairs(
					AuthProxyCookieName, "anything", // do not copy over deprecated auth-proxy
					cookies.CookieHeaderName, "something",
					cookies.CookieHeaderName, "something2",
				),
			),
			wantMd: metadata.MD{
				"cookie": []string{"something", "something2"},
			},
			wantChanged: true,
		},
		{
			desc: "incoming authorization header",
			incoming: metadata.NewIncomingContext(t.Context(),
				metadata.Pairs(authHeaderName, "some-token"),
			),
			wantMd: metadata.MD{
				authHeaderName: []string{"some-token"},
			},
			wantChanged: true,
		},
		{
			desc: "multiple incoming authorization headers are ignored",
			incoming: metadata.NewIncomingContext(t.Context(),
				metadata.Pairs(
					authHeaderName, "some-token",
					authHeaderName, "other-token",
				),
			),
			wantError: true,
		},
		{
			desc: "cookie comparison",
			incoming: metadata.NewIncomingContext(
				metadata.NewOutgoingContext(t.Context(),
					metadata.MD{
						"cookie": []string{"something=somethingelse", "something2=somethingelse2"},
					},
				),
				metadata.MD{
					"cookie":       []string{"something2=somethingelse2", "something=somethingelse"},
					authHeaderName: []string{"some-token"},
				},
			),
			wantMd: metadata.MD{
				"cookie":       []string{"something=somethingelse", "something2=somethingelse2"},
				authHeaderName: []string{"some-token"},
			},
			wantChanged: true,
		},
		{
			desc: "cookie comparison reverse",
			incoming: metadata.NewIncomingContext(
				metadata.NewOutgoingContext(t.Context(),
					metadata.MD{
						"cookie": []string{"something2=somethingelse2", "something=somethingelse"},
					},
				),
				metadata.MD{
					"cookie":       []string{"something=somethingelse", "something2=somethingelse2"},
					authHeaderName: []string{"some-token"},
				},
			),
			wantMd: metadata.MD{
				"cookie":       []string{"something2=somethingelse2", "something=somethingelse"},
				authHeaderName: []string{"some-token"},
			},
			wantChanged: true,
		},
		{
			desc: "cookies are not compared inside the cookie-MDstring",
			incoming: metadata.NewIncomingContext(
				metadata.NewOutgoingContext(t.Context(),
					metadata.MD{
						"cookie": []string{"something=somethingelse; something2=somethingelse2"},
					},
				),
				metadata.MD{
					"cookie": []string{"something2=somethingelse2; something=somethingelse"},
				},
			),
			wantError: true,
		},
		{
			desc: "incoming authorization header collides with outgoing",
			incoming: metadata.NewIncomingContext(
				metadata.NewOutgoingContext(t.Context(), metadata.Pairs(authHeaderName, "some-token")),
				metadata.Pairs(
					authHeaderName, "other-token",
				),
			),
			wantError: true,
		},
		{
			desc: "incoming malformed cookies",
			incoming: metadata.NewIncomingContext(
				metadata.NewOutgoingContext(
					t.Context(), metadata.Pairs(cookies.CookieHeaderName, "anything"),
				),
				metadata.Pairs(cookies.CookieHeaderName, "something"),
			),
			wantError: true,
		},
		{
			desc: "multiple incoming apikey tokens",
			incoming: metadata.NewIncomingContext(
				t.Context(),
				metadata.Pairs(ApikeyTokenHeaderName, "something", ApikeyTokenHeaderName, "something2"),
			),
			wantError: true,
		},
		{
			desc: "incoming apikey token and outgoing collide",
			incoming: metadata.NewIncomingContext(
				metadata.NewOutgoingContext(t.Context(), metadata.Pairs(ApikeyTokenHeaderName, "anything")),
				metadata.Pairs(ApikeyTokenHeaderName, "something"),
			),
			wantError: true,
		},
		{
			desc: "incoming org header",
			incoming: metadata.NewIncomingContext(t.Context(),
				metadata.Pairs(org.OrgIDHeader, "org1"),
			),
			wantMd: metadata.MD{
				org.OrgIDHeader: []string{"org1"},
			},
			wantChanged: true,
		},
		{
			desc: "multiple incoming org headers are ignored",
			incoming: metadata.NewIncomingContext(t.Context(),
				metadata.Pairs(
					org.OrgIDHeader, "org1",
					org.OrgIDHeader, "org2",
				),
			),
			wantError: true,
		},
		{
			desc: "incoming org header collides with outgoing",
			incoming: metadata.NewIncomingContext(
				metadata.NewOutgoingContext(t.Context(), metadata.Pairs(org.OrgIDHeader, "org1")),
				metadata.Pairs(
					org.OrgIDHeader, "org2",
				),
			),
			wantError: true,
		},
		{
			desc: "happy case - merge incoming metadata with outgoing metadata",
			incoming: metadata.NewIncomingContext(
				metadata.NewOutgoingContext(
					t.Context(),
					metadata.Pairs(
						cookies.CookieHeaderName, "something=somethingelse; something2=somethingelse2",
						"other-key", "other-value", // present pairs in outgoing should not be overwritten
					),
				),
				metadata.Pairs(
					cookies.CookieHeaderName, "something=somethingelse; something2=somethingelse2",
					authHeaderName, "some-token",
					ApikeyTokenHeaderName, "something2",
					"non identity relevant header", "irrelevant-value",
				),
			),
			wantMd: metadata.MD{
				"cookie":              []string{"something=somethingelse; something2=somethingelse2"},
				authHeaderName:        []string{"some-token"},
				ApikeyTokenHeaderName: []string{"something2"},
				"other-key":           []string{"other-value"},
			},
			wantChanged: true,
		},
		{
			desc: "happy case - multiple cookies",
			incoming: metadata.NewIncomingContext(
				t.Context(),
				metadata.MD{
					cookies.CookieHeaderName: []string{"something", "something2"},
				},
			),
			wantMd: metadata.MD{
				"cookie": []string{"something", "something2"},
			},
			wantChanged: true,
		},
		{
			desc: "duplicate incoming org headers are ignored",
			incoming: metadata.NewIncomingContext(t.Context(),
				metadata.Pairs(
					org.OrgIDHeader, "org1",
					org.OrgIDHeader, "org1",
				),
			),
			wantMd: metadata.MD{
				org.OrgIDHeader: []string{"org1"},
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

func TestContextToRequest(t *testing.T) {
	testCases := []struct {
		name   string
		noAuth bool
		meta   metadata.MD
	}{
		{
			name:   "no-auth",
			noAuth: true,
		},
		{
			name:   "with-auth-coookie",
			noAuth: false,
			meta:   metadata.New(map[string]string{cookies.CookieHeaderName: AuthProxyCookieName + "=" + token + "; " + org.OrgIDCookie + "=" + "testorg"}),
		},
		{
			name:   "with-duplicated-auth-coookie",
			noAuth: false,
			meta:   metadata.New(map[string]string{cookies.CookieHeaderName: AuthProxyCookieName + "=" + token2 + "; " + AuthProxyCookieName + "=" + token + "; " + org.OrgIDCookie + "=" + "wrongorg" + "; " + org.OrgIDCookie + "=" + "testorg"}),
		},
		{
			name:   "with-apikey-meta",
			noAuth: false,
			meta: metadata.New(map[string]string{
				ApikeyTokenHeaderName:    token,
				cookies.CookieHeaderName: org.OrgIDCookie + "=testorg",
			}),
		},
		{
			name:   "with-org-header",
			noAuth: false,
			meta: metadata.New(map[string]string{
				ApikeyTokenHeaderName:    token,
				cookies.CookieHeaderName: org.OrgIDCookie + "=testorg",
				org.OrgIDHeader:          "headerorg",
			}),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(t.Context(), tc.meta)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			err := ContextToRequest(ctx, req)
			if err != nil {
				t.Errorf("ToRequest(..) = _, %v, wantError=%v", err, tc.noAuth)
			}

			if tc.noAuth {
				return
			}

			uo, err := UserOrgFromRequest(req)
			if err != nil {
				t.Errorf("FromRequest(..) = _, %v, wantError=%v", err, tc.noAuth)
			}

			if uo.User.EmailCanonicalized() != email {
				t.Errorf("user.Email() = %q, want %q", uo.User.EmailCanonicalized(), "testorg")
			}

			if uo.Org.GetID() != "testorg" {
				t.Errorf("identity.OrgFromRequest(..) = %q, want %q", uo.Org, "testorg")
			}

			if len(req.Cookies()) > 2 {
				t.Errorf("ContextToRequest(..) = %q, want max 2 cookies", req.Cookies())
			}

			// Check if OrgIDHeader is propagated
			if wantVals := tc.meta.Get(org.OrgIDHeader); len(wantVals) > 0 {
				if got := req.Header.Get(org.OrgIDHeader); got != wantVals[0] {
					t.Errorf("Header %q = %q, want %q", org.OrgIDHeader, got, wantVals[0])
				}
			}
		})
	}
}

func TestOrgToRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	OrgToRequest(r, "wrongorg") // cookie set here should get overwritten
	OrgToRequest(r, "testorg")
	if len(r.Cookies()) != 1 {
		t.Errorf("OrgToRequest(..) = %q, want 1 cookie", r.Cookies())
	}
	if r.Cookies()[0].Name != org.OrgIDCookie {
		t.Errorf("OrgToRequest(..) = %q, want %q", r.Cookies()[0].Name, org.OrgIDCookie)
	}
	if r.Cookies()[0].Value != "testorg" {
		t.Errorf("OrgToRequest(..) = %q, want %q", r.Cookies()[0].Value, "testorg")
	}
	if got := r.Header.Get(org.OrgIDHeader); got != "testorg" {
		t.Errorf("OrgToRequest(..) header %q = %q, want %q", org.OrgIDHeader, got, "testorg")
	}
}

func TestOrgToContext(t *testing.T) {
	ctx := t.Context()
	ctx, err := OrgToContext(ctx, "testorg")
	if err != nil {
		t.Fatalf("OrgToContext(..) = _, %v, want no error", err)
	}
	md, _ := metadata.FromOutgoingContext(ctx)
	if md.Get(cookies.CookieHeaderName)[0] != org.OrgIDCookie+"=testorg" {
		t.Errorf("UserToContext(..) did not add the user's identity to the context")
	}
	if md.Get(org.OrgIDHeader)[0] != "testorg" {
		t.Errorf("OrgToContext(..) did not add the org header to the context")
	}
}

func TestOrgToIncomingContext(t *testing.T) {
	ctx := t.Context()
	ctx, err := OrgToIncomingContext(ctx, "testorg")
	if err != nil {
		t.Fatalf("OrgToIncomingContext(..) = _, %v, want no error", err)
	}
	md, _ := metadata.FromIncomingContext(ctx)
	if md.Get(cookies.CookieHeaderName)[0] != org.OrgIDCookie+"=testorg" {
		t.Errorf("UserToIncomingContext(..) did not add the user's identity to the context")
	}
}

func TestOrgFromRequest(t *testing.T) {
	tests := []struct {
		name    string
		cookie  string
		want    string
		wantErr bool
	}{
		{
			name:    "with cookie",
			cookie:  "testorg",
			want:    "testorg",
			wantErr: false,
		},
		{
			name:    "without cookie",
			want:    "",
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.cookie != "" {
				r.AddCookie(&http.Cookie{Name: org.OrgIDCookie, Value: tc.cookie})
			}
			o, err := OrgFromRequest(r)
			if tc.wantErr {
				if err == nil {
					t.Fatal("got err == nil, wanted error")
				}
				return
			}

			if o.GetID() != tc.want {
				t.Fatalf("org.GetID(), got %q, want %q", o.GetID(), tc.want)
			}
		})
	}
}

func TestOrgFromContext(t *testing.T) {
	tests := []struct {
		name    string
		md      metadata.MD
		want    string
		wantErr bool
	}{
		{
			name: "with cookie",
			md:   metadata.Pairs(cookies.CookieHeaderName, org.OrgIDCookie+"=testorg"),
			want: "testorg",
		},
		{
			name: "with header",
			md:   metadata.Pairs(org.OrgIDHeader, "headerorg"),
			want: "headerorg",
		},
		{
			name: "header precedence",
			md: metadata.Pairs(
				org.OrgIDHeader, "headerorg",
				cookies.CookieHeaderName, org.OrgIDCookie+"=cookieorg",
			),
			want: "headerorg",
		},
		{
			name: "multiple headers",
			md: metadata.Pairs(
				org.OrgIDHeader, "headerorg",
				org.OrgIDHeader, "headerorg2",
			),
			wantErr: true,
		},
		{
			name: "empty org header",
			md: metadata.Pairs(
				org.OrgIDHeader, "",
				cookies.CookieHeaderName, org.OrgIDCookie+"=cookieorg",
			),
			want: "cookieorg",
		},
		{
			name:    "with metadata",
			md:      metadata.Pairs(org.OrgIDCookie, "wrongorg"),
			wantErr: true,
		},
		{
			name:    "empty metadata",
			md:      metadata.Pairs(),
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(t.Context(), tc.md)
			o, err := OrgFromContext(ctx)
			if tc.wantErr {
				if err == nil {
					t.Fatal("got err == nil, wanted error")
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}
			if o.GetID() != tc.want {
				t.Fatalf("org.GetID(), got %q, want %q", o.GetID(), tc.want)
			}
		})
	}
}

func TestLoggableID(t *testing.T) {
	tests := []struct {
		desc  string
		email string
		want  string
	}{
		{
			desc:  "both_short",
			email: "a@b.uk",
			want:  "***@***",
		},
		{
			desc:  "short_prefix",
			email: "alex@ai.com",
			want:  "***@a***m",
		},
		{
			desc:  "short_domain",
			email: "henry@i.ai",
			want:  "h***y@***",
		},
		{
			desc:  "both_long",
			email: "somelongprefix@somelongdomain.com",
			want:  "s***x@s***m",
		},
		{
			desc:  "invalid",
			email: "bob",
			want:  "***@***",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			u, err := UserFromJWT(jwttesting.MustMintToken(jwttesting.WithEmail(tc.email)))
			if err != nil {
				t.Fatalf("FromJWT(%q) failed: %v", token, err)
			}
			id := u.LoggableID()
			if id != tc.want {
				t.Errorf("user.LoggableID() = %q, want %q", id, tc.want)
			}
		})
	}
}

func TestClearRequestOrg(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{Name: org.OrgIDCookie, Value: "testorg"})
	r.AddCookie(&http.Cookie{Name: "othercookie", Value: "othervalue"})
	r.Header.Set(org.OrgIDCookie, "testorg") // legacy header
	r.Header.Set(org.OrgIDHeader, "testorg")
	ClearRequestOrg(r)
	if c, err := r.Cookie(org.OrgIDCookie); err == nil {
		t.Errorf("ClearRequestOrg(..) = %q, did not clear org cookie", c)
	}
	if c, err := r.Cookie("othercookie"); err != nil {
		t.Errorf("ClearRequestOrg(..) = %q, cleared valid cookie", c)
	}
	if r.Header.Get(org.OrgIDCookie) != "" {
		t.Errorf("ClearRequestOrg(..) = %q, want empty legacy org header", r.Header.Get(org.OrgIDCookie))
	}
	if r.Header.Get(org.OrgIDHeader) != "" {
		t.Errorf("ClearRequestOrg(..) = %q, want empty org header", r.Header.Get(org.OrgIDHeader))
	}
}

func TestClearRequestUser(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	for _, name := range cookieHeaders {
		r.AddCookie(&http.Cookie{Name: name, Value: "testvalue"})
	}
	r.AddCookie(&http.Cookie{Name: "othercookie", Value: "othervalue"})
	r.Header.Set(authHeaderName, "testuser")
	r.Header.Set(ApikeyTokenHeaderName, "testuser")
	ClearRequestUser(r)
	for _, name := range cookieHeaders {
		if c, err := r.Cookie(name); err == nil {
			t.Errorf("ClearRequestUser(..) = %q, did not clear cookie %v", c, name)
		}
	}
	if c, err := r.Cookie("othercookie"); err != nil {
		t.Errorf("ClearRequestUser(..) = %q, cleared valid cookie", c)
	}
	if r.Header.Get(authHeaderName) != "" {
		t.Errorf("ClearRequestUser(..) = %q, want empty auth header", r.Header.Get(authHeaderName))
	}
	if r.Header.Get(ApikeyTokenHeaderName) != "" {
		t.Errorf("ClearRequestUser(..) = %q, want empty apikey header", r.Header.Get(ApikeyTokenHeaderName))
	}
}

func TestClearContextOrg(t *testing.T) {
	// setup context with all org cookies and headers
	testCookies := []*http.Cookie{
		{Name: "othercookie", Value: "othervalue"},
		{Name: org.OrgIDCookie, Value: "testorg"},
	}
	testMD := metadata.Pairs(cookies.ToMDString(testCookies...)...)
	testMD.Set(org.OrgIDHeader, "testorg")
	testMD.Set(org.OrgIDCookie, "testorg") // Legacy header

	ctx := metadata.NewOutgoingContext(t.Context(), testMD)
	ctx, err := ClearContextOrg(ctx)
	if err != nil {
		t.Fatalf("ClearContextOrg(..) = _, %v, want no error", err)
	}

	// check cookies
	cs, err := cookies.FromContext(ctx)
	if err != nil {
		t.Fatalf("FromContext(..) = _, %v, want no error", err)
	}
	for _, c := range cs {
		if c.Name == "othercookie" {
			continue
		}
		t.Errorf("ClearContextOrg(..) = %q, did not clear cookie %v", c, c.Name)
	}

	// check headers
	md, _ := metadata.FromOutgoingContext(ctx)
	if len(md.Get(org.OrgIDHeader)) > 0 {
		t.Errorf("ClearContextOrg(..) = %q, want empty org header", md.Get(org.OrgIDHeader)[0])
	}
	if len(md.Get(org.OrgIDCookie)) > 0 {
		t.Errorf("ClearContextOrg(..) = %q, want empty org cookie header", md.Get(org.OrgIDCookie)[0])
	}
}

func TestClearContextUser(t *testing.T) {
	// setup context with all user cookies and headers
	testCookies := []*http.Cookie{
		{Name: "othercookie", Value: "othervalue"},
	}
	for _, name := range cookieHeaders {
		testCookies = append(testCookies, &http.Cookie{Name: name, Value: "testvalue"})
	}
	testMD := metadata.Pairs(cookies.ToMDString(testCookies...)...)
	testMD.Set(authHeaderName, "testuser")
	testMD.Set(ApikeyTokenHeaderName, "testuser")
	ctx := metadata.NewOutgoingContext(t.Context(), testMD)
	ctx, err := ClearContextUser(ctx)
	if err != nil {
		t.Fatalf("ClearContextUser(..) = _, %v, want no error", err)
	}
	// check cookies
	cs, err := cookies.FromContext(ctx)
	if err != nil {
		t.Fatalf("FromContext(..) = _, %v, want no error", err)
	}
	for _, c := range cs {
		if c.Name == "othercookie" {
			continue
		}
		t.Errorf("ClearContextUser(..) = %q, did not clear cookie %v", c, c.Name)
	}
	// check headers
	md, _ := metadata.FromOutgoingContext(ctx)
	if len(md.Get(authHeaderName)) > 0 {
		t.Errorf("ClearContextUser(..) = %q, want empty auth header", md.Get(authHeaderName)[0])
	}
	if len(md.Get(ApikeyTokenHeaderName)) > 0 {
		t.Errorf("ClearContextUser(..) = %q, want empty apikey header", md.Get(ApikeyTokenHeaderName)[0])
	}
}

type UserOrgRetrievalTest struct {
	name    string
	user    *User
	org     string
	wantErr bool
}

func TestUserOrgRetrieval(t *testing.T) {
	tests := []*UserOrgRetrievalTest{
		{
			name: "with user and org",
			user: &User{
				jwt: token,
			},
			org:     "testorg",
			wantErr: false,
		},
		{
			name:    "without user",
			org:     "testorg",
			user:    &User{},
			wantErr: true,
		},
		{
			name: "without org",
			user: &User{
				jwt: token,
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run("Context/"+tc.name, func(t *testing.T) {
			runUserOrgRetrievalContext(t, tc)
		})
		t.Run("Request/"+tc.name, func(t *testing.T) {
			runUserOrgRetrievalRequest(t, tc)
		})
	}
}

func runUserOrgRetrievalContext(t *testing.T, tc *UserOrgRetrievalTest) {
	ctx, err := ToIncomingContext(t.Context(), tc.user, tc.org)
	if err != nil {
		t.Fatal(err)
	}
	uo, err := UserOrgFromContext(ctx)
	if tc.wantErr {
		if err == nil {
			t.Fatal("got err == nil, wanted error")
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if uo.User.EmailCanonicalized() != tc.user.EmailCanonicalized() {
		t.Errorf("user.Email() = %q, want %q", uo.User.EmailCanonicalized(), tc.user.EmailCanonicalized())
	}
	if uo.Org.GetID() != tc.org {
		t.Errorf("identity.OrgFromContext(..) = %q, want %q", uo.Org.GetID(), tc.org)
	}
}

func runUserOrgRetrievalRequest(t *testing.T, tc *UserOrgRetrievalTest) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	ToRequest(r, tc.user, tc.org)
	uo, err := UserOrgFromRequest(r)
	if tc.wantErr {
		if err == nil {
			t.Fatal("got err == nil, wanted error")
		}
		return
	}
	if err != nil {
		t.Fatal(err)
	}
	if uo.User.EmailCanonicalized() != tc.user.EmailCanonicalized() {
		t.Errorf("user.Email() = %q, want %q", uo.User.EmailCanonicalized(), tc.user.EmailCanonicalized())
	}
	if uo.Org.GetID() != tc.org {
		t.Errorf("identity.OrgFromRequest(..) = %q, want %q", uo.Org.GetID(), tc.org)
	}
}

func TestErrStatusCodes(t *testing.T) {
	tests := []struct {
		name     string
		exec     func() error
		wantHTTP int
		wantGRPC codes.Code
	}{
		{
			name:     "missing org id",
			exec:     func() error { return ErrMissingOrgID },
			wantHTTP: http.StatusBadRequest,
			wantGRPC: codes.InvalidArgument,
		},
		{
			name:     "unauthenticated",
			exec:     func() error { return ErrUnauthenticated },
			wantHTTP: http.StatusUnauthorized,
			wantGRPC: codes.Unauthenticated,
		},
		{
			name:     "unknown",
			exec:     func() error { return errors.New("unknown error") },
			wantHTTP: http.StatusInternalServerError,
			wantGRPC: codes.Unknown,
		},
		{
			name: "UserFromContext no metadata",
			exec: func() error {
				_, err := UserFromContext(t.Context())
				return err
			},
			wantHTTP: http.StatusUnauthorized,
			wantGRPC: codes.Unauthenticated,
		},
		{
			name: "OrgFromContext no metadata",
			exec: func() error {
				_, err := OrgFromContext(t.Context())
				return err
			},
			wantHTTP: http.StatusBadRequest,
			wantGRPC: codes.InvalidArgument,
		},
	}

	for _, tc := range tests {
		t.Run("HTTP/"+tc.name, func(t *testing.T) {
			hr := httptest.NewRecorder()
			ErrHTTP(hr, tc.exec())
			if hr.Code != tc.wantHTTP {
				t.Errorf("http.StatusUnauthorized = %v, want %v", hr.Code, tc.wantHTTP)
			}
		})
		t.Run("GRPC/"+tc.name, func(t *testing.T) {
			c := status.Code(ErrGRPC(tc.exec()))
			if c != tc.wantGRPC {
				t.Errorf("status.Code(err) = %v, want %v", c, tc.wantGRPC)
			}
		})
	}
}
