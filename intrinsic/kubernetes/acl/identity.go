// Copyright 2023 Intrinsic Innovation LLC

// Package identity provides helpers to work with user identities inside the Intrinsic stack.
package identity

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"intrinsic/kubernetes/acl/cookies"
	"intrinsic/kubernetes/acl/headers"
	"intrinsic/kubernetes/acl/jwt"
	"intrinsic/kubernetes/acl/org"

	log "github.com/golang/glog"
	"go.opencensus.io/trace"
)

var (
	// ErrUnauthenticated indicates that the request was not authenticated.
	ErrUnauthenticated = errors.New("unauthenticated")
	// ErrMissingOrgID indicates that the there was no org-id found.
	ErrMissingOrgID = errors.New("no org-id found")
	// ErrInvalidRequest indicates that the request is invalid.
	ErrInvalidRequest = errors.New("invalid request")
)

// The following vars are implementation details and should not be used by a consumer of this lib.
var (
	// errEmailInvalid indicates that the email address is invalid.
	errEmailInvalid = errors.New("email address is invalid")
	// errJWTNotVerified indicates that the JWT verification failed.
	errJWTNotVerified = errors.New("the jwt was not able to be verified")

	// errNoJWT indicates that there was no JWT found.
	errNoJWT = errors.New("no jwt found")
	// errJWTUnmarshal indicates that the JWT could not be unmarshalled.
	errJWTUnmarshal = errors.New("the jwt could not be unmarshalled")
	// errCookiesParse indicates that the cookie string could not be parsed.
	errCookiesParse = errors.New("failed to parse the provivided cookies")
	// errNoIntrinsicCookie indicates that the request had no user auth cookie in the expected locations.
	errNoIntrinsicCookie = errors.New("request has no auth-proxy or onprem-token or portal-token cookie (try refreshing your browser)")
	// errNoMetadata indicates that the request had no metadata.
	errNoMetadata = errors.New("no metadata found")
	// errNoOrgIDCookie indicates that the request had no org-id cookie.
	errNoOrgIDCookie = errors.New("org-id cookie missing")
	// errOrdIDEmpty indicates that the org-id cookie was empty.
	errOrgIDEmpty = errors.New("org-id cookie is empty")
	// errMetadataKeyConflict indicates that multiple possible values were found in context metadata for a single key.
	errMetadataKeyConflict = errors.New("multiple possible values found in context metadata for a single key")
	// errOrgInMetadata indicates that the org-id was found as metadata key, it should be a cookie instead.
	errOrgInMetadata = errors.New("org-id found in metadata keys, use org-id cookie instead")

	emailRegex     = regexp.MustCompile(`(^(?P<prefix>[^@]+)@(?P<domain>.+)$)`)
	obfuscateRegex = regexp.MustCompile(`(^(.).*(.)$)`)
)

// User represents a user inside the Intrinsic stack.
type User struct {
	jwt string
	// Unmarshalled jwt cache.
	data *jwt.Data
}

// UserOrg is a return type used to answer combined requests.
type UserOrg struct {
	User *User
	Org  *org.Organization
}

const (
	// AuthProxyCookieName is the name of the cookie for storing the auth proxy token
	AuthProxyCookieName   = "auth-proxy"
	onpremTokenCookieName = "onprem-token"
	portalCookieName      = "portal-token"
	authHeaderName        = "authorization"
	// ApikeyTokenHeaderName is the metadata key for api key-based authorization
	ApikeyTokenHeaderName = "apikey-token"
	authProjectHeaderName = "x-intrinsic-auth-project"

	// IntrinsicIPCEmailSuffix is the email domain for IPC accounts.
	IntrinsicIPCEmailSuffix = "@ipc.intrinsic.ai"
)

var cookieHeaders = []string{AuthProxyCookieName, onpremTokenCookieName, portalCookieName}

// UserToRequest adds the user's identity to an HTTP request.
func UserToRequest(r *http.Request, u *User) {
	_, span := trace.StartSpan(r.Context(), "identity.UserToRequest")
	defer span.End()

	cookies.AddToRequest(r, &http.Cookie{Name: AuthProxyCookieName, Value: u.jwt})
}

// Email retrieves the canonalized user or service email of an identity.
// Alias for EmailCanonicalized.
// Deprecated: use EmailCanonicalized instead or better EmailRaw for non-ACL use cases.
func (i *User) Email() string {
	return i.EmailCanonicalized()
}

