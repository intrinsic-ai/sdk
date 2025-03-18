// Copyright 2023 Intrinsic Innovation LLC

// Package cookies provides shared utility functions to deal with context and header embedded information.
package cookies

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	log "github.com/golang/glog"
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

	// If there's more than one cookie header, we attempt to merge them.
	if len(cookies) > 1 {
		log.WarningContextf(ctx, "Multiple cookie headers found, attempting to merge them...")
		return mergeCookies(cookies...)
	}

	return parseCookies(cookies[0])
}

func mergeCookies(cs ...string) ([]*http.Cookie, error) {
	mergedCookies := make(map[string]*http.Cookie)

	for _, cval := range cs {
		cookies, err := parseCookies(cval)
		if err != nil {
			return nil, fmt.Errorf("failed to parse cookies: %v", err)
		}

		for _, cookie := range cookies {
			mc, ok := mergedCookies[cookie.Name]
			// If the cookie already exists, check if the value is the same.
			if ok && mc.Value != cookie.Value { // unhappy path
				return nil, fmt.Errorf("conflicting cookie values for key %s", cookie.Name)
			}
			// If the cookie does not exist or the value is the same, just add/overwrite it.
			mergedCookies[cookie.Name] = cookie
		}
	}

	result := make([]*http.Cookie, 0, len(mergedCookies))
	for _, cookie := range mergedCookies {
		result = append(result, cookie)
	}

	return result, nil
}
