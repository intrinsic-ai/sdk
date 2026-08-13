// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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