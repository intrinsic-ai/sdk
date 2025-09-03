// Copyright 2023 Intrinsic Innovation LLC

package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc/credentials"
	"intrinsic/kubernetes/acl/jwt"
)

const defaultMinTokenLifetime = time.Minute

var _ credentials.PerRPCCredentials = &APIKeyTokenSource{}

// timeNow can be overridden in tests.
var timeNow = time.Now

// APIKeyTokenProvider provides a token for an API key.
type APIKeyTokenProvider interface {
	Token(ctx context.Context, apiKey string) (string, error)
}

// APIKeyTokenSource provides a JWT token retrieved using an API key. Can be
// used as [credentials.PerRPCCredentials] with gRPC clients.
type APIKeyTokenSource struct {
	cts    cachedTokenSource
	perRPC perRPCCreds
}

// APIKeyTokenSourceOption configures an [APIKeyTokenSource].
type APIKeyTokenSourceOption = func(s *APIKeyTokenSource)

// WithAllowInsecure enables the token source to add credentials on insecure
// connections. This can be necessary to pass credentials over local connections
// that use insecure transport security.
func WithAllowInsecure() APIKeyTokenSourceOption {
	return func(s *APIKeyTokenSource) {
		s.perRPC.allowInsecure = true
	}
}

// WithMinTokenLifetime specifies the minimum amount of time that a token must
// still be valid at the time of the request. Defaults to 1 minute. The token
// source will retrieve a new token if the expiry is too close. Specifying a
// duration larger than the initial lifetime of the token means that it will be
// refreshed on every request.
func WithMinTokenLifetime(d time.Duration) APIKeyTokenSourceOption {
	return func(s *APIKeyTokenSource) {
		s.cts.minTokenLifetime = d
	}
}

// WithAdditionalMetadata adds additional metadata to the request.
func WithAdditionalMetadata(md *AddMetadata) APIKeyTokenSourceOption {
	return func(s *APIKeyTokenSource) {
		s.perRPC.md = md
	}
}

// NewAPIKeyTokenSource creates and configures an [APIKeyTokenSource].
func NewAPIKeyTokenSource(apiKey string, tp APIKeyTokenProvider, opts ...APIKeyTokenSourceOption) *APIKeyTokenSource {
	s := &APIKeyTokenSource{
		cts: cachedTokenSource{
			tp:               tp,
			apiKey:           apiKey,
			minTokenLifetime: defaultMinTokenLifetime,
		},
	}
	s.perRPC.ts = &s.cts
	for _, opt := range opts {
		opt(s)
	}
	if s.perRPC.md == nil {
		s.perRPC.md = &AddMetadata{}
	}
	return s
}

// GetRequestMetadata returns request metadata that authenticates the request
// using a JWT retrieved using the API key.
func (s *APIKeyTokenSource) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	return s.perRPC.GetRequestMetadata(ctx)
}

// RequireTransportSecurity returns the configured level of transport security.
// A token source requires transport security unless it was explicitly
// configured using [WithAllowInsecure].
func (s *APIKeyTokenSource) RequireTransportSecurity() bool {
	return s.perRPC.RequireTransportSecurity()
}

// Token returns a JWT token retrieved using the API key.
func (s *APIKeyTokenSource) Token(ctx context.Context) (string, error) {
	return s.cts.Token(ctx)
}

// cachedTokenSource adds caching to an APIKeyTokenProvider.
type cachedTokenSource struct {
	tp               APIKeyTokenProvider
	apiKey           string
	minTokenLifetime time.Duration

	mu sync.Mutex
	c  *tokenCache
}

type tokenCache struct {
	t      string
	expiry time.Time
}

func (s *cachedTokenSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.c == nil || s.c.expiry.Add(-s.minTokenLifetime).Before(timeNow()) {
		t, err := s.tp.Token(ctx, s.apiKey)
		if err != nil {
			return "", fmt.Errorf("could not get account token: %v", err)
		}
		d, err := jwt.UnmarshalUnsafe(t)
		if err != nil {
			return "", fmt.Errorf("could not unmarshal account token: %v", err)
		}
		s.c = &tokenCache{
			t:      t,
			expiry: time.Unix(d.ExpiresAt, 0),
		}
	}
	return s.c.t, nil
}
