// Copyright 2023 Intrinsic Innovation LLC

package cookies

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc/metadata"
)

func TestCookiesFromContext(t *testing.T) {
	t.Run("no-metadata", func(t *testing.T) {
		result, err := FromContext(context.Background())
		if len(result) != 0 {
			t.Errorf("TestCookiesFromContext() = %v, want empty", result)
		}

		if err != nil {
			t.Errorf("TestCookiesFromContext() = %v, want nil", err)
		}
	})

	t.Run("no-cookie-header", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{})

		result, err := FromContext(ctx)
		if len(result) != 0 {
			t.Errorf("TestCookiesFromContext() = %v, want empty", result)
		}

		if err != nil {
			t.Errorf("TestCookiesFromContext() = %v, want nil", err)
		}
	})

	t.Run("empty cookie header", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{CookieHeaderName: ""}))

		result, err := FromContext(ctx)
		if len(result) != 0 {
			t.Errorf("TestCookiesFromContext() = %v, want empty", result)
		}

		if err != nil {
			t.Errorf("TestCookiesFromContext() = %v, want nil", err)
		}
	})

	t.Run("merge-cookie-headers", func(t *testing.T) {
		md := metadata.New(map[string]string{CookieHeaderName: "org-id=exampleorg; user-id=doe@example.com"})
		md.Append(CookieHeaderName, "org-id=exampleorg")
		ctx := metadata.NewIncomingContext(context.Background(), md)

		result, err := FromContext(ctx)
		if err != nil {
			t.Errorf("Error in TestCookiesFromContext() = %v, want no error", err)
		}
		if len(result) != 2 {
			t.Errorf("TestCookiesFromContext() = %v, want merged cookie", result)
		}
		if result[0].Name != "org-id" || result[0].Value != "exampleorg" {
			t.Errorf("TestCookiesFromContext() = %v, want merged cookie", result)
		}
		if result[1].Name != "user-id" || result[1].Value != "doe@example.com" {
			t.Errorf("TestCookiesFromContext() = %v, want merged cookie", result)
		}
	})

	t.Run("too-many-cookie-headers", func(t *testing.T) {
		md := metadata.New(map[string]string{CookieHeaderName: "org-id=exampleorg; user-id=doe@example.com"})
		md.Append(CookieHeaderName, "org-id=exampleorg; user-id=john@example.com")
		ctx := metadata.NewIncomingContext(context.Background(), md)

		result, err := FromContext(ctx)
		if len(result) != 0 {
			t.Errorf("TestCookiesFromContext() = %v, want empty", result)
		}

		// We expect an error in this case!
		if err == nil {
			t.Errorf("TestCookiesFromContext() was nil")
		}
	})

	t.Run("too-many-cookie-headers", func(t *testing.T) {
		expected := []*http.Cookie{
			&http.Cookie{Name: "one", Value: "val1"},
			&http.Cookie{Name: "two", Value: "val2"},
		}
		md := metadata.New(map[string]string{CookieHeaderName: "one=val1; two=val2"})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		result, err := FromContext(ctx)
		if err != nil {
			t.Errorf("TestCookiesFromContext() = %v, want nil", err)
		}

		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("TestCookiesFromContext() returned diff (-want +got):\n%s", diff)
		}
	})
}

func makeRequest(t *testing.T, cookies ...*http.Cookie) *http.Request {
	t.Helper()
	r, err := http.NewRequest("GET", "http://localhost", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	for _, c := range cookies {
		r.AddCookie(c)
	}
	return r
}

func TestFromRequestNamed(t *testing.T) {
	tests := []struct {
		name  string
		req   *http.Request
		names []string
		want  []*http.Cookie
	}{
		{
			name:  "no-cookie-no-name-header",
			req:   &http.Request{},
			names: []string{},
			want:  nil,
		},
		{
			name:  "no-cookie-header",
			req:   &http.Request{},
			names: []string{"hello"},
			want:  nil,
		},
		{
			name:  "cookie-name-not-in-request",
			req:   makeRequest(t, &http.Cookie{Name: "one", Value: "val1"}),
			names: []string{"some-name"},
			want:  nil,
		},
		{
			name:  "cookie-name-match-request",
			req:   makeRequest(t, &http.Cookie{Name: "one", Value: "val1"}),
			names: []string{"one"},
			want:  []*http.Cookie{&http.Cookie{Name: "one", Value: "val1"}},
		},
		{
			name:  "two-cookies-one-name",
			req:   makeRequest(t, &http.Cookie{Name: "one", Value: "val1"}, &http.Cookie{Name: "two", Value: "val2"}),
			names: []string{"one"},
			want:  []*http.Cookie{&http.Cookie{Name: "one", Value: "val1"}},
		},
		{
			name:  "one-cookies-no-name",
			req:   makeRequest(t, &http.Cookie{Name: "one", Value: "val1"}),
			names: []string{},
			want:  nil,
		},
		{
			name:  "duplicate-cookies",
			req:   makeRequest(t, &http.Cookie{Name: "one", Value: "val1"}, &http.Cookie{Name: "one", Value: "val2"}),
			names: []string{"one"},
			want:  []*http.Cookie{&http.Cookie{Name: "one", Value: "val1"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := FromRequestNamed(tc.req, tc.names)
			if diff := cmp.Diff(tc.want, result, cmpopts.SortSlices(func(a, b *http.Cookie) bool { return a.Name < b.Name })); diff != "" {
				t.Errorf("FromRequestNamed() returned diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestToMDString(t *testing.T) {
	tests := []struct {
		name    string
		cookies []*http.Cookie
		want    []string
	}{
		{
			name:    "no-cookies",
			cookies: []*http.Cookie{},
			want:    []string{CookieHeaderName, ""},
		},
		{
			name: "one-cookie",
			cookies: []*http.Cookie{
				&http.Cookie{Name: "one", Value: "val1"},
			},
			want: []string{CookieHeaderName, "one=val1"},
		},
		{
			name: "two-cookies",
			cookies: []*http.Cookie{
				&http.Cookie{Name: "one", Value: "val1"},
				&http.Cookie{Name: "two", Value: "val2"},
			},
			want: []string{CookieHeaderName, "one=val1; two=val2"},
		},
		{
			name: "duplicate-cookies",
			cookies: []*http.Cookie{
				&http.Cookie{Name: "one", Value: "val1"},
				&http.Cookie{Name: "one", Value: "val1"},
			},
			want: []string{CookieHeaderName, "one=val1; one=val1"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ToMDString(tc.cookies...)
			if diff := cmp.Diff(tc.want, result, cmpopts.SortSlices(func(a, b string) bool { return a < b })); diff != "" {
				t.Errorf("ToMDString() returned diff (-want +got):\n%s", diff)
			}
		})
	}
}
