// Copyright 2023 Intrinsic Innovation LLC

// Package identity provides helpers to work with user identities inside the Intrinsic stack.
package identity

import (
	"context"
	"errors"
	"fmt"
	"intrinsic/frontend/go/origin"
	"intrinsic/kubernetes/acl/cookies"
	"intrinsic/kubernetes/acl/headers"
	"intrinsic/kubernetes/acl/jwt"
	"intrinsic/kubernetes/acl/org"
	"intrinsic/stats/go/telemetry"
	"net/http"
	"regexp"
	"slices"
	"strings"

	log "github.com/golang/glog"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	// ErrUnauthenticated indicates that the request was not authenticated.
	ErrUnauthenticated = errors.New("unauthenticated")
	// ErrMissingOrgID indicates that the there was no org-id found.
	ErrMissingOrgID = errors.New("no org-id found")
	// ErrMissingProject indicates that there was no project-id found.
	ErrMissingProject = errors.New("no project-id found")
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
	// errNoProjectInMetadata indicates that there was no project-id found in metadata.
	errNoProjectInMetadata = errors.New("project-id metadata key missing")
	// errNoProjectInHeader indicates that there was no project-id found in headers.
	errNoProjectInHeader = errors.New("x-intrinsic-auth-project header missing")

	emailRegex     = regexp.MustCompile(`(^(?P<prefix>[^@]+)@(?P<domain>.+)$)`)
	obfuscateRegex = regexp.MustCompile(`(^(.).*(.)$)`)
)

// ErrGRPC converts errors from the identity package to the corresponding gRPC error.
func ErrGRPC(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, ErrMissingOrgID), errors.Is(err, ErrInvalidRequest):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}

// ErrHTTP converts errors from the identity package to the corresponding HTTP status code and writes it to the response writer.
func ErrHTTP(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	switch {
	case errors.Is(err, ErrMissingOrgID), errors.Is(err, ErrInvalidRequest):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrUnauthenticated):
		http.Error(w, err.Error(), http.StatusUnauthorized)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// User represents a user inside the Intrinsic stack.
type User struct {
	jwt string
	// Unmarshalled jwt cache.
	data *jwt.Data
	// Org is the organization of the user, populated dynamically when extracted from context/request.
	org *org.Organization
	// Project is the compute project of the user, populated dynamically when extracted from context/request.
	project string
}

// Org returns the organization of the user.
func (i *User) Org() *org.Organization {
	if i == nil {
		return nil
	}
	return i.org
}

// Project returns the compute project of the user.
func (i *User) Project() string {
	if i == nil {
		return ""
	}
	return i.project
}

// UserOrg is a return type used to answer combined requests.
type UserOrg struct {
	User *User
	Org  *org.Organization
}

const (
	// authProxyCookieName is the name of the cookie for storing the auth proxy token
	authProxyCookieName   = "auth-proxy"
	onpremTokenCookieName = "onprem-token"
	portalCookieName      = "portal-token"
	authHeaderName        = "authorization"
	// apikeyTokenHeaderName is the metadata key for api key-based authorization
	apikeyTokenHeaderName = "apikey-token"
	authProjectHeaderName = "x-intrinsic-auth-project"

	// Real user accounts have @google.com email addresses
	googleEmailSuffix = "@google.com"
	// IntrinsicServiceAccountEmailSuffix is the email domain for service accounts (such as automated processes)
	IntrinsicServiceAccountEmailSuffix = "@serviceaccount.intrinsic.ai"
	// IntrinsicTestAccountEmailSuffix is the email domain for owned test accounts (OTA, see go/rhea)
	IntrinsicTestAccountEmailSuffix = "@gmail.com"

	// IntrinsicIPCEmailSuffix is the email domain for IPC accounts.
	IntrinsicIPCEmailSuffix = "@ipc.intrinsic.ai"
)

func (i *User) populateOrgAndProject(ctx context.Context) error {
	if i == nil {
		return nil
	}
	o, err := OrgFromContext(ctx)
	if err == nil {
		i.org = o
	} else if !errors.Is(err, ErrMissingOrgID) {
		return err
	}
	p, err := extractProjectFromContext(ctx)
	if err == nil {
		i.project = p
	} else if !errors.Is(err, ErrMissingProject) {
		return err
	}
	return nil
}

func (i *User) populateOrgAndProjectFromRequest(r *http.Request) error {
	if i == nil {
		return nil
	}
	o, err := OrgFromRequest(r)
	if err == nil {
		i.org = o
	} else if !errors.Is(err, ErrMissingOrgID) {
		return err
	}
	p, err := extractProjectFromRequest(r)
	if err == nil {
		i.project = p
	} else if !errors.Is(err, ErrMissingProject) {
		return err
	}
	return nil
}

func extractProjectFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.Join(ErrMissingProject, errNoMetadata)
	}
	projectVal := md.Get(authProjectHeaderName)
	if len(projectVal) == 0 || projectVal[0] == "" {
		return "", errors.Join(ErrMissingProject, errNoProjectInMetadata)
	}
	return projectVal[0], nil
}

func extractProjectFromRequest(r *http.Request) (string, error) {
	if r == nil {
		return "", errors.Join(ErrMissingProject, ErrInvalidRequest)
	}
	projectVal := r.Header.Get(authProjectHeaderName)
	if projectVal == "" {
		return "", errors.Join(ErrMissingProject, errNoProjectInHeader)
	}
	return projectVal, nil
}

var cookieHeaders = []string{authProxyCookieName, onpremTokenCookieName, portalCookieName}

