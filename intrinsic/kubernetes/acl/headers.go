// Copyright 2023 Intrinsic Innovation LLC

package headers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"intrinsic/kubernetes/acl/org"
	"intrinsic/stats/go/telemetry"

	log "github.com/golang/glog"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/metadata"
)

var (
	// ErrMultipleOrgsInHeader indicates that there are multiple organizations specified as part of the X-Intrinsic-Org header.
	ErrMultipleOrgsInHeader = errors.New("multiple organizations specified in X-Intrinsic-Org header")
)

// OrgFromRequest extracts the organization identifier from the HTTP request.
func OrgFromRequest(r *http.Request) (*org.Organization, error) {
	if r == nil {
		return nil, nil
	}
	ctx, span := trace.StartSpan(r.Context(), "headers.OrgFromRequest")
	defer span.End()

	orgs := r.Header.Values(org.OrgIDHeader)
	return validateAndExtractOrgID(ctx, orgs)
}

// AddOrgToRequest adds the org header to the request.
// It overwrites the existing header should it be set.
func AddOrgToRequest(r *http.Request, orgID string) {
	if r == nil {
		return
	}

	_, span := trace.StartSpan(r.Context(), "headers.AddOrgToRequest")
	span.AddAttributes(trace.StringAttribute("org_id", orgID))
	defer span.End()

	if r.Header == nil {
		r.Header = make(http.Header)
	}

	r.Header.Set(org.OrgIDHeader, orgID)
}

// AddOrgToContext adds the org header to the context.
// It overwrites the existing header should it be set.
func AddOrgToContext(ctx context.Context, orgID string) context.Context {
	_, span := trace.StartSpan(ctx, "headers.AddOrgToContext")
	span.AddAttributes(trace.StringAttribute("org_id", orgID))
	defer span.End()

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.NewOutgoingContext(ctx, metadata.Pairs(org.OrgIDHeader, orgID))
	}

	newMD := md.Copy()
	if newMD == nil {
		newMD = metadata.MD{}
	}
	newMD.Set(org.OrgIDHeader, orgID)
	return metadata.NewOutgoingContext(ctx, newMD)
}

// OrgFromContext extracts the organization identifier from the gRPC context.
func OrgFromContext(ctx context.Context) (*org.Organization, error) {
	ctx, span := trace.StartSpan(ctx, "headers.OrgFromContext")
	defer span.End()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, nil
	}
	orgs := md.Get(org.OrgIDHeader)
	return validateAndExtractOrgID(ctx, orgs)
}

func validateAndExtractOrgID(ctx context.Context, orgs []string) (*org.Organization, error) {
	ctx, span := trace.StartSpan(ctx, "headers.validateAndExtractOrgID")
	defer span.End()

	if len(orgs) > 1 {
		telemetry.SetError(span, trace.StatusCodeInvalidArgument, fmt.Sprintf("Multiple organizations specified in the %q header.", org.OrgIDHeader), ErrMultipleOrgsInHeader)
		log.ErrorContextf(ctx, "Multiple organizations specified in the %q header, only a single organization value is supported.", org.OrgIDHeader)
		return nil, ErrMultipleOrgsInHeader
	}
	if len(orgs) == 1 {
		if orgs[0] == "" {
			log.WarningContextf(ctx, "Header %q specifies an empty organization. This is likely an implementation error. Falling back to using the organization from cookies.", org.OrgIDHeader)
		} else {
			log.V(2).InfoContextf(ctx, "Using org from header %q", org.OrgIDHeader)
			return &org.Organization{ID: orgs[0]}, nil
		}
	}

	return nil, nil
}
