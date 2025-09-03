// Copyright 2023 Intrinsic Innovation LLC

package auth

import (
	"context"
	"fmt"
	"net/http"

	"intrinsic/kubernetes/acl/cookies"
)

// AddMetadata contains additional metadata to be added to the request.
// Example:
//
//	md := &AddMetadata{
//	  metadata: map[string]string{"custom-header": "something"},
//	  cookies:  map[string]string{"org-id": "intrinsic-dev"},
//	}
//	NewAPIKeyTokenSource("api-key", tp, WithAdditionalMetadata(md))
type AddMetadata struct {
	metadata map[string]string
	cookies  map[string]string
}

type tokenSource interface {
	Token(ctx context.Context) (string, error)
}

// perRPCCreds is a gRPC credentials.PerRPCCredentials that adds the given token source and
// metadata to the request.
type perRPCCreds struct {
	ts            tokenSource
	md            *AddMetadata
	allowInsecure bool
}

// GetRequestMetadata returns request metadata that authenticates the request
// using a JWT retrieved using the API key.
func (s *perRPCCreds) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	t, err := s.ts.Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get account token: %v", err)
	}
	cks := []*http.Cookie{
		&http.Cookie{Name: "auth-proxy", Value: t},
	}
	// add additional cookies if provided
	for k, v := range s.md.cookies {
		if k == "" || v == "" {
			continue
		}
		cks = append(cks, &http.Cookie{Name: k, Value: v})
	}
	mdkv := cookies.ToMDString(cks...)
	metadata := map[string]string{mdkv[0]: mdkv[1]}
	// add additional metadata if provided
	for k, v := range s.md.metadata {
		if k == "" || v == "" {
			continue
		}
		metadata[k] = v
	}
	return metadata, nil
}

// RequireTransportSecurity returns the configured level of transport security.
// A token source requires transport security unless it was explicitly
// configured using [WithAllowInsecure].
func (s *perRPCCreds) RequireTransportSecurity() bool {
	return !s.allowInsecure
}
