// Copyright 2023 Intrinsic Innovation LLC

// Package datamanifest provides utils for working with Data Asset manifests.
package datamanifest

import (
	"fmt"

	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"

	"google.golang.org/protobuf/reflect/protoregistry"

	dmpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
)

type validateDataManifestOptions struct {
	types *protoregistry.Types
}

// ValidateDataManifestOption is an option for validating a DataManifest.
type ValidateDataManifestOption func(*validateDataManifestOptions)

// WithTypes provides a Types for validating proto messages.
func WithTypes(types *protoregistry.Types) ValidateDataManifestOption {
	return func(opts *validateDataManifestOptions) {
		opts.types = types
	}
}

// ValidateDataManifest validates a DataManifest.
func ValidateDataManifest(m *dmpb.DataManifest, options ...ValidateDataManifestOption) error {
	opts := &validateDataManifestOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if m == nil {
		return fmt.Errorf("DataManifest must not be nil")
	}

	if err := metadatautils.ValidateManifestMetadata(m.GetMetadata()); err != nil {
		return fmt.Errorf("invalid DataManifest metadata: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetMetadata().GetId())

	if m.GetData() == nil {
		return fmt.Errorf("data payload must be specified for %q", id)
	}

	if opts.types != nil {
		if name := m.GetData().MessageName(); name != "" {
			if _, err := opts.types.FindMessageByName(name); err != nil {
				return fmt.Errorf("cannot find data message %q for %q: %w", name, id, err)
			}
		}
	}

	return nil
}
