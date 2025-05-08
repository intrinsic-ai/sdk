// Copyright 2023 Intrinsic Innovation LLC

// Package origin provides utilities to report where requests came from.
// See https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Common_non-standard_request_fields
// for many of the used headers.
package origin

import (
	"net/http"
)

func fromHostAndURI(r *http.Request) string {
	var res string
	if h := r.Header.Get("X-Forwarded-Host"); h != "" {
		res = "https://" + h
	}
	if u := r.Header.Get("X-Original-Uri"); u != "" {
		res += u
	}
	if res != "" {
		return res
	}
	return "<unknown>"
}

// Description returns a human readable request origin suitable for logging.
func Description(r *http.Request) string {
	if originURL := r.Header.Get("X-Original-Url"); originURL != "" {
		return "original-url=" + originURL
	}
	return "forwarded-host/original-uri=" + fromHostAndURI(r)
}

// URL returns the address of the request origin.
func URL(r *http.Request) string {
	if originURL := r.Header.Get("X-Original-Url"); originURL != "" {
		return originURL
	}
	return fromHostAndURI(r)
}