// GetJWTFromRequest returns the JWT from a request.
func GetJWTFromRequest(r *http.Request) (string, error) {
	_, span := trace.StartSpan(r.Context(), "identity.GetJWTFromRequest")
	defer span.End()

	for _, cn := range cookieHeaders {
		jwt, err := r.Cookie(cn)
		if err == nil {
			log.V(2).Infof("Using jwt from cookie %q", cn)
			return jwt.Value, nil
		}
	}
	if token := r.Header.Get(apikeyTokenHeaderName); token != "" {
		log.V(2).Infof("Using jwt from header %q", apikeyTokenHeaderName)
		return token, nil
	}
	// Retrieving the JWT from the authorization header.
	// The authorization header is set for service-to-service communication.
	if token := trimBearer(r.Header.Get(authHeaderName)); token != "" {
		log.V(2).Infof("Using jwt from header %q", authHeaderName)
		return token, nil
	}
	telemetry.SetError(span, trace.StatusCodeUnauthenticated, "GetJWTFromRequest", errNoJWT)
	return "", errors.Join(ErrUnauthenticated, errNoJWT)
}

// GetJWTFromContext returns the JWT from a context.
func GetJWTFromContext(ctx context.Context) (string, error) {
	_, span := trace.StartSpan(ctx, "identity.GetJWTFromContext")
	defer span.End()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		telemetry.SetError(span, trace.StatusCodeUnauthenticated, "GetJWTFromContext", errNoMetadata)
		return "", errors.Join(ErrUnauthenticated, errNoMetadata)
	}

	contextCookies, err := cookies.FromContext(ctx)
	if err != nil {
		telemetry.SetError(span, trace.StatusCodeInvalidArgument, "GetJWTFromContext: Failed to parse cookies", errCookiesParse)
		return "", errors.Join(ErrUnauthenticated, errCookiesParse, err)
	}
	for _, c := range contextCookies {
		for _, cn := range cookieHeaders {
			if c.Name == cn {
				return c.Value, nil
			}
		}
	}
	// Retrieving the JWT from the apikey-token header.
	if jwtMD, ok := md[apikeyTokenHeaderName]; ok && len(jwtMD) > 0 && jwtMD[0] != "" {
		return jwtMD[0], nil
	}
	// Retrieving the JWT from the authorization header.
	// The authorization header is set for service-to-service communication.
	if tks, ok := md[authHeaderName]; ok && len(tks) > 0 && trimBearer(tks[0]) != "" {
		return trimBearer(tks[0]), nil
	}
	telemetry.SetError(span, trace.StatusCodeUnauthenticated, "GetJWTFromContext", errNoJWT)
	return "", errors.Join(ErrUnauthenticated, errNoJWT)
}

// trimBearer strips the "Bearer " or "bearer " prefix from the authorization header.
// Returns empty string if no prefix is there.
func trimBearer(authHeader string) string {
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	if strings.HasPrefix(authHeader, "bearer ") {
		return strings.TrimPrefix(authHeader, "bearer ")
	}
	return ""
}

// TokenVerifier verifies the given JWT.
// Required for From{Request,Context}Verified.
type TokenVerifier interface {
	VerifyIDToken(ctx context.Context, token string) error
}

// UserFromContextVerified returns the verified user's identity from an incoming context.
func UserFromContextVerified(ctx context.Context, tv TokenVerifier) (*User, error) {
	ctx, span := trace.StartSpan(ctx, "identity.UserFromContextVerified")
	defer span.End()

	t, err := GetJWTFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err := tv.VerifyIDToken(ctx, t); err != nil {
		telemetry.SetError(span, trace.StatusCodeUnauthenticated, "UserFromContextVerified: Failed to verify JWT", err)
		return nil, errors.Join(ErrUnauthenticated, errJWTNotVerified, err)
	}
	u, err := UserFromJWT(t)
	if err != nil {
		return nil, err
	}
	if err := u.populateOrgAndProject(ctx); err != nil {
		return nil, err
	}
	return u, nil
}

// UserFromRequestVerified returns the verified user's identity from an incoming HTTP request.
func UserFromRequestVerified(r *http.Request, tv TokenVerifier) (*User, error) {
	_, span := trace.StartSpan(r.Context(), "identity.UserFromRequestVerified")
	defer span.End()

	jwt, err := GetJWTFromRequest(r)
	if err != nil {
		return nil, err
	}
	if err := tv.VerifyIDToken(r.Context(), jwt); err != nil {
		telemetry.SetError(span, trace.StatusCodeUnauthenticated, "UserFromRequestVerified: Failed to verify JWT", err)
		return nil, errors.Join(ErrUnauthenticated, errJWTNotVerified, err)
	}
	u, err := UserFromJWT(jwt)
	if err != nil {
		return nil, err
	}
	if err := u.populateOrgAndProjectFromRequest(r); err != nil {
		return nil, err
	}
	return u, nil
}

// UserFromRequest returns the user's identity from an HTTP request.
// No verification of the given JWT is performed.
func UserFromRequest(r *http.Request) (*User, error) {
	_, span := trace.StartSpan(r.Context(), "identity.UserFromRequest")
	defer span.End()

	jwt, err := GetJWTFromRequest(r)
	if err != nil {
		return nil, err
	}
	u, err := UserFromJWT(jwt)
	if err != nil {
		return nil, err
	}
	if err := u.populateOrgAndProjectFromRequest(r); err != nil {
		return nil, err
	}
	return u, nil
}