// EmailRaw retrieves the user's email as it is stored in the JWT (uncanonicalized).
// The JWT stores the primary email address of the user as defined by the login provider.
func (i *User) EmailRaw() string {
	m, _ := jwt.Email(i.jwt) // verified when created
	return m
}

// EmailCanonicalized retrieves the canonicalized user's email.
// Only use for ACL lookups.
// We are removing canonalization and ACLs are the only use case left for it.
func (i *User) EmailCanonicalized() string {
	m, _ := jwt.Email(i.jwt)    // verified when created
	m, _ = CanonicalizeEmail(m) // only error would be if m is not an email, can't happen here
	return m
}

// UserToContext adds the user's identity to a gRPC context.
func UserToContext(ctx context.Context, u *User) (context.Context, error) {
	return cookies.AddToContext(ctx, &http.Cookie{Name: AuthProxyCookieName, Value: u.jwt})
}

var (
	stripNonEmailChars = regexp.MustCompile(`[^a-zA-Z0-9!#$%&'*+\-/=?^_{|}~` + "`" + `]`)
	stripNonGmailChars = regexp.MustCompile(`[^a-zA-Z0-9]`)
)

// CanonicalizeEmail ensures that different valid forms of emails map to the same user account.
func CanonicalizeEmail(email string) (string, error) {
	// convert everything to lowercase (RFC 5321)
	parts := strings.Split(strings.ToLower(email), "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("%w: email address has an incorrect number of '@' character: %q", errEmailInvalid, email)
	}
	user, provider := parts[0], parts[1]
	if user == "" {
		return "", fmt.Errorf("%w: email address is missing the user part: %q", errEmailInvalid, email)
	}
	if provider == "" {
		return "", fmt.Errorf("%w: email address is missing the provider part: %q", errEmailInvalid, email)
	}

	// canonicalize provider
	if provider == "googlemail.com" {
		provider = "gmail.com"
	}

	// canonicalize user
	// cut everything starting with '+' on the part before the @ (including the '+') (RFC 5233)
	user = strings.SplitN(user, "+", 2)[0]

	// canonicalize user based on provider
	if provider == "gmail.com" {
		// replace all non gmail-supported characters before the @ sign
		user = stripNonGmailChars.ReplaceAllString(user, "")
	} else {
		// replace all non email-supported characters before the @ sign
		user = stripNonEmailChars.ReplaceAllString(user, "")
	}

	return user + "@" + provider, nil
}

// OrgToRequest adds the organization identifier to the HTTP request.
func OrgToRequest(r *http.Request, orgID string) {
	_, span := trace.StartSpan(r.Context(), "identity.OrgToRequest")
	defer span.End()
	cookies.AddToRequest(r, org.IDCookie(orgID))
	headers.AddOrgToRequest(r, orgID)
}

// OrgToContext returns a new context that has the org-id stored in its metadata.
func OrgToContext(ctx context.Context, orgID string) (context.Context, error) {
	if orgID == "" {
		log.WarningContextf(ctx, "OrgToContext: orgID is empty, returning unchanged context")
		return ctx, nil
	}
	ctx = headers.AddOrgToContext(ctx, orgID)
	return cookies.AddToContext(ctx, org.IDCookie(orgID))
}

// ToRequest adds the user and org metadata to the HTTP request.
func ToRequest(r *http.Request, u *User, orgID string) {
	_, span := trace.StartSpan(r.Context(), "identity.ToRequest")
	defer span.End()

	UserToRequest(r, u)
	OrgToRequest(r, orgID)
}

// ToContext adds the user and org metadata to the context.
func ToContext(ctx context.Context, u *User, orgID string) (context.Context, error) {
	_, span := trace.StartSpan(ctx, "identity.ToContext")
	defer span.End()

	ctx, err := UserToContext(ctx, u)
	if err != nil {
		return ctx, err
	}
	return OrgToContext(ctx, orgID)
}

// IPCEmail returns the IPC email based on its identifier.
// Example: "my-robot" -> "my-robot@ipc.intrinsic.ai"
func IPCEmail(name string) string {
	return name + IntrinsicIPCEmailSuffix
}

// UserFromJWT retrieves an Identity from a given JWT.
func UserFromJWT(t string) (*User, error) {
	_, err := jwt.Email(t)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get email from JWT: %v", ErrUnauthenticated, err)
	}
	return &User{jwt: t}, nil
}
