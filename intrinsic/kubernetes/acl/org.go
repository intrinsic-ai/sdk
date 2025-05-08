// Copyright 2023 Intrinsic Innovation LLC

// Package org provides helpers to deal with organizations in requests and code.
package org

import (
	"net/http"
	"regexp"
	"strings"

	log "github.com/golang/glog"
)

// OrgIDCookie is the cookie key for the organization identifier.
const OrgIDCookie = "org-id"

// IntrinsicOrgID is the organization identifier used for the Intrinsic in multi-tenant projects.
const IntrinsicOrgID = "intrinsic"

// PublicOrgID is an ACL-only organization used for public resources.
const PublicOrgID = "publicorg"

// Organization represents an organization inside the Intrinsic stack.
type Organization struct {
	ID string
}

var (
	// use email as a base for org id, as per section 2.3 of RFC 3986 replace other chars by '_'.
	// The result must be a valid label value for compatibility with quota tracking:
	// https://cloud.google.com/compute/docs/labeling-resources#requirements
	stripEmail = regexp.MustCompile(`[^a-zA-Z0-9_-]`)
)

// IDCookie returns a cookie with the given orgID.
func IDCookie(orgID string) *http.Cookie {
	return &http.Cookie{Name: OrgIDCookie, Value: orgID}
}

// GetID returns the identifier of the organization.
func (o *Organization) GetID() string {
	return o.ID
}

// NewFromUID makes a new org with an ID based on an email.
func NewFromUID(uid string) *Organization {
	orgid := strings.ToLower(stripEmail.ReplaceAllString(uid, "_"))
	if len(orgid) > 63 {
		log.Warningf("OrgID %q is too long, truncating to %q", orgid, orgid[:63])
		orgid = orgid[:63]
	}
	return &Organization{ID: orgid}
}