// UserFromContext returns the user's identity from a gRPC context.
// No verification of the given JWT is performed.
func UserFromContext(ctx context.Context) (*User, error) {
	ctx, span := trace.StartSpan(ctx, "identity.UserFromContext")
	defer span.End()

	t, err := GetJWTFromContext(ctx)
	if err != nil {
		return nil, err
	}
	u, err := UserFromJWT(t)
	if err != nil {
		return nil, err
	}
	if err := u.populateOrgAndProject(ctx); err != nil {
		return nil, err
	}
	return u, nil
}

// UserToRequest adds the user's identity to an HTTP request.
// Deprecated: use ToRequest with WithUser option.
func UserToRequest(r *http.Request, u *User) {
	_, span := trace.StartSpan(r.Context(), "identity.UserToRequest")
	defer span.End()

	ToRequest(r, WithUser(u))
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
// Deprecated: use AppendToOutgoingContext with WithUser option.
func UserToContext(ctx context.Context, u *User) (context.Context, error) {
	return AppendToOutgoingContext(ctx, WithUser(u))
}

func (i *User) lazyUnmarshalAndReturnData() (*jwt.Data, error) {
	if i.data != nil {
		return i.data, nil
	}
	d, err := jwt.UnmarshalUnsafe(i.jwt)
	if err != nil {
		return nil, errors.Join(ErrUnauthenticated, errJWTUnmarshal, err)
	}
	i.data = d // last write wins for parallel access
	return i.data, nil
}

// LoggableID returns an obfuscated name by which to identify a user in a log message.
func (i *User) LoggableID() string {
	submatches := emailRegex.FindStringSubmatch(i.EmailCanonicalized())
	var prefix, domain string
	if submatches == nil {
		prefix = ""
		domain = ""
	} else {
		groups := emailRegex.SubexpNames()
		prefix = submatches[slices.Index(groups, "prefix")]
		domain = submatches[slices.Index(groups, "domain")]
	}

	return fmt.Sprintf("%s@%s", obfuscateString(prefix), obfuscateString(domain))
}

// FirebaseUserID retrieves the user's user_id as defined by Firebase Auth.
// This is not named UserID or UID because we may want to introduce a different ID in the future
// which will be stored as claims in the JWT.
func (i *User) FirebaseUserID() string {
	dt, err := i.lazyUnmarshalAndReturnData()
	if err != nil {
		return ""
	}
	return dt.UserID
}

// JWTData retrieves the user's JWT data as defined by the JWT spec.
func (i *User) JWTData() (*jwt.Data, error) {
	return i.lazyUnmarshalAndReturnData()
}

// IsIntrinsicUser tests whether a user belongs to the Intrinsic organization.
// This is a temporary function for cases where we have no ACL checks yet.
func (i *User) IsIntrinsicUser() bool {
	return isIntrinsicUser(strings.ToLower(i.EmailRaw()))
}

func isIntrinsicUser(email string) bool {
	return strings.HasSuffix(email, googleEmailSuffix) ||
		strings.HasSuffix(email, IntrinsicServiceAccountEmailSuffix)
}

// UserToIncomingContext adds the user's identity to an incoming gRPC context.
// This is useful for testing and some special cases.
// Deprecated: use AppendToIncomingContext with WithUser option.
func UserToIncomingContext(ctx context.Context, u *User) (context.Context, error) {
	return AppendToIncomingContext(ctx, WithUser(u))
}

// RequestToContext returns a new context that has the user's identity, audience, and org-id stored in its
// metadata. The signature of the JWT is not verified.
// Use [RequestToContext] when chaining HTTP -> GRPC.
// Prefer [RequestToIncomingContext] over [RequestToContext] to avoid oversharing data with downstream services.
// See [ToContextFromIncoming] when chaining GRPC -> GRPC.
// Set the ctx parameter to r.Context() if there is no other more suitable context (e.g. a trace
// span derived from r.Context()).
// We require the user JWT to be on metadata["Cookie"]["auth-proxy"] because that is where the
// auth proxy service reads it from.
func RequestToContext(ctx context.Context, r *http.Request) (context.Context, error) {
	lctx, span := trace.StartSpan(ctx, "identity.RequestToContext")
	defer span.End()

	md, err := requestToMetadata(lctx, r)
	if err != nil {
		telemetry.SetError(span, trace.StatusCodeInvalidArgument, "RequestToContext: Failed to transform http requests identity to metadata", err)
		return ctx, err
	}
	return metadata.AppendToOutgoingContext(ctx, md...), nil
}

// RequestToIncomingContext returns a new context that has the user's identity, audience, and org-id stored in its
// metadata. The signature of the JWT is not verified.
func RequestToIncomingContext(ctx context.Context, r *http.Request) (context.Context, error) {
	lctx, span := trace.StartSpan(ctx, "identity.RequestToIncomingContext")
	defer span.End()

	md, err := requestToMetadata(lctx, r)
	if err != nil {
		return ctx, err
	}
	return metadata.NewIncomingContext(ctx, metadata.Pairs(md...)), nil
}

func requestToMetadata(ctx context.Context, r *http.Request) ([]string, error) {
	ctx, span := trace.StartSpan(ctx, "identity.requestToMetadata")
	defer span.End()

	userJWT, err := GetJWTFromRequest(r)
	if err != nil {
		return nil, err
	}
	// Ensure that authProxyCookieName is set to GetJWTFromRequest if it is not set yet.
	// This way downstream services do not have to check for the onpremTokenCookieName.
	if _, err := r.Cookie(authProxyCookieName); err != nil { // if not set yet, add it
		cookies.AddToRequest(r, &http.Cookie{Name: authProxyCookieName, Value: userJWT})
	}

	// Copy the relevant cookies to the context metadata
	cs := cookies.FromRequestNamed(r, []string{authProxyCookieName, org.OrgIDCookie})
	// backfill org-id cookie if it is not set, but passed as a header
	if _, err := r.Cookie(org.OrgIDCookie); err != nil {
		orgID := r.Header.Get(org.OrgIDCookie)
		if orgID != "" {
			log.WarningContextf(ctx, "Legacy org-id detected, request=%v, %s", r.URL, origin.Description(r))
			cs = append(cs, &http.Cookie{Name: org.OrgIDCookie, Value: orgID})
		} else {
			log.WarningContextf(ctx, "No org-id in request, request=%v, %s", r.URL, origin.Description(r))
		}
	}
	md := cookies.ToMDString(cs...)

	// Copy organization header to the context metadata
	o, err := headers.OrgFromRequest(r)
	if err != nil {
		return nil, err
	}
	if o != nil {
		md = append(md, org.OrgIDHeader, o.ID)
	}

	// Copy the audience to the context metadata.
	if aud, err := jwt.Aud(userJWT); err == nil {
		md = append(md, authProjectHeaderName, aud)
	} else {
		log.WarningContext(ctx, "no aud info in request")
	}
	return md, err
}

// ContextToRequest adds the auth related metadata from a GRPC context to an http request.
// Use ToRequest when chaining GRPC -> HTTP.
func ContextToRequest(ctx context.Context, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "identity.ContextToRequest")
	defer span.End()

	cookiesToCopy := []string{authProxyCookieName, portalCookieName, org.OrgIDCookie}
	// Copy the relevant cookies to the context request
	possibleCookies, err := cookies.FromContext(ctx)
	if err != nil {
		telemetry.SetError(span, trace.StatusCodeInvalidArgument, "ContextToRequest: Failed to get cookies from context", err)
		return errors.Join(ErrUnauthenticated, errCookiesParse, err)
	}
	var filteredCookies []*http.Cookie
	for _, c := range possibleCookies {
		if slices.Contains(cookiesToCopy, c.Name) {
			filteredCookies = append(filteredCookies, c)
		}
	}
	cookies.AddToRequest(r, filteredCookies...)

	metaToCopy := []string{apikeyTokenHeaderName, org.OrgIDHeader}
	md, _ := metadata.FromIncomingContext(ctx)
	for _, m := range metaToCopy {
		if val := md.Get(m); len(val) == 1 {
			r.Header.Set(m, val[0])
		}
	}

	return nil
}

