// Copyright 2023 Intrinsic Innovation LLC

// Package origin provides utilities to report where requests came from.
// See https://en.wikipedia.org/wiki/List_of_HTTP_header_fields#Common_non-standard_request_fields
// for many of the used headers.
package origin

import (
	"net/http"

	"github.com/golang/glog"
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
	if originPath := r.Header.Get("X-Envoy-Original-Path"); originPath != "" {
		return "envoy-original-path=" + originPath
	}
	return "forwarded-host/original-uri=" + fromHostAndURI(r)
}

// ExtractOriginalURL returns the original URL or path from proxy headers.
func ExtractOriginalURL(r *http.Request) string {
	u := r.Header.Get("X-Original-Url")
	p := r.Header.Get("X-Envoy-Original-Path")
	if u != "" && p != "" {
		glog.Warningf("Ambiguous request headers: both X-Original-Url (%q) and X-Envoy-Original-Path (%q) are set", u, p)
		return ""
	}
	if u != "" {
		return u
	}
	return p
}

// URL returns the address of the request origin.
func URL(r *http.Request) string {
	if u := ExtractOriginalURL(r); u != "" {
		return u
	}
	return fromHostAndURI(r)
}
