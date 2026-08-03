// Copyright 2023 Intrinsic Innovation LLC

package origin

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestExtractOriginalURL tests ExtractOriginalURL directly, verifying header selection
// precedence (X-Original-Url vs X-Envoy-Original-Path) and ambiguity handling when both
// headers are present, without testing any secondary fallbacks.
func TestExtractOriginalURL(t *testing.T) {
	tests := []struct {
		name               string
		xOriginalURL       string
		xEnvoyOriginalPath string
		wantURL            string
	}{
		{
			name:               "X-Original-Url set",
			xOriginalURL:       "/foo/bar",
			xEnvoyOriginalPath: "",
			wantURL:            "/foo/bar",
		},
		{
			name:               "X-Envoy-Original-Path set",
			xOriginalURL:       "",
			xEnvoyOriginalPath: "/envoy/path",
			wantURL:            "/envoy/path",
		},
		{
			name:               "Both headers set -> ambiguous, returns empty string",
			xOriginalURL:       "/foo/bar",
			xEnvoyOriginalPath: "/envoy/path",
			wantURL:            "",
		},
		{
			name:               "Neither set",
			xOriginalURL:       "",
			xEnvoyOriginalPath: "",
			wantURL:            "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
			if tt.xOriginalURL != "" {
				req.Header.Set("X-Original-Url", tt.xOriginalURL)
			}
			if tt.xEnvoyOriginalPath != "" {
				req.Header.Set("X-Envoy-Original-Path", tt.xEnvoyOriginalPath)
			}

			got := ExtractOriginalURL(req)
			if got != tt.wantURL {
				t.Errorf("ExtractOriginalURL() = %q, want %q", got, tt.wantURL)
			}
		})
	}
}

// TestURL tests the high-level URL entrypoint, verifying that it uses ExtractOriginalURL
// when available and correctly falls back to X-Forwarded-Host and X-Original-Uri when
// ExtractOriginalURL returns an empty string (e.g. when headers are ambiguous or missing).
func TestURL(t *testing.T) {
	tests := []struct {
		name               string
		xOriginalURL       string
		xEnvoyOriginalPath string
		xForwardedHost     string
		xOriginalURI       string
		wantURL            string
	}{
		{
			name:         "ExtractOriginalURL priority",
			xOriginalURL: "/foo/bar",
			wantURL:      "/foo/bar",
		},
		{
			name:               "ExtractOriginalURL envoy priority",
			xEnvoyOriginalPath: "/envoy/path",
			wantURL:            "/envoy/path",
		},
		{
			name:               "Fallback to host and uri when dual headers ambiguous",
			xOriginalURL:       "/foo/bar",
			xEnvoyOriginalPath: "/envoy/path",
			xForwardedHost:     "example.com",
			xOriginalURI:       "/fallback",
			wantURL:            "https://example.com/fallback",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)
			if tt.xOriginalURL != "" {
				req.Header.Set("X-Original-Url", tt.xOriginalURL)
			}
			if tt.xEnvoyOriginalPath != "" {
				req.Header.Set("X-Envoy-Original-Path", tt.xEnvoyOriginalPath)
			}
			if tt.xForwardedHost != "" {
				req.Header.Set("X-Forwarded-Host", tt.xForwardedHost)
			}
			if tt.xOriginalURI != "" {
				req.Header.Set("X-Original-Uri", tt.xOriginalURI)
			}

			got := URL(req)
			if got != tt.wantURL {
				t.Errorf("URL() = %q, want %q", got, tt.wantURL)
			}
		})
	}
}