// EnsureAuthProxyCookie ensures that the input request has an auth-proxy token cookie, copying it
// from the onprem-token cookie if needed.
func EnsureAuthProxyCookie(r *http.Request) error {
	_, span := trace.StartSpan(r.Context(), "identity.EnsureAuthProxyCookie")
	defer span.End()

	// Return early if the auth-proxy cookie already exists.
	if _, err := r.Cookie(authProxyCookieName); err == nil {
		return nil
	}

	// Try to copy the onprem-token or portal-token cookie to the auth-proxy cookie.
	for _, cn := range []string{onpremTokenCookieName, portalCookieName} {
		c, err := r.Cookie(cn)
		if err != nil {
			continue
		}
		cookies.AddToRequest(r, &http.Cookie{Name: authProxyCookieName, Value: c.Value})
		break
	}
	if _, err := r.Cookie(authProxyCookieName); err != nil {
		telemetry.SetError(span, trace.StatusCodeUnauthenticated, "EnsureAuthProxyCookie: No auth-proxy cookie found in request", err)
		return errors.Join(ErrUnauthenticated, errNoIntrinsicCookie, err)
	}

	return nil
}

var (
	stripNonEmailChars = regexp.MustCompile(`[^a-zA-Z0-9!#$%&'*+\-/=?^_{|}~` + "`" + `]`)
	// see gaia/data/email_util.cc;l=227
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

// Note: We are removing this as there should not be a need for clients using SDK
// to use any kind of metadata conversion methods.

// ToContextFromIncoming copies auth-related incoming GRPC metadata to outgoing
// metadata. The method does not error if auth-related information is not
// present. Use [ToContextFromIncomingChecked] to check if the context contains
// incoming authentication info. See [ToContextFromIncomingChecked] for more
// details.
func ToContextFromIncoming(ctx context.Context) (context.Context, error) {
	_, span := trace.StartSpan(ctx, "identity.ToContextFromIncoming")
	defer span.End()

	ctx, _, err := ToContextFromIncomingChecked(ctx)
	return ctx, err
}

// ToContextFromIncomingChecked copies auth-related incoming GRPC metadata to
// outgoing metadata. Returns false (and an unchanged context) if no relevant
// metadata was found. Use this when chaining GRPC requests
// (HTTP/GRPC->GRPC->GRPC).
//
// If any relevant values are already present on the outgoing context metadata,
// the values from the incoming context metadata will be appended to the
// existing values. This may be problematic for headers where web server may
// expect only a single value, such as the "Authorization" header. A warning
// will be logged if certain headers have more than one value after propagating
// incoming metadata.
func ToContextFromIncomingChecked(ctx context.Context) (context.Context, bool, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, false, nil
	}

	var changed bool

	cookieHeaders := md.Get(cookies.CookieHeaderName)
	if len(cookieHeaders) >= 1 { // only act if a cookie header is present in incoming
		newCtx, csChanged, err := setOutgoingValueCollisionAware(ctx, cookies.CookieHeaderName, cookieHeaders...)
		if err != nil {
			return ctx, false, err
		}
		changed = changed || csChanged
		ctx = newCtx
	}

	authHeaders := md.Get(authHeaderName)
	if len(authHeaders) > 1 {
		log.WarningContextf(ctx, "ToContextFromIncomingChecked: Multiple auth headers found in incoming context metadata: %v", authHeaders)
		return ctx, false, fmt.Errorf("%w: %w for %q in incoming context metadata", ErrInvalidRequest, errMetadataKeyConflict, authHeaderName)
	}
	if len(authHeaders) == 1 { // only act if a auth header is present in incoming
		newCtx, authChanged, err := setOutgoingValueCollisionAware(ctx, authHeaderName, authHeaders...)
		if err != nil {
			return ctx, false, err
		}
		changed = changed || authChanged
		ctx = newCtx
	}

	apikeyHeaders := md.Get(apikeyTokenHeaderName)
	if len(apikeyHeaders) > 1 {
		log.WarningContextf(ctx, "ToContextFromIncomingChecked: Multiple apikey headers found in incoming context metadata: %v", apikeyHeaders)
		return ctx, false, fmt.Errorf("%w: %w for %q in incoming context metadata", ErrInvalidRequest, errMetadataKeyConflict, apikeyTokenHeaderName)
	}
	if len(apikeyHeaders) == 1 { // only act if a apikey header is present in incoming
		newCtx, apikeyChanged, err := setOutgoingValueCollisionAware(ctx, apikeyTokenHeaderName, apikeyHeaders...)
		if err != nil {
			return ctx, false, err
		}
		changed = changed || apikeyChanged
		ctx = newCtx
	}

	orgHeaders := md.Get(org.OrgIDHeader)
	if len(orgHeaders) > 1 {
		orgHeaders = slices.Clone(orgHeaders)
		slices.Sort(orgHeaders)
		orgHeaders = slices.Compact(orgHeaders)
	}
	if len(orgHeaders) > 1 {
		log.WarningContextf(ctx, "ToContextFromIncomingChecked: Multiple org headers found in incoming context metadata: %v", orgHeaders)
		return ctx, false, fmt.Errorf("%w: %w for %q in incoming context metadata", ErrInvalidRequest, errMetadataKeyConflict, org.OrgIDHeader)
	}
	if len(orgHeaders) == 1 { // only act if a org header is present in incoming
		newCtx, orgChanged, err := setOutgoingValueCollisionAware(ctx, org.OrgIDHeader, orgHeaders...)
		if err != nil {
			return ctx, false, err
		}
		changed = changed || orgChanged
		ctx = newCtx
	}

	if changed {
		// Headers (except for "cookie") are not generally expected to have multiple
		// values. This can cause issues at target services. Printing a warning
		// might make odd looking errors easier to root cause.
		warnIfMultipleOutgoingValues(ctx, authHeaderName, apikeyTokenHeaderName, org.OrgIDHeader)
	}

	return ctx, changed, nil
}

