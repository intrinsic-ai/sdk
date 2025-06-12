// Copyright 2023 Intrinsic Innovation LLC

// Package org provides helpers to deal with organizations in requests and code.
package org

import (
	"net/http"

)

// OrgIDCookie is the cookie key for the organization identifier.
const OrgIDCookie = "org-id"

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

//