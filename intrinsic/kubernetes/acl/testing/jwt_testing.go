// Copyright 2023 Intrinsic Innovation LLC

// Package jwttesting is a test-only helper for creating JWT tokens.
package jwttesting

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/golang/glog"
	"github.com/pborman/uuid"
)

type customClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	UserID        string `json:"user_id,omitempty"`
	UID           string `json:"uid"`
	// intrinsic custom claims
	Authorized bool     `json:"authorized,omitempty"`
	Projects   []string `json:"ps,omitempty"`
	// standard claims
	jwt.StandardClaims
}

type config struct {
	signingKey []byte
	claims     customClaims
}

func defaultConfig() *config {
	return &config{
		signingKey: []byte(uuid.New()),
		claims: customClaims{
			StandardClaims: jwt.StandardClaims{},
		},
	}
}

// Option configures the JWT token creation.
type Option func(*config)

// Options is a list of [Option]s.
type Options []Option

// WithSigningKey configures a custom JWT signing key.
func WithSigningKey(k string) Option {
	return func(c *config) {
		c.signingKey = []byte(k)
	}
}

// WithEmail sets the custom "email" field in the JWT payload.
func WithEmail(e string) Option {
	return func(c *config) {
		c.claims.Email = e
	}
}

// WithEmailVerified sets the custom "email_verified" field in the JWT payload.
func WithEmailVerified(ev bool) Option {
	return func(c *config) {
		c.claims.EmailVerified = ev
	}
}

// WithUserID sets the Firebase "user_id" field in the JWT payload.
func WithUserID(u string) Option {
	return func(c *config) {
		c.claims.UserID = u
	}
}

// WithUID sets the custom "uid" field in the JWT payload.
func WithUID(u string) Option {
	return func(c *config) {
		c.claims.UID = u
	}
}

// WithExpiresAt sets the standard "exp" field in the JWT payload.
func WithExpiresAt(t time.Time) Option {
	return func(c *config) {
		c.claims.StandardClaims.ExpiresAt = t.Unix()
	}
}

// WithIssuedAt sets the standard "IssuedAt"/"iat" field in the JWT payload.
func WithIssuedAt(t time.Time) Option {
	return func(c *config) {
		c.claims.StandardClaims.IssuedAt = t.Unix()
	}
}

// WithIssuer sets the standard "Issuer"/"iss" field in the JWT payload.
func WithIssuer(i string) Option {
	return func(c *config) {
		c.claims.StandardClaims.Issuer = i
	}
}

// WithAudience sets the standard "Audience"/"aud" field in the JWT payload.
func WithAudience(a string) Option {
	return func(c *config) {
		c.claims.StandardClaims.Audience = a
	}
}

// WithSubject sets the standard "Subject"/"sub" field in the JWT payload.
func WithSubject(s string) Option {
	return func(c *config) {
		c.claims.StandardClaims.Subject = s
	}
}

// WithAuthorized sets the custom "authorized" field in the JWT payload.
func WithAuthorized(a bool) Option {
	return func(c *config) {
		c.claims.Authorized = a
	}
}

// WithProjects sets the custom "ps" field in the JWT payload.
func WithProjects(ps []string) Option {
	return func(c *config) {
		c.claims.Projects = ps
	}
}

func impl(opts ...Option) (string, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cfg.claims)
	return token.SignedString(cfg.signingKey)
}

// MintToken creates a new signed JWT token with the given options.
func MintToken(t *testing.T, opts ...Option) string {
	t.Helper()

	ss, err := impl(opts...)
	if err != nil {
		t.Fatalf("could not create a signed JWT: %v", err)
	}
	return ss
}

// MustMintToken creates a new signed JWT token with the given options. Panics if the token cannot
// be created.
func MustMintToken(opts ...Option) string {
	ss, err := impl(opts...)
	if err != nil {
		log.Exitf("Could not create a signed JWT: %v", err)
	}
	return ss
}