func warnIfMultipleOutgoingValues(ctx context.Context, headers ...string) {
	mdOut, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return
	}
	for _, h := range headers {
		if vals := mdOut.Get(h); len(vals) > 1 {
			log.WarningContextf(ctx, "Header %q has %d values in outgoing metadata. Multiple values for this header may cause target services to reject requests with somewhat cryptic errors.", h, len(vals))
		}
	}
}

func obfuscateString(s string) string {
	if len(s) < 5 {
		return "***"
	}
	return obfuscateRegex.ReplaceAllString(s, `$2***$3`)
}

func setOutgoingValueCollisionAware(ctx context.Context, key string, vals ...string) (c context.Context, changed bool, err error) {
	lctx, span := trace.StartSpan(ctx, "identity.setOutgoingValueCollisionAware")
	defer span.End()
	span.AddAttributes(trace.StringAttribute("key", key))

	omd, ok := metadata.FromOutgoingContext(ctx)
	if !ok { // outgoing context is absent
		omd = metadata.MD{}
		omd.Set(key, vals...)
		return metadata.NewOutgoingContext(ctx, omd), true, nil
	}

	presentValues := omd.Get(key)

	// set the value if it's not present
	if len(presentValues) == 0 {
		omd.Set(key, vals...)
		return metadata.NewOutgoingContext(ctx, omd), true, nil
	}
	// return if the values are already present
	slices.Sort(presentValues)
	slices.Sort(vals)
	if slices.Equal(presentValues, vals) {
		return ctx, false, nil
	}

	log.WarningContextf(lctx, "Collision detected when setting values on outgoing context metadata for key %q. present outgoing values: %v, values that should get set: %v", key, presentValues, vals)
	telemetry.SetError(span, trace.StatusCodeInvalidArgument, "setOutgoingValueCollisionAware: Collision detected when setting values on outgoing context metadata", errMetadataKeyConflict)
	return ctx, false, fmt.Errorf("%w: %w for %q in outgoing context metadata", ErrInvalidRequest, errMetadataKeyConflict, key)
}

// OrgToRequest adds the organization identifier to the HTTP request.
// Deprecated: use ToRequest with WithOrg option.
func OrgToRequest(r *http.Request, orgID string) {
	ToRequest(r, WithOrg(orgID))
}

// OrgToContext returns a new context that has the org-id stored in its metadata.
// Deprecated: use AppendToOutgoingContext with WithOrg option.
func OrgToContext(ctx context.Context, orgID string) (context.Context, error) {
	if orgID == "" {
		log.WarningContextf(ctx, "OrgToContext: orgID is empty, returning unchanged context")
		return ctx, nil
	}
	return AppendToOutgoingContext(ctx, WithOrg(orgID))
}

