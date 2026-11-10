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

package cachecontrol_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"intrinsic/httpjson/cachecontrol"
)

func assertNoCacheHeaders(t *testing.T, headers http.Header) {
	t.Helper()
	if got, want := headers.Get("Cache-Control"), "no-store, no-cache"; got != want {
		t.Errorf("Cache-Control header = %q, want %q", got, want)
	}
	if got, want := headers.Get("Pragma"), "no-cache"; got != want {
		t.Errorf("Pragma header = %q, want %q", got, want)
	}
}

func TestNoCacheHeaders(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok")
	})
	mux.HandleFunc("POST /submit", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, "created")
	})
	mux.HandleFunc("GET /error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})

	handler := cachecontrol.NoCacheHeaders(http.NewCrossOriginProtection().Handler(mux))
	server := httptest.NewServer(handler)
	defer server.Close()

	tests := []struct {
		name       string
		method     string
		path       string
		headers    map[string]string
		wantStatus int
	}{
		{
			name:       "GET 200 OK",
			method:     http.MethodGet,
			path:       "/ok",
			wantStatus: http.StatusOK,
		},
		{
			name:       "POST 201 Created",
			method:     http.MethodPost,
			path:       "/submit",
			wantStatus: http.StatusCreated,
		},
		{
			name:       "GET 500 Internal Server Error",
			method:     http.MethodGet,
			path:       "/error",
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "GET 404 Not Found",
			method:     http.MethodGet,
			path:       "/missing",
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "Cross-origin POST rejected with 403 Forbidden",
			method: http.MethodPost,
			path:   "/submit",
			headers: map[string]string{
				"Origin":         "https://evil.example",
				"Sec-Fetch-Site": "cross-site",
			},
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, server.URL+tc.path, strings.NewReader("{}"))
			if err != nil {
				t.Fatalf("http.NewRequest() failed: %v", err)
			}
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}

			resp, err := server.Client().Do(req)
			if err != nil {
				t.Fatalf("client.Do() failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tc.wantStatus)
			}
			assertNoCacheHeaders(t, resp.Header)
		})
	}
}
