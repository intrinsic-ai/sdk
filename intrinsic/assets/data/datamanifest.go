// Copyright 2023 Intrinsic Innovation LLC

// Package datamanifest contains tools for working with DataManifest.
package datamanifest

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoregistry"
	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"

	dmpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
)

// ValidateDataManifestOptions contains options for validating a DataManifest.
type ValidateDataManifestOptions struct {
	types *protoregistry.Types
}

// ValidateDataManifestOption is an option for validating a DataManifest.
type ValidateDataManifestOption func(*ValidateDataManifestOptions)

// WithTypes sets the Types option.
func WithTypes(types *protoregistry.Types) ValidateDataManifestOption {
	return func(opts *ValidateDataManifestOptions) {
		opts.types = types
	}
}

// ValidateDataManifest validates a data manifest.
func ValidateDataManifest(m *dmpb.DataManifest, options ...ValidateDataManifestOption) error {
	opts := &ValidateDataManifestOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if err := metadatautils.ValidateManifestMetadata(m.GetMetadata()); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetMetadata().GetId())

	if m.GetData() == nil {
		return fmt.Errorf("data must be specified for Data %q", id)
	}

	if opts.types != nil {
		if name := m.GetData().MessageName(); name != "" {
			if _, err := opts.types.FindMessageByName(name); err != nil {
				return fmt.Errorf("cannot find data message %q for Data %q: %w", name, id, err)
			}
		}
	}

	return nil
}