// OrgFromRequest extracts the organization identifier from the HTTP header. If the header is not
// present, it will look for the org-id cookie.
func OrgFromRequest(r *http.Request) (*org.Organization, error) {
	_, span := trace.StartSpan(r.Context(), "identity.OrgFromRequest")
	defer span.End()

	if o, _ := headers.OrgFromRequest(r); o != nil {
		return o, nil
	}
	log.Warningf("OrgFromRequest could not retrieve the org from 'X-Intrinsic-Org'. Falling back to using cookies. URL=%v, %s", r.URL, origin.Description(r))

	organization, err := r.Cookie(org.OrgIDCookie)
	if err != nil {
		orgName := r.Header.Get(org.OrgIDCookie)
		if orgName != "" {
			log.Warningf("Legacy org-id detected in request=%v, %s", r.URL, origin.Description(r))
			return &org.Organization{ID: orgName}, nil
		}
		telemetry.SetError(span, trace.StatusCodeInvalidArgument, "OrgFromRequest", errNoOrgIDCookie)
		return nil, errors.Join(ErrMissingOrgID, errNoOrgIDCookie)
	}
	if organization.Value == "" {
		telemetry.SetError(span, trace.StatusCodeInvalidArgument, "OrgFromRequest", errOrgIDEmpty)
		return nil, errors.Join(ErrMissingOrgID, errOrgIDEmpty)
	}

	return &org.Organization{ID: organization.Value}, nil
}

// OrgToIncomingContext returns a new context that has the org-id stored in its metadata.
// This is useful for testing and some special cases.
// Deprecated: use AppendToIncomingContext with WithOrg option.
func OrgToIncomingContext(ctx context.Context, orgID string) (context.Context, error) {
	if orgID == "" {
		log.WarningContextf(ctx, "OrgToIncomingContext: orgID is empty, returning unchanged context")
		return ctx, nil
	}
	return AppendToIncomingContext(ctx, WithOrg(orgID))
}

// OrgFromContext extracts the organization identifier from the  gRPC context.
func OrgFromContext(ctx context.Context) (*org.Organization, error) {
	ctx, span := trace.StartSpan(ctx, "identity.OrgFromContext")
	defer span.End()

	if o, err := headers.OrgFromContext(ctx); err != nil {
		return nil, err
	} else if o != nil {
		return o, nil
	}
	log.WarningContext(ctx, "OrgFromContext: No org header found. Falling back to using cookies.")

	contextCookies, err := cookies.FromContext(ctx)
	if err != nil {
		telemetry.SetError(span, trace.StatusCodeInvalidArgument, "OrgFromContext: Failed to get cookies from context", err)
		return nil, errors.Join(ErrMissingOrgID, errCookiesParse, err)
	}
	for _, cookie := range contextCookies {
		if cookie.Name == org.OrgIDCookie {
			log.V(2).InfoContextf(ctx, "Using org from cookie %q", cookie)
			return &org.Organization{ID: cookie.Value}, nil
		}
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		telemetry.SetError(span, trace.StatusCodeInvalidArgument, "OrgFromContext", errNoMetadata)
		return nil, errors.Join(ErrMissingOrgID, errNoMetadata)
	}
	orgMD, ok := md[org.OrgIDCookie]
	if ok && len(orgMD) > 0 && orgMD[0] != "" {
		log.ErrorContextf(ctx, "Tried using org from context metadata instead of a cookie %q. Update your code to use cookies instead.", org.OrgIDCookie)
		return nil, errors.Join(ErrMissingOrgID, errOrgInMetadata, errNoOrgIDCookie)
	}

	telemetry.SetError(span, trace.StatusCodeInvalidArgument, "OrgFromContext", errNoOrgIDCookie)
	return nil, errors.Join(ErrMissingOrgID, errNoOrgIDCookie)
}

// UserOrgFromContext extracts the user and organization identifier from the gRPC context.
// Both the user and organization must be present.
// On absence [ErrUnauthenticated] or [ErrMissingOrgID] will be returned .
// Deprecated: use UserFromContext and the returned User's Org() getter.
func UserOrgFromContext(ctx context.Context) (*UserOrg, error) {
	ctx, span := trace.StartSpan(ctx, "identity.UserOrgFromContext")
	defer span.End()

	u, err := UserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	o, err := OrgFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return &UserOrg{User: u, Org: o}, nil
}

// UserOrgFromRequest extracts the user and organization identifier from the HTTP request.
// Both the user and organization must be present.
// On absence [ErrUnauthenticated] or [ErrMissingOrgID] will be returned .
func UserOrgFromRequest(r *http.Request) (*UserOrg, error) {
	_, span := trace.StartSpan(r.Context(), "identity.UserOrgFromRequest")
	defer span.End()

	u, err := UserFromRequest(r)
	if err != nil {
		return nil, err
	}
	o, err := OrgFromRequest(r)
	if err != nil {
		return nil, err
	}
	return &UserOrg{User: u, Org: o}, nil
}

// ToRequestDeprecated adds the user and org metadata to the HTTP request.
// Deprecated: use ToRequest with WithUser and WithOrg options.
func ToRequestDeprecated(r *http.Request, u *User, orgID string) {
	_, span := trace.StartSpan(r.Context(), "identity.ToRequest")
	defer span.End()

	UserToRequest(r, u)
	OrgToRequest(r, orgID)
}

// ToContext adds the user and org metadata to the context.
// Deprecated: use AppendToOutgoingContext with WithUser and WithOrg options.
func ToContext(ctx context.Context, u *User, orgID string) (context.Context, error) {
	var opts []Option
	if u != nil {
		opts = append(opts, WithUser(u))
	}
	if orgID != "" {
		opts = append(opts, WithOrg(orgID))
	}
	return AppendToOutgoingContext(ctx, opts...)
}

