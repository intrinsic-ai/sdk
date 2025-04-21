// Copyright 2023 Intrinsic Innovation LLC

// Package servicemanifest contains tools for working with ServiceManifest.
package servicemanifest

import (
	"fmt"
	"slices"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"
	smpb "intrinsic/assets/services/proto/service_manifest_go_proto"
	"intrinsic/util/proto/names"
)

var (
)

// ValidateServiceManifestOptions contains options for validating a ServiceManifest.
type ValidateServiceManifestOptions struct {
	files *protoregistry.Files
}

// ValidateServiceManifestOption is an option for validating a ServiceManifest.
type ValidateServiceManifestOption func(*ValidateServiceManifestOptions)

// WithFiles adds the proto files to the validation options.
func WithFiles(files *protoregistry.Files) ValidateServiceManifestOption {
	return func(opts *ValidateServiceManifestOptions) {
		opts.files = files
	}
}

// ValidateServiceManifest checks that a ServiceManifest is consistent and valid.
func ValidateServiceManifest(m *smpb.ServiceManifest, options ...ValidateServiceManifestOption) error {
	opts := &ValidateServiceManifestOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if err := metadatautils.ValidateManifestMetadata(m.GetMetadata()); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetMetadata().GetId())

	if m.GetServiceDef() != nil && m.GetServiceDef().GetSimSpec() == nil {
		return fmt.Errorf("a sim_spec must be specified if a service_def is provided for Service %q", id)
	}

	for _, p := range m.GetServiceDef().GetServiceProtoPrefixes() {
		if err := names.ValidateProtoPrefix(p); err != nil {
			return fmt.Errorf("service proto prefix %q is not valid for Service %q: %w", p, id, err)
		}
	}

	expectedImagePaths := map[string]struct{}{}
	if name := m.GetServiceDef().GetRealSpec().GetImage().GetArchiveFilename(); name != "" {
		expectedImagePaths[name] = struct{}{}
	}
	if name := m.GetServiceDef().GetSimSpec().GetImage().GetArchiveFilename(); name != "" {
		expectedImagePaths[name] = struct{}{}
	}
	for p := range expectedImagePaths {
		if !slices.Contains(m.GetAssets().GetImageFilenames(), p) {
			return fmt.Errorf("image %q in the manifest for Service %q is not listed in the service assets", p, id)
		}
	}
	for _, p := range m.GetAssets().GetImageFilenames() {
		if _, ok := expectedImagePaths[p]; !ok {
			return fmt.Errorf("image %q in the service assets for Service %q is not used in the manifest", p, id)
		}
	}

	if opts.files != nil {
		for _, prefix := range m.GetServiceDef().GetServiceProtoPrefixes() {
			strippedPrefix := strings.TrimSuffix(strings.TrimPrefix(prefix, "/"), "/")
			name := protoreflect.FullName(strippedPrefix)
			if _, err := opts.files.FindDescriptorByName(name); err != nil {
				return fmt.Errorf("could not find service proto prefix %q in provided descriptors for Service %q: %w", prefix, id, err)
			}
		}
	}

	return nil
}

func setDifference(slice1, slice2 []string) []string {
	var difference []string
	for _, val := range slice1 {
		if !slices.Contains(slice2, val) {
			difference = append(difference, val)
		}
	}
	return difference
}
