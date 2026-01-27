// Copyright 2023 Intrinsic Innovation LLC

// Package org provides helpers to deal with organizations in requests and code.
package org

import (
	"net/http"

	"github.com/rs/xid"
)

// OrgIDCookie is the cookie key for the organization identifier.
const OrgIDCookie = "org-id"

const OrgIDHeader = "x-intrinsic-org"

// Organization represents an organization inside the Intrinsic stack.
type Organization struct {
	ID string
}

// IDCookie returns a cookie with the given orgID.
func IDCookie(orgID string) *http.Cookie {
	return &http.Cookie{Name: OrgIDCookie, Value: orgID}
}

// GetID returns the identifier of the organization.
func (o *Organization) GetID() string {
	return o.ID
}

// ID returns a random organization ID.
func ID() string {
	return xid.New().String()
}

//