// ToIncomingContext adds the user and org metadata to the incoming context.
// This is useful for testing and some special cases.
// Deprecated: use AppendToIncomingContext with WithUser and WithOrg options.
func ToIncomingContext(ctx context.Context, u *User, orgID string) (context.Context, error) {
	var opts []Option
	if u != nil {
		opts = append(opts, WithUser(u))
	}
	if orgID != "" {
		opts = append(opts, WithOrg(orgID))
	}
	return AppendToIncomingContext(ctx, opts...)
}

// ClearRequest removes the user and org metadata from the HTTP request.
func ClearRequest(r *http.Request) {
	_, span := trace.StartSpan(r.Context(), "identity.ClearRequest")
	defer span.End()

	ToRequest(r, WithClearUser())
}

func filterCookies(cookies []*http.Cookie, filter []string) []*http.Cookie {
	var validCookies []*http.Cookie
	for _, c := range cookies {
		if !slices.Contains(filter, c.Name) {
			validCookies = append(validCookies, c)
		}
	}
	return validCookies
}

// ClearRequestOrg removes the org metadata from the HTTP request.
func ClearRequestOrg(r *http.Request) {
	_, span := trace.StartSpan(r.Context(), "identity.ClearRequestOrg")
	defer span.End()

	ToRequest(r, withClearUserOrg())
}

// ClearRequestUser removes the user metadata from the HTTP request.
func ClearRequestUser(r *http.Request) {
	_, span := trace.StartSpan(r.Context(), "identity.ClearRequestUser")
	defer span.End()

	ToRequest(r, withClearUserAuth())
}

// ClearContext removes the user and org metadata from the outgoing context.
func ClearContext(ctx context.Context) (context.Context, error) {
	_, span := trace.StartSpan(ctx, "identity.ClearContext")
	defer span.End()

	return AppendToOutgoingContext(ctx, WithClearUser())
}

// ClearContextOrg removes the org metadata from the outgoing context.
func ClearContextOrg(ctx context.Context) (context.Context, error) {
	_, span := trace.StartSpan(ctx, "identity.ClearContextOrg")
	defer span.End()

	return AppendToOutgoingContext(ctx, withClearUserOrg())
}

