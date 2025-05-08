// Copyright 2023 Intrinsic Innovation LLC

package identity

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
	"intrinsic/kubernetes/acl/cookies"
	"intrinsic/kubernetes/acl/org"
)

func TestOrgToContext(t *testing.T) {
	ctx := context.Background()
	ctx = OrgToContext(ctx, "testorg")
	md, _ := metadata.FromOutgoingContext(ctx)
	if md.Get(cookies.CookieHeaderName)[0] != org.OrgIDCookie+"=testorg" {
		t.Errorf("UserToContext(..) did not add the user's identity to the context")
	}
}
