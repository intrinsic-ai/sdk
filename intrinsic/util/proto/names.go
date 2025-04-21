// Copyright 2023 Intrinsic Innovation LLC

// Package names provides utilities for proto names.
package names

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// See: https://protobuf.com/docs/language-spec.
	nameRegex = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_]*\.)+[A-Za-z_][A-Za-z0-9_]*$`)

	errInvalidProtoName   = errors.New("invalid proto name")
	errInvalidProtoPrefix = errors.New("invalid proto prefix")
)

// ValidateProtoName validates a proto name.
func ValidateProtoName(protoName string) error {
	if !nameRegex.MatchString(protoName) {
		return fmt.Errorf("%w: expected name formatted as '<package>.<message>', got %q", errInvalidProtoName, protoName)
	}
	return nil
}

// ValidateProtoPrefix validates a proto prefix.
func ValidateProtoPrefix(protoPrefix string) error {
	if len(protoPrefix) < 2 || !strings.HasPrefix(protoPrefix, "/") || !strings.HasSuffix(protoPrefix, "/") {
		return fmt.Errorf("%w: expected prefix formatted as '/<package>.<message>/', got %q", errInvalidProtoPrefix, protoPrefix)
	}
	protoName := protoPrefix[1 : len(protoPrefix)-1]
	if !nameRegex.MatchString(protoName) {
		return fmt.Errorf("%w: expected prefix formatted as '/<package>.<message>/', got %q", errInvalidProtoPrefix, protoPrefix)
	}
	return nil
}
