// Copyright 2023 Intrinsic Innovation LLC

package identity

import (
	"context"
	"fmt"
	"intrinsic/kubernetes/acl/cookies"
	"intrinsic/kubernetes/acl/org"
	"intrinsic/kubernetes/acl/testing/jwttesting"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
		r.AddCookie(&http.Cookie{Name: authProxyCookieName, Value: token})
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
		r.Header.Add(apikeyTokenHeaderName, token)
		u, err := UserFromRequest(r)
		if err != nil {
			t.Fatal(err)
		}
		if u.EmailCanonicalized() != email {
			t.Fatalf("email, got %s, want %s", u.EmailCanonicalized(), email)
		}
	})
}

func TestUserFromMetadata(t *testing.T) {
	ctx := metadata.NewIncomingContext(t.Context(),
		metadata.Pairs(cookies.CookieHeaderName, authProxyCookieName+"="+token))

	u, err := UserFromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if u.EmailCanonicalized() != email {
		t.Fatalf("email, got %s, want %s", u.EmailCanonicalized(), email)
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
	requestWithAuthProxy.AddCookie(&http.Cookie{Name: authProxyCookieName, Value: token})

	requestWithOnPrem := httptest.NewRequest(http.MethodGet, "/", nil)
	requestWithOnPrem.AddCookie(&http.Cookie{Name: onpremTokenCookieName, Value: token2})

	requestWithPortal := httptest.NewRequest(http.MethodGet, "/", nil)
	requestWithPortal.AddCookie(&http.Cookie{Name: portalCookieName, Value: token2})

	requestWithBoth := httptest.NewRequest(http.MethodGet, "/", nil)
	requestWithBoth.AddCookie(&http.Cookie{Name: authProxyCookieName, Value: token})
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
				cookie, err := tc.request.Cookie(authProxyCookieName)
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
	withAuthCookie.AddCookie(&http.Cookie{Name: authProxyCookieName, Value: token})
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
		metadata.Pairs(cookies.CookieHeaderName, authProxyCookieName+"="+token))
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
				metadata.NewOutgoingContext(
					t.Context(),
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
			incoming: metadata.NewIncomingContext(
				t.Context(),
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
			incoming: metadata.NewIncomingContext(
				t.Context(),
				metadata.Pairs(
					authProxyCookieName, "anything", // do not copy over deprecated auth-proxy
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
			incoming: metadata.NewIncomingContext(
				t.Context(),
				metadata.Pairs(authHeaderName, "some-token"),
			),
			wantMd: metadata.MD{
				authHeaderName: []string{"some-token"},
			},
			wantChanged: true,
		},
		{
			desc: "multiple incoming authorization headers are ignored",
			incoming: metadata.NewIncomingContext(
				t.Context(),
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
				metadata.NewOutgoingContext(
					t.Context(),
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
				metadata.NewOutgoingContext(
					t.Context(),
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
				metadata.NewOutgoingContext(
					t.Context(),
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
				metadata.Pairs(apikeyTokenHeaderName, "something", apikeyTokenHeaderName, "something2"),
			),
			wantError: true,
		},
		{
			desc: "incoming apikey token and outgoing collide",
			incoming: metadata.NewIncomingContext(
				metadata.NewOutgoingContext(t.Context(), metadata.Pairs(apikeyTokenHeaderName, "anything")),
				metadata.Pairs(apikeyTokenHeaderName, "something"),
			),
			wantError: true,
		},
		{
			desc: "incoming org header",
			incoming: metadata.NewIncomingContext(
				t.Context(),
				metadata.Pairs(org.OrgIDHeader, "org1"),
			),
			wantMd: metadata.MD{
				org.OrgIDHeader: []string{"org1"},
			},
			wantChanged: true,
		},
		{
			desc: "multiple incoming org headers are ignored",
			incoming: metadata.NewIncomingContext(
				t.Context(),
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
					apikeyTokenHeaderName, "something2",
					"non identity relevant header", "irrelevant-value",
				),
			),
			wantMd: metadata.MD{
				"cookie":              []string{"something=somethingelse; something2=somethingelse2"},
				authHeaderName:        []string{"some-token"},
				apikeyTokenHeaderName: []string{"something2"},
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
			incoming: metadata.NewIncomingContext(
				t.Context(),
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
			meta:   metadata.New(map[string]string{cookies.CookieHeaderName: authProxyCookieName + "=" + token + "; " + org.OrgIDCookie + "=" + "testorg"}),
		},
		{
			name:   "with-duplicated-auth-coookie",
			noAuth: false,
			meta:   metadata.New(map[string]string{cookies.CookieHeaderName: authProxyCookieName + "=" + token2 + "; " + authProxyCookieName + "=" + token + "; " + org.OrgIDCookie + "=" + "wrongorg" + "; " + org.OrgIDCookie + "=" + "testorg"}),
		},
		{
			name:   "with-apikey-meta",
			noAuth: false,
			meta: metadata.New(map[string]string{
				apikeyTokenHeaderName:    token,
				cookies.CookieHeaderName: org.OrgIDCookie + "=testorg",
			}),
		},
		{
			name:   "with-org-header",
			noAuth: false,
			meta: metadata.New(map[string]string{
				apikeyTokenHeaderName:    token,
				cookies.CookieHeaderName: org.OrgIDCookie + "=cookieorg",
				org.OrgIDHeader:          "testorg",
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

			u, err := UserFromRequest(req)
			if err != nil {
				t.Fatalf("UserFromRequest(..) failed: %v", err)
			}

			if u.EmailCanonicalized() != email {
				t.Errorf("user.Email() = %q, want %q", u.EmailCanonicalized(), email)
			}

			gotOrg := ""
			if u.Org() != nil {
				gotOrg = u.Org().GetID()
			}
			if gotOrg != "testorg" {
				t.Errorf("user.Org().GetID() = %q, want %q", gotOrg, "testorg")
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

func TestClearContextChained(t *testing.T) {
	// Setup context with both org and user cookies, and other cookies
	testCookies := []*http.Cookie{
		{Name: "othercookie", Value: "othervalue"},
		{Name: org.OrgIDCookie, Value: "testorg"},
	}
	for _, name := range cookieHeaders {
		testCookies = append(testCookies, &http.Cookie{Name: name, Value: "testvalue"})
	}
	testMD := metadata.Pairs(cookies.ToMDString(testCookies...)...)
	testMD.Set(org.OrgIDHeader, "testorg")
	testMD.Set(org.OrgIDCookie, "testorg")
	testMD.Set(authHeaderName, "testuser")
	testMD.Set(apikeyTokenHeaderName, "testuser")

	ctx := metadata.NewOutgoingContext(t.Context(), testMD)

	// Call ClearContext (which chains ClearContextOrg and ClearContextUser)
	ctx, err := ClearContext(ctx)
	if err != nil {
		t.Fatalf("ClearContext(..) = _, %v, want no error", err)
	}

	// Verify only othercookie remains
	cs, err := cookies.FromOutgoingContext(ctx)
	if err != nil {
		t.Fatalf("FromOutgoingContext(..) = _, %v, want no error", err)
	}
	foundOther := false
	for _, c := range cs {
		if c.Name == "othercookie" {
			foundOther = true
			if c.Value != "othervalue" {
				t.Errorf("expected othercookie value 'othervalue', got %q", c.Value)
			}
			continue
		}
		t.Errorf("ClearContext(..) did not clear cookie %v", c.Name)
	}
	if !foundOther {
		t.Errorf("expected othercookie to be preserved, but it was lost")
	}

	// Verify headers are cleared
	md, _ := metadata.FromOutgoingContext(ctx)
	if len(md.Get(org.OrgIDHeader)) > 0 {
		t.Errorf("ClearContext(..) = %q, want empty org header", md.Get(org.OrgIDHeader)[0])
	}
	if len(md.Get(org.OrgIDCookie)) > 0 {
		t.Errorf("ClearContext(..) = %q, want empty org cookie header", md.Get(org.OrgIDCookie)[0])
	}
	if len(md.Get(authHeaderName)) > 0 {
		t.Errorf("ClearContext(..) = %q, want empty auth header", md.Get(authHeaderName)[0])
	}
	if len(md.Get(apikeyTokenHeaderName)) > 0 {
		t.Errorf("ClearContext(..) = %q, want empty apikey header", md.Get(apikeyTokenHeaderName)[0])
	}
	if len(md.Get(cookies.CookieHeaderName)) != 1 {
		t.Errorf("expected exactly 1 Cookie header list in metadata, got %d", len(md.Get(cookies.CookieHeaderName)))
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

func TestNewOutgoingContext(t *testing.T) {
	u := &User{jwt: token, org: &org.Organization{ID: "testorg"}, project: "testproject"}
	newUser := &User{jwt: "newtoken", org: &org.Organization{ID: "neworg"}, project: "newproject"}

	tests := []struct {
		name    string
		opts    []Option
		wantMD  metadata.MD
		wantErr error
	}{
		{
			name: "with user containing org and project",
			opts: []Option{WithUser(u)},
			wantMD: metadata.MD{
				authProjectHeaderName: []string{"testproject"},
				org.OrgIDHeader:       []string{"testorg"},
				"cookie":              []string{"auth-proxy=" + token + "; org-id=testorg"},
			},
		},
		{
			name: "explicit options combined",
			opts: []Option{WithUserJWT(token), WithOrg("testorg"), WithComputeProject("testproject")},
			wantMD: metadata.MD{
				authProjectHeaderName: []string{"testproject"},
				org.OrgIDHeader:       []string{"testorg"},
				"cookie":              []string{"auth-proxy=" + token + "; org-id=testorg"},
			},
		},
		{
			name:   "WithClearUser clears credentials on empty context",
			opts:   []Option{WithClearUser()},
			wantMD: metadata.MD{},
		},
		{
			name: "WithClearUser combined with WithUser atomic swap",
			opts: []Option{WithClearUser(), WithUser(newUser)},
			wantMD: metadata.MD{
				authProjectHeaderName: []string{"newproject"},
				org.OrgIDHeader:       []string{"neworg"},
				"cookie":              []string{"auth-proxy=newtoken; org-id=neworg"},
			},
		},
		{
			name: "empty options returns context unchanged",
			opts: []Option{},
		},
		{
			name:    "empty org error",
			opts:    []Option{WithOrg("")},
			wantErr: ErrMissingOrgID,
		},
		{
			name:    "empty jwt error",
			opts:    []Option{WithUserJWT("")},
			wantErr: ErrUnauthenticated,
		},
		{
			name:    "empty project error",
			opts:    []Option{WithComputeProject("")},
			wantErr: ErrMissingProject,
		},
		{
			name:    "nil user error",
			opts:    []Option{WithUser(nil)},
			wantErr: ErrInvalidRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, err := NewOutgoingContext(t.Context(), tc.opts...)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("NewOutgoingContext() error = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewOutgoingContext() unexpected error = %v", err)
			}

			gotMD, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				if len(tc.wantMD) > 0 {
					t.Fatal("FromOutgoingContext() returned ok = false but expected some metadata")
				}
				return
			}

			if diff := cmp.Diff(tc.wantMD, gotMD, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("metadata.FromOutgoingContext() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAppendToOutgoingContext(t *testing.T) {
	u := &User{jwt: token, org: &org.Organization{ID: "testorg"}, project: "testproject"}
	newUser := &User{jwt: "newtoken", org: &org.Organization{ID: "neworg"}, project: "newproject"}

	initialCookies := []*http.Cookie{
		{Name: "othercookie", Value: "othervalue"},
		{Name: authProxyCookieName, Value: "oldtoken"},
		{Name: org.OrgIDCookie, Value: "oldorg"},
	}
	initialMD := metadata.Pairs(cookies.ToMDString(initialCookies...)...)
	initialMD.Set("x-initial-key", "initial-value")
	initialMD.Set(authHeaderName, "Bearer oldtoken")
	initialMD.Set(apikeyTokenHeaderName, "oldapikey")
	initialMD.Set(org.OrgIDHeader, "oldorg")
	initialMD.Set(authProjectHeaderName, "oldproject")

	tests := []struct {
		name      string
		initialMD metadata.MD
		opts      []Option
		wantMD    metadata.MD
		wantErr   error
	}{
		{
			name:      "with user containing org and project",
			initialMD: initialMD.Copy(),
			opts:      []Option{WithUser(u)},
			wantMD: metadata.MD{
				"x-initial-key":       []string{"initial-value"},
				authHeaderName:        []string{"Bearer oldtoken"},
				apikeyTokenHeaderName: []string{"oldapikey"},
				org.OrgIDHeader:       []string{"testorg"},
				authProjectHeaderName: []string{"testproject"},
				"cookie":              []string{"auth-proxy=" + token + "; org-id=testorg; othercookie=othervalue"},
			},
		},
		{
			name:      "explicit options combined",
			initialMD: initialMD.Copy(),
			opts:      []Option{WithUserJWT(token), WithOrg("testorg"), WithComputeProject("testproject")},
			wantMD: metadata.MD{
				"x-initial-key":       []string{"initial-value"},
				authHeaderName:        []string{"Bearer oldtoken"},
				apikeyTokenHeaderName: []string{"oldapikey"},
				org.OrgIDHeader:       []string{"testorg"},
				authProjectHeaderName: []string{"testproject"},
				"cookie":              []string{"auth-proxy=" + token + "; org-id=testorg; othercookie=othervalue"},
			},
		},
		{
			name:      "WithClearUser clears credentials but preserves unrelated config",
			initialMD: initialMD.Copy(),
			opts:      []Option{WithClearUser()},
			wantMD: metadata.MD{
				"x-initial-key": []string{"initial-value"},
				"cookie":        []string{"othercookie=othervalue"},
			},
		},
		{
			name:      "WithClearUser combined with WithUser atomic swap",
			initialMD: initialMD.Copy(),
			opts:      []Option{WithClearUser(), WithUser(newUser)},
			wantMD: metadata.MD{
				"x-initial-key":       []string{"initial-value"},
				org.OrgIDHeader:       []string{"neworg"},
				authProjectHeaderName: []string{"newproject"},
				"cookie":              []string{"auth-proxy=newtoken; org-id=neworg; othercookie=othervalue"},
			},
		},
		{
			name:      "empty options returns context unchanged",
			initialMD: initialMD.Copy(),
			opts:      []Option{},
			wantMD: metadata.MD{
				"x-initial-key":       []string{"initial-value"},
				authHeaderName:        []string{"Bearer oldtoken"},
				apikeyTokenHeaderName: []string{"oldapikey"},
				org.OrgIDHeader:       []string{"oldorg"},
				authProjectHeaderName: []string{"oldproject"},
				"cookie":              []string{"othercookie=othervalue; auth-proxy=oldtoken; org-id=oldorg"},
			},
		},
		{
			name:      "empty org error",
			initialMD: initialMD.Copy(),
			opts:      []Option{WithOrg("")},
			wantErr:   ErrMissingOrgID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := metadata.NewOutgoingContext(t.Context(), tc.initialMD)

			ctx, err := AppendToOutgoingContext(ctx, tc.opts...)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("AppendToOutgoingContext() error = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("AppendToOutgoingContext() unexpected error = %v", err)
			}

			gotMD, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				if len(tc.wantMD) > 0 {
					t.Fatal("FromOutgoingContext() returned ok = false but expected some metadata")
				}
				return
			}

			if diff := cmp.Diff(tc.wantMD, gotMD, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("metadata.FromOutgoingContext() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestToRequest(t *testing.T) {
	u := &User{jwt: token, org: &org.Organization{ID: "testorg"}, project: "testproject"}

	tests := []struct {
		name        string
		opts        []Option
		wantHeaders map[string]string
		wantCookies map[string]string
		wantErr     error
	}{
		{
			name: "with user",
			opts: []Option{WithUser(u)},
			wantHeaders: map[string]string{
				authProjectHeaderName: "testproject",
				org.OrgIDHeader:       "testorg",
			},
			wantCookies: map[string]string{
				authProxyCookieName: token,
				org.OrgIDCookie:     "testorg",
			},
		},
		{
			name:    "empty org error",
			opts:    []Option{WithOrg("")},
			wantErr: ErrMissingOrgID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			err := ToRequest(r, tc.opts...)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("ToRequest() error = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("ToRequest() unexpected error = %v", err)
			}

			for k, v := range tc.wantHeaders {
				if got := r.Header.Get(k); got != v {
					t.Errorf("Header %q = %q, want %q", k, got, v)
				}
			}

			for k, v := range tc.wantCookies {
				c, err := r.Cookie(k)
				if err != nil {
					t.Errorf("Cookie %q missing: %v", k, err)
					continue
				}
				if c.Value != v {
					t.Errorf("Cookie %q = %q, want %q", k, c.Value, v)
				}
			}
		})
	}
}

func TestAppendToIncomingContext(t *testing.T) {
	u := &User{jwt: token, org: &org.Organization{ID: "testorg"}, project: "testproject"}
	newUser := &User{jwt: "newtoken", org: &org.Organization{ID: "neworg"}, project: "newproject"}

	initialCookies := []*http.Cookie{
		{Name: "othercookie", Value: "othervalue"},
		{Name: authProxyCookieName, Value: "oldtoken"},
		{Name: org.OrgIDCookie, Value: "oldorg"},
	}
	initialMD := metadata.Pairs(cookies.ToMDString(initialCookies...)...)
	initialMD.Set("x-initial-key", "initial-value")
	initialMD.Set(authHeaderName, "Bearer oldtoken")
	initialMD.Set(apikeyTokenHeaderName, "oldapikey")
	initialMD.Set(org.OrgIDHeader, "oldorg")
	initialMD.Set(authProjectHeaderName, "oldproject")

	tests := []struct {
		name      string
		initialMD metadata.MD
		opts      []Option
		wantMD    metadata.MD
	}{
		{
			name:      "append user containing org and project to incoming context",
			initialMD: initialMD.Copy(),
			opts:      []Option{WithUser(u)},
			wantMD: metadata.MD{
				"x-initial-key":       []string{"initial-value"},
				authHeaderName:        []string{"Bearer oldtoken"},
				apikeyTokenHeaderName: []string{"oldapikey"},
				org.OrgIDHeader:       []string{"testorg"},
				authProjectHeaderName: []string{"testproject"},
				"cookie":              []string{"auth-proxy=" + token + "; org-id=testorg; othercookie=othervalue"},
			},
		},
		{
			name:      "explicit options combined",
			initialMD: initialMD.Copy(),
			opts:      []Option{WithUserJWT(token), WithOrg("testorg"), WithComputeProject("testproject")},
			wantMD: metadata.MD{
				"x-initial-key":       []string{"initial-value"},
				authHeaderName:        []string{"Bearer oldtoken"},
				apikeyTokenHeaderName: []string{"oldapikey"},
				org.OrgIDHeader:       []string{"testorg"},
				authProjectHeaderName: []string{"testproject"},
				"cookie":              []string{"auth-proxy=" + token + "; org-id=testorg; othercookie=othervalue"},
			},
		},
		{
			name:      "WithClearUser clears credentials but preserves system metadata and other cookies",
			initialMD: initialMD.Copy(),
			opts:      []Option{WithClearUser()},
			wantMD: metadata.MD{
				"x-initial-key": []string{"initial-value"},
				"cookie":        []string{"othercookie=othervalue"},
			},
		},
		{
			name:      "WithClearUser combined with WithUser - atomic swap",
			initialMD: initialMD.Copy(),
			opts:      []Option{WithClearUser(), WithUser(newUser)},
			wantMD: metadata.MD{
				"x-initial-key":       []string{"initial-value"},
				org.OrgIDHeader:       []string{"neworg"},
				authProjectHeaderName: []string{"newproject"},
				"cookie":              []string{"auth-proxy=newtoken; org-id=neworg; othercookie=othervalue"},
			},
		},
		{
			name:      "empty options returns context unchanged",
			initialMD: initialMD.Copy(),
			opts:      []Option{},
			wantMD: metadata.MD{
				"x-initial-key":       []string{"initial-value"},
				authHeaderName:        []string{"Bearer oldtoken"},
				apikeyTokenHeaderName: []string{"oldapikey"},
				org.OrgIDHeader:       []string{"oldorg"},
				authProjectHeaderName: []string{"oldproject"},
				"cookie":              []string{"othercookie=othervalue; auth-proxy=oldtoken; org-id=oldorg"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := metadata.NewIncomingContext(t.Context(), tc.initialMD)

			ctx, err := AppendToIncomingContext(ctx, tc.opts...)
			if err != nil {
				t.Fatalf("AppendToIncomingContext() error = %v", err)
			}

			gotMD, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				if len(tc.wantMD) > 0 {
					t.Fatal("FromIncomingContext() returned ok = false but expected some metadata")
				}
				return
			}

			if diff := cmp.Diff(tc.wantMD, gotMD, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("metadata.FromIncomingContext() returned unexpected diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestToMetadata(t *testing.T) {
	token := jwttesting.MintToken(t, jwttesting.WithEmail("test@google.com"))
	u := &User{jwt: token, org: &org.Organization{ID: "testorg"}, project: "testproject"}
	newUser := &User{jwt: "newtoken", org: &org.Organization{ID: "neworg"}, project: "newproject"}

	tests := []struct {
		name       string
		opts       []Option
		wantMD     metadata.MD
		wantErrSub string
	}{
		{
			name: "with user containing org and project",
			opts: []Option{WithUser(u)},
			wantMD: metadata.MD{
				authProjectHeaderName: []string{"testproject"},
				org.OrgIDHeader:       []string{"testorg"},
				"cookie":              []string{"auth-proxy=" + token + "; org-id=testorg"},
			},
		},
		{
			name: "explicit options combined",
			opts: []Option{WithUserJWT(token), WithOrg("testorg"), WithComputeProject("testproject")},
			wantMD: metadata.MD{
				authProjectHeaderName: []string{"testproject"},
				org.OrgIDHeader:       []string{"testorg"},
				"cookie":              []string{"auth-proxy=" + token + "; org-id=testorg"},
			},
		},
		{
			name:   "WithClearUser clears credentials on empty context",
			opts:   []Option{WithClearUser()},
			wantMD: metadata.MD{},
		},
		{
			name: "WithClearUser combined with WithUser atomic swap",
			opts: []Option{WithClearUser(), WithUser(newUser)},
			wantMD: metadata.MD{
				authProjectHeaderName: []string{"newproject"},
				org.OrgIDHeader:       []string{"neworg"},
				"cookie":              []string{"auth-proxy=newtoken; org-id=neworg"},
			},
		},
		{
			name:   "empty options returns empty metadata",
			opts:   []Option{},
			wantMD: metadata.MD{},
		},
		{
			name:       "option evaluation error propagated",
			opts:       []Option{WithOrg("")},
			wantErrSub: "no org-id found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotMD, err := ToMetadata(tc.opts...)
			if tc.wantErrSub != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErrSub) {
					t.Fatalf("ToMetadata() returned error = %v, want error containing %q", err, tc.wantErrSub)
				}
				return
			}
			if err != nil {
				t.Fatalf("ToMetadata() unexpected error = %v", err)
			}
			if diff := cmp.Diff(tc.wantMD, gotMD, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ToMetadata() returned unexpected metadata diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAppendToMetadata(t *testing.T) {
	token := jwttesting.MintToken(t, jwttesting.WithEmail("test@google.com"))
	u := &User{jwt: token, org: &org.Organization{ID: "testorg"}, project: "testproject"}
	newUser := &User{jwt: "newtoken", org: &org.Organization{ID: "neworg"}, project: "newproject"}

	initialCookies := []*http.Cookie{
		{Name: "othercookie", Value: "othervalue"},
		{Name: authProxyCookieName, Value: "oldtoken"},
		{Name: org.OrgIDCookie, Value: "oldorg"},
	}
	initialMD := metadata.Pairs(cookies.ToMDString(initialCookies...)...)
	initialMD.Set("x-initial-key", "initial-value")
	initialMD.Set(authHeaderName, "Bearer oldtoken")
	initialMD.Set(apikeyTokenHeaderName, "oldapikey")
	initialMD.Set(org.OrgIDHeader, "oldorg")
	initialMD.Set(authProjectHeaderName, "oldproject")

	tests := []struct {
		name       string
		initialMD  metadata.MD
		opts       []Option
		wantMD     metadata.MD
		wantErrSub string
	}{
		{
			name:      "with user containing org and project",
			initialMD: initialMD.Copy(),
			opts:      []Option{WithUser(u)},
			wantMD: metadata.MD{
				"x-initial-key":       []string{"initial-value"},
				authHeaderName:        []string{"Bearer oldtoken"},
				apikeyTokenHeaderName: []string{"oldapikey"},
				org.OrgIDHeader:       []string{"testorg"},
				authProjectHeaderName: []string{"testproject"},
				"cookie":              []string{"auth-proxy=" + token + "; org-id=testorg; othercookie=othervalue"},
			},
		},
		{
			name:      "explicit options combined",
			initialMD: initialMD.Copy(),
			opts:      []Option{WithUserJWT(token), WithOrg("testorg"), WithComputeProject("testproject")},
			wantMD: metadata.MD{
				"x-initial-key":       []string{"initial-value"},
				authHeaderName:        []string{"Bearer oldtoken"},
				apikeyTokenHeaderName: []string{"oldapikey"},
				org.OrgIDHeader:       []string{"testorg"},
				authProjectHeaderName: []string{"testproject"},
				"cookie":              []string{"auth-proxy=" + token + "; org-id=testorg; othercookie=othervalue"},
			},
		},
		{
			name:      "WithClearUser clears credentials but preserves unrelated config",
			initialMD: initialMD.Copy(),
			opts:      []Option{WithClearUser()},
			wantMD: metadata.MD{
				"x-initial-key": []string{"initial-value"},
				"cookie":        []string{"othercookie=othervalue"},
			},
		},
		{
			name:      "WithClearUser combined with WithUser atomic swap",
			initialMD: initialMD.Copy(),
			opts:      []Option{WithClearUser(), WithUser(newUser)},
			wantMD: metadata.MD{
				"x-initial-key":       []string{"initial-value"},
				org.OrgIDHeader:       []string{"neworg"},
				authProjectHeaderName: []string{"newproject"},
				"cookie":              []string{"auth-proxy=newtoken; org-id=neworg; othercookie=othervalue"},
			},
		},
		{
			name:      "empty options returns input unchanged",
			initialMD: initialMD.Copy(),
			opts:      []Option{},
			wantMD: metadata.MD{
				"x-initial-key":       []string{"initial-value"},
				authHeaderName:        []string{"Bearer oldtoken"},
				apikeyTokenHeaderName: []string{"oldapikey"},
				org.OrgIDHeader:       []string{"oldorg"},
				authProjectHeaderName: []string{"oldproject"},
				"cookie":              []string{"othercookie=othervalue; auth-proxy=oldtoken; org-id=oldorg"},
			},
		},
		{
			name:       "option evaluation error propagated",
			initialMD:  nil,
			opts:       []Option{WithOrg("")},
			wantErrSub: "no org-id found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotMD, err := AppendToMetadata(tc.initialMD, tc.opts...)
			if tc.wantErrSub != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErrSub) {
					t.Fatalf("AppendToMetadata() returned error = %v, want error containing %q", err, tc.wantErrSub)
				}
				return
			}
			if err != nil {
				t.Fatalf("AppendToMetadata() unexpected error = %v", err)
			}
			if diff := cmp.Diff(tc.wantMD, gotMD, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("AppendToMetadata() returned unexpected metadata diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestToMetadataMap(t *testing.T) {
	token := jwttesting.MintToken(t, jwttesting.WithEmail("test@google.com"))
	u := &User{jwt: token, org: &org.Organization{ID: "testorg"}, project: "testproject"}
	newUser := &User{jwt: "newtoken", org: &org.Organization{ID: "neworg"}, project: "newproject"}

	tests := []struct {
		name       string
		opts       []Option
		wantMap    map[string]string
		wantErrSub string
	}{
		{
			name: "with user containing org and project",
			opts: []Option{WithUser(u)},
			wantMap: map[string]string{
				authProjectHeaderName: "testproject",
				org.OrgIDHeader:       "testorg",
				"cookie":              "auth-proxy=" + token + "; org-id=testorg",
			},
		},
		{
			name: "explicit options combined",
			opts: []Option{WithUserJWT(token), WithOrg("testorg"), WithComputeProject("testproject")},
			wantMap: map[string]string{
				authProjectHeaderName: "testproject",
				org.OrgIDHeader:       "testorg",
				"cookie":              "auth-proxy=" + token + "; org-id=testorg",
			},
		},
		{
			name:    "WithClearUser clears credentials on empty context",
			opts:    []Option{WithClearUser()},
			wantMap: map[string]string{},
		},
		{
			name: "WithClearUser combined with WithUser atomic swap",
			opts: []Option{WithClearUser(), WithUser(newUser)},
			wantMap: map[string]string{
				authProjectHeaderName: "newproject",
				org.OrgIDHeader:       "neworg",
				"cookie":              "auth-proxy=newtoken; org-id=neworg",
			},
		},
		{
			name:    "empty options returns empty map",
			opts:    []Option{},
			wantMap: map[string]string{},
		},
		{
			name:       "option evaluation error propagated",
			opts:       []Option{WithOrg("")},
			wantErrSub: "no org-id found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotMap, err := ToMetadataMap(tc.opts...)
			if tc.wantErrSub != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErrSub) {
					t.Fatalf("ToMetadataMap() returned error = %v, want error containing %q", err, tc.wantErrSub)
				}
				return
			}
			if err != nil {
				t.Fatalf("ToMetadataMap() unexpected error = %v", err)
			}
			if diff := cmp.Diff(tc.wantMap, gotMap, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ToMetadataMap() returned unexpected map diff (-want +got):\n%s", diff)
			}
		})
	}
}

func copyMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	res := make(map[string]string, len(m))
	for k, v := range m {
		res[k] = v
	}
	return res
}

func TestAppendToMetadataMap(t *testing.T) {
	token := jwttesting.MintToken(t, jwttesting.WithEmail("test@google.com"))
	u := &User{jwt: token, org: &org.Organization{ID: "testorg"}, project: "testproject"}
	newUser := &User{jwt: "newtoken", org: &org.Organization{ID: "neworg"}, project: "newproject"}

	initialMapWithCreds := map[string]string{
		"x-initial-key":       "initial-value",
		"cookie":              "othercookie=othervalue; auth-proxy=oldtoken; org-id=oldorg",
		"authorization":       "Bearer oldtoken",
		"apikey-token":        "oldapikey",
		org.OrgIDHeader:       "oldorg",
		authProjectHeaderName: "oldproject",
		"x-api-key":           "oldapikey",  // unrelated
		"x-intrinsic-project": "oldproject", // unrelated
	}

	tests := []struct {
		name       string
		initialMap map[string]string
		opts       []Option
		wantMap    map[string]string
		wantErrSub string
	}{
		{
			name:       "with user containing org and project",
			initialMap: copyMap(initialMapWithCreds),
			opts:       []Option{WithUser(u)},
			wantMap: map[string]string{
				"x-initial-key":       "initial-value",
				"authorization":       "Bearer oldtoken",
				"apikey-token":        "oldapikey",
				org.OrgIDHeader:       "testorg",
				authProjectHeaderName: "testproject",
				"cookie":              "auth-proxy=" + token + "; org-id=testorg; othercookie=othervalue",
				"x-api-key":           "oldapikey",
				"x-intrinsic-project": "oldproject",
			},
		},
		{
			name:       "explicit options combined",
			initialMap: copyMap(initialMapWithCreds),
			opts:       []Option{WithUserJWT(token), WithOrg("testorg"), WithComputeProject("testproject")},
			wantMap: map[string]string{
				"x-initial-key":       "initial-value",
				"authorization":       "Bearer oldtoken",
				"apikey-token":        "oldapikey",
				org.OrgIDHeader:       "testorg",
				authProjectHeaderName: "testproject",
				"cookie":              "auth-proxy=" + token + "; org-id=testorg; othercookie=othervalue",
				"x-api-key":           "oldapikey",
				"x-intrinsic-project": "oldproject",
			},
		},
		{
			name:       "WithClearUser clears credentials but preserves unrelated config",
			initialMap: copyMap(initialMapWithCreds),
			opts:       []Option{WithClearUser()},
			wantMap: map[string]string{
				"x-initial-key":       "initial-value",
				"cookie":              "othercookie=othervalue",
				"x-api-key":           "oldapikey",
				"x-intrinsic-project": "oldproject",
			},
		},
		{
			name:       "WithClearUser combined with WithUser atomic swap",
			initialMap: copyMap(initialMapWithCreds),
			opts:       []Option{WithClearUser(), WithUser(newUser)},
			wantMap: map[string]string{
				"x-initial-key":       "initial-value",
				org.OrgIDHeader:       "neworg",
				authProjectHeaderName: "newproject",
				"cookie":              "auth-proxy=newtoken; org-id=neworg; othercookie=othervalue",
				"x-api-key":           "oldapikey",
				"x-intrinsic-project": "oldproject",
			},
		},
		{
			name:       "empty options returns input unchanged",
			initialMap: copyMap(initialMapWithCreds),
			opts:       []Option{},
			wantMap: map[string]string{
				"x-initial-key":       "initial-value",
				"authorization":       "Bearer oldtoken",
				"apikey-token":        "oldapikey",
				org.OrgIDHeader:       "oldorg",
				authProjectHeaderName: "oldproject",
				"cookie":              "othercookie=othervalue; auth-proxy=oldtoken; org-id=oldorg",
				"x-api-key":           "oldapikey",
				"x-intrinsic-project": "oldproject",
			},
		},
		{
			name:       "option evaluation error propagated",
			initialMap: nil,
			opts:       []Option{WithOrg("")},
			wantErrSub: "no org-id found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotMap, err := AppendToMetadataMap(tc.initialMap, tc.opts...)
			if tc.wantErrSub != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErrSub) {
					t.Fatalf("AppendToMetadataMap() returned error = %v, want error containing %q", err, tc.wantErrSub)
				}
				return
			}
			if err != nil {
				t.Fatalf("AppendToMetadataMap() unexpected error = %v", err)
			}
			if diff := cmp.Diff(tc.wantMap, gotMap, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("AppendToMetadataMap() returned unexpected map diff (-want +got):\n%s", diff)
			}
		})
	}
}
