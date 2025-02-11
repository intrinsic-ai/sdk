// Copyright 2023 Intrinsic Innovation LLC

// Package cookies provides shared utility functions to deal with context and header embedded information.
package cookies

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	// CookieHeaderName is the name of the header / metadata field used for cookies
	CookieHeaderName = "Cookie"
)

// parseCookies parses cookies from a string into a list of http.Cookie objects.
// Follows closely the semantics for HTTP requests.
func parseCookies(cookies string) ([]*http.Cookie, error) {
	if cookies == "" {
		return nil, nil
	}
	r, err := http.NewRequest("GET", "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %v", err)
	}
	r.Header.Set(CookieHeaderName, cookies)
	return r.Cookies(), nil
}

// FromRequestNamed returns the named cookies from a request.
// Returns only one cookie per name, ignores names that are not found.
func FromRequestNamed(r *http.Request, names []string) []*http.Cookie {
	var cs []*http.Cookie
	for _, name := range names {
		if cookie, err := r.Cookie(name); err == nil {
			cs = append(cs, cookie)
		}
	}
	return cs
}

// ToMDString converts a list of http.Cookie objects to a string that can be used as a metadata
// value.
func ToMDString(cs ...*http.Cookie) []string {
	var cookiesKV []string
	for _, c := range cs {
		cookiesKV = append(cookiesKV, (&http.Cookie{Name: c.Name, Value: c.Value}).String())
	}
	return []string{CookieHeaderName, strings.Join(cookiesKV, "; ")}
}

// FromContext extracts the "Cookies" from a GRPC incoming context.
// Cookie here refers to a mapped metadata that mirrors http cookies and is used to unify handling
// of http and GRPC based metadata in our stack.
func FromContext(ctx context.Context) ([]*http.Cookie, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	// If there's no context, we have an empty list.
	if !ok {
		return nil, nil
	}

	cookies := md.Get(CookieHeaderName)
	// If there's no cookies set, it's an empty list.
	if len(cookies) == 0 {
		return nil, nil
	}

	// If there's more than one cookie header, reject them.
	// This runs the danger of confusion when different services use first/last/arbitrary.
	if len(cookies) > 1 {
		return nil, errors.New("multiple cookies in header")
	}

	return parseCookies(cookies[0])
}
