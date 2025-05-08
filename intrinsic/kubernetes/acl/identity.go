// Copyright 2023 Intrinsic Innovation LLC

// Package identity provides helpers to work with user identities inside the Intrinsic stack.
package identity

import (
	"context"

	log "github.com/golang/glog"
	"google.golang.org/grpc/metadata"
	"intrinsic/kubernetes/acl/cookies"
	"intrinsic/kubernetes/acl/org"
)

// OrgToContext returns a new context that has the org-id stored in its metadata.
func OrgToContext(ctx context.Context, orgID string) context.Context {
	if orgID == "" {
		log.WarningContextf(ctx, "No org-id in context, returning unchanged context")
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, cookies.ToMDString(org.IDCookie(orgID))...)
}
