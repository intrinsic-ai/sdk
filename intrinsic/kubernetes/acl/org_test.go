// Copyright 2023 Intrinsic Innovation LLC

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