// ClearContextUser removes the user metadata from the outgoing context.
func ClearContextUser(ctx context.Context) (context.Context, error) {
	_, span := trace.StartSpan(ctx, "identity.ClearContextUser")
	defer span.End()

	return AppendToOutgoingContext(ctx, withClearUserAuth())
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

// Option configures the context/request propagation.
type Option func(*update) error

type update struct {
	headers      map[string]string
	cookies      []*http.Cookie
	clearHeaders []string
	clearCookies []string
}

func (u *update) addHeader(k, v string) {
	if u.headers == nil {
		u.headers = make(map[string]string)
	}
	u.headers[k] = v
}

func (u *update) addCookie(c *http.Cookie) {
	if c != nil {
		u.cookies = append(u.cookies, c)
	}
}

func (u *update) clearHeader(k string) {
	u.clearHeaders = append(u.clearHeaders, k)
}

func (u *update) clearCookie(k string) {
	u.clearCookies = append(u.clearCookies, k)
}

func (u *update) applyToMD(md metadata.MD) (metadata.MD, error) {
	newMD := md.Copy()
	if newMD == nil {
		newMD = metadata.MD{}
	}

	// 1. Apply all header changes first
	for _, kh := range u.clearHeaders {
		newMD.Delete(kh)
	}
	for k, v := range u.headers {
		newMD.Set(k, v)
	}

	// 2. Guard: If no cookie changes are needed, return early
	if len(u.clearCookies) == 0 && len(u.cookies) == 0 {
		return newMD, nil
	}

	// 3. Guard: If we only have additions (no clearing), handle and return early
	if len(u.clearCookies) == 0 {
		return cookies.AddToMD(newMD, u.cookies...)
	}

	// 4. Handle clearing (and potential merging with additions)
	existing, err := cookies.FromMD(newMD)
	if err != nil {
		return nil, err
	}

	finalCookies := filterCookies(existing, u.clearCookies)
	finalCookies = append(finalCookies, u.cookies...)
	newMD.Delete(cookies.CookieHeaderName)
	return cookies.AddToMD(newMD, finalCookies...)
}

func withCookie(c *http.Cookie) Option {
	return func(u *update) error {
		u.addCookie(c)
		return nil
	}
}

func withHeader(k, v string) Option {
	return func(u *update) error {
		u.addHeader(k, v)
		return nil
	}
}

// WithOrg scopes the organization identifier.
func WithOrg(orgID string) Option {
	return func(u *update) error {
		if orgID == "" {
			return ErrMissingOrgID
		}
		if err := withCookie(org.IDCookie(orgID))(u); err != nil {
			return err
		}
		return withHeader(org.OrgIDHeader, orgID)(u)
	}
}

// WithUserJWT scopes the user's JWT token.
func WithUserJWT(jwt string) Option {
	return func(u *update) error {
		if jwt == "" {
			return ErrUnauthenticated
		}
		return withCookie(&http.Cookie{Name: authProxyCookieName, Value: jwt})(u)
	}
}

// WithComputeProject scopes the compute project identifier.
func WithComputeProject(projectID string) Option {
	return func(u *update) error {
		if projectID == "" {
			return ErrMissingProject
		}
		return withHeader(authProjectHeaderName, projectID)(u)
	}
}

// WithUser scopes the user's identity, organization, and compute project.
func WithUser(u *User) Option {
	return func(up *update) error {
		if u == nil {
			return fmt.Errorf("%w: user cannot be nil", ErrInvalidRequest)
		}
		if err := WithUserJWT(u.jwt)(up); err != nil {
			return err
		}
		if u.org != nil {
			if err := WithOrg(u.org.ID)(up); err != nil {
				return err
			}
		}
		if u.project != "" {
			if err := WithComputeProject(u.project)(up); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithClearUser removes all user identity credentials, API keys, and
// authentication cookies from the context or request.
//
// Note: if options to add identity information are also present, this will
// clear identity data first, then add the new data. This allows clients to
// clear -> add in a single call, instead of two separate calls.
func WithClearUser() Option {
	return func(u *update) error {
		if err := withClearUserOrg()(u); err != nil {
			return err
		}
		return withClearUserAuth()(u)
	}
}

func withClearUserOrg() Option {
	return func(u *update) error {
		u.clearHeader(org.OrgIDHeader)
		u.clearHeader(org.OrgIDCookie) // clears legacy org-id header
		u.clearCookie(org.OrgIDCookie) // clears org-id cookie
		return nil
	}
}

func withClearUserAuth() Option {
	return func(u *update) error {
		u.clearHeader(authHeaderName)
		u.clearHeader(apikeyTokenHeaderName)
		u.clearHeader(authProjectHeaderName)
		for _, ch := range cookieHeaders {
			u.clearCookie(ch)
		}
		return nil
	}
}

// NewOutgoingContext creates a new outgoing context with the provided options.
func NewOutgoingContext(ctx context.Context, opts ...Option) (context.Context, error) {
	if len(opts) == 0 {
		return ctx, nil
	}

	_, span := trace.StartSpan(ctx, "identity.NewOutgoingContext")
	defer span.End()

	u := &update{}
	for _, opt := range opts {
		if err := opt(u); err != nil {
			return ctx, err
		}
	}

	newMD, err := u.applyToMD(metadata.MD{})
	if err != nil {
		return ctx, err
	}

	return metadata.NewOutgoingContext(ctx, newMD), nil
}

// AppendToOutgoingContext appends/merges the provided options into the outgoing context metadata.
func AppendToOutgoingContext(ctx context.Context, opts ...Option) (context.Context, error) {
	if len(opts) == 0 {
		return ctx, nil
	}

	_, span := trace.StartSpan(ctx, "identity.AppendToOutgoingContext")
	defer span.End()

	md, _ := metadata.FromOutgoingContext(ctx)
	newMD, err := AppendToMetadata(md, opts...)
	if err != nil {
		return ctx, err
	}

	return metadata.NewOutgoingContext(ctx, newMD), nil
}

// ToRequest modifies the HTTP request headers and cookies based on options.
func ToRequest(r *http.Request, opts ...Option) error {
	if r == nil {
		return fmt.Errorf("%w: request cannot be nil", ErrInvalidRequest)
	}
	if len(opts) == 0 {
		return nil
	}

	_, span := trace.StartSpan(r.Context(), "identity.ToRequest")
	defer span.End()

	u := &update{}
	for _, opt := range opts {
		if err := opt(u); err != nil {
			return err
		}
	}

	for _, kh := range u.clearHeaders {
		r.Header.Del(kh)
	}

	if len(u.clearCookies) > 0 {
		cs := filterCookies(r.Cookies(), u.clearCookies)
		r.Header.Del(cookies.CookieHeaderName)
		cookies.AddToRequest(r, cs...)
	}

	for k, v := range u.headers {
		r.Header.Set(k, v)
	}

	if len(u.cookies) > 0 {
		cookies.AddToRequest(r, u.cookies...)
	}
	return nil
}

// AppendToIncomingContext appends/merges the provided options to the incoming gRPC context.
func AppendToIncomingContext(ctx context.Context, opts ...Option) (context.Context, error) {
	if len(opts) == 0 {
		return ctx, nil
	}

	_, span := trace.StartSpan(ctx, "identity.AppendToIncomingContext")
	defer span.End()

	md, _ := metadata.FromIncomingContext(ctx)
	newMD, err := AppendToMetadata(md, opts...)
	if err != nil {
		return ctx, err
	}

	return metadata.NewIncomingContext(ctx, newMD), nil
}

// AppendToMetadata merges the provided identity options into an existing gRPC metadata map.
func AppendToMetadata(md metadata.MD, opts ...Option) (metadata.MD, error) {
	if len(opts) == 0 {
		return md, nil
	}

	u := &update{}
	for _, opt := range opts {
		if err := opt(u); err != nil {
			return md, err
		}
	}

	return u.applyToMD(md)
}

// ToMetadata serializes the provided identity options directly into a gRPC metadata map.
func ToMetadata(opts ...Option) (metadata.MD, error) {
	return AppendToMetadata(nil, opts...)
}

// AppendToMetadataMap merges identity options into an existing flat metadata map.
func AppendToMetadataMap(m map[string]string, opts ...Option) (map[string]string, error) {
	if len(opts) == 0 {
		return m, nil
	}

	md := make(metadata.MD, len(m))
	for k, v := range m {
		md.Set(k, v)
	}

	newMD, err := AppendToMetadata(md, opts...)
	if err != nil {
		return m, err
	}

	if m == nil {
		m = make(map[string]string, len(newMD))
	} else {
		for k := range m {
			delete(m, k)
		}
	}
	for k, vs := range newMD {
		m[k] = strings.Join(vs, ", ")
	}
	return m, nil
}

// ToMetadataMap flattens identity options into a metadata key-value map from scratch.
func ToMetadataMap(opts ...Option) (map[string]string, error) {
	return AppendToMetadataMap(nil, opts...)
}
