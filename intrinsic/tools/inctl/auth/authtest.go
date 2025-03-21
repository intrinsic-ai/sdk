// Copyright 2023 Intrinsic Innovation LLC

// Package authtest provides test helpers.
package authtest

import (
	"testing"

	"intrinsic/tools/inctl/auth/auth"
)

// NewStoreForTest creates a new auth.Store for use in tests.
func NewStoreForTest(t *testing.T) *auth.Store {
	configDir := t.TempDir()
	return &auth.Store{GetConfigDirFx: func() (string, error) { return configDir, nil }}
}
