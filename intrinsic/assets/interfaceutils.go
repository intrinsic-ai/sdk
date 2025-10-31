// Copyright 2023 Intrinsic Innovation LLC

// Package interfaceutils provides utilities for working with Asset URI interfaces.
package interfaceutils

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

const (
	// GRPCURIPrefix is the prefix used for gRPC service dependencies.
	GRPCURIPrefix = "grpc://"
	// DataURIPrefix is the prefix used for proto-based data dependencies.
	DataURIPrefix = "data://"
)

var (
	uriRegex = regexp.MustCompile(`^(grpc://|data://)([A-Za-z_][A-Za-z0-9_]*\.)+[A-Za-z_][A-Za-z0-9_]*$`)

	// ErrInvalidInterfaceName is returned when an interface name is invalid.
	ErrInvalidInterfaceName = errors.New("invalid interface name")
)

// ValidateInterfaceName validates an interface name with a protocol prefix.
func ValidateInterfaceName(uri string) error {
	if !uriRegex.MatchString(uri) {
		return fmt.Errorf("%w: expected URI to be formatted as '<protocol>://<package>.<message>', got %q", ErrInvalidInterfaceName, uri)
	}
	return nil
}
