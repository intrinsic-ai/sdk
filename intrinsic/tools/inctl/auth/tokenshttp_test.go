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

package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

const (
	testAPIKey = "test-api-key"
	wantToken  = "test-token"
	wantDomain = "flowstate-dev.intrinsic.ai"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func TestNewTokensServiceClient(t *testing.T) {
	var fakeTokensService RoundTripFunc = func(req *http.Request) *http.Response {
		if req.URL.Host != wantDomain {
			t.Errorf("Unexpected request host: %q", req.URL.Host)
		}
		if req.URL.Path != "/api/v1/accountstokens:idtoken" {
			t.Errorf("Unexpected request URL: %q", req.URL.Path)
		}
		if req.Method != http.MethodPost {
			t.Errorf("Unexpected request method: %q", req.Method)
		}
		if req.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Unexpected request content type: %q", req.Header.Get("Content-Type"))
		}
		var gotReq GetIDTokenRequest
		if err := json.NewDecoder(req.Body).Decode(&gotReq); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}
		if gotReq.APIKey != testAPIKey {
			t.Errorf("Unexpected API key: %q", gotReq.APIKey)
		}
		if !gotReq.DoFanOut {
			t.Errorf("DoFanOut is false")
		}
		rb := fmt.Sprintf("{\"idToken\": %q}", wantToken)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(rb)),
		}
	}
	htcl := http.Client{Transport: fakeTokensService}
	cl, err := NewTokensServiceClient(&htcl, wantDomain)
	if err != nil {
		t.Fatal(err)
	}
	tk, err := cl.Token(context.Background(), testAPIKey)
	if err != nil {
		t.Fatal(err)
	}
	if tk != wantToken {
		t.Errorf("Token: got %q, want %q", tk, wantToken)
	}
}
