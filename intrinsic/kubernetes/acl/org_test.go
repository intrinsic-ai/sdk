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

package org

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIDCookie(t *testing.T) {
	t.Run("org-id-cookie", func(t *testing.T) {
		orgID := "testorg"
		expectedCookie := &http.Cookie{Name: OrgIDCookie, Value: orgID}

		c := IDCookie(orgID)
		if diff := cmp.Diff(expectedCookie, c); diff != "" {
			t.Errorf("IDCookie(%q) returned an unexpected diff (-want +got): %v", orgID, diff)
		}
	})
}

func TestID(t *testing.T) {
	id1 := ID()
	id2 := ID()
	if id1 == id2 {
		t.Errorf("ID() returned the same ID: %v", id1)
	}
	if len(id1) > 63 {
		t.Errorf("ID() returns an id that is too long: %v, %d", id1, len(id1))
	}
}
