// Copyright 2023 Intrinsic Innovation LLC

package org

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewFromUID(t *testing.T) {
	tests := []struct {
		uid  string
		want string
	}{
		{
			uid:  "johndoe@gmail.com",
			want: "johndoe_gmail_com",
		},
	}

	for _, tc := range tests {
		org := NewFromUID(tc.uid)
		got := org.GetID()
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("NewFromUID(%v) returned an unexpected diff (-want +got): %v", tc.uid, diff)
		}
	}
}

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
