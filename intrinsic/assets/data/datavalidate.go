// Copyright 2023 Intrinsic Innovation LLC

// Package datavalidate provides utils for validating Data Assets.
package datavalidate

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"intrinsic/assets/data/utils"
	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"

	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	dmpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
	atpb "intrinsic/assets/proto/asset_type_go_proto"
)

type dataManifestOptions struct {
	files *protoregistry.Files
}

// DataManifestOption is an option for validating a DataManifest.
type DataManifestOption func(*dataManifestOptions)

// WithFiles provides a Files for validating proto messages.
func WithFiles(files *protoregistry.Files) DataManifestOption {
	return func(opts *dataManifestOptions) {
		opts.files = files
	}
}

// DataManifest validates a DataManifest.
func DataManifest(m *dmpb.DataManifest, options ...DataManifestOption) error {
	opts := &dataManifestOptions{}
	for _, opt := range options {
		opt(opts)
	}
	if opts.files == nil {
		return fmt.Errorf("files option must be specified")
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
	if name := m.GetData().MessageName(); name == "" {
		return fmt.Errorf("data payload must not be an empty Any for %q", id)
	} else if _, err := opts.files.FindDescriptorByName(name); err != nil {
		return fmt.Errorf("cannot find data payload message %q for %q: %w", name, id, err)
	}

	return nil
}

type dataAssetOptions struct {
	referencedDataOptions []ReferencedDataOption
}

// DataAssetOption is an option for validating a DataAsset.
type DataAssetOption func(*dataAssetOptions)

// WithReferencedDataOptions specifies the options for validating ReferencedData within the Data
// Asset.
func WithReferencedDataOptions(referencedDataOptions ...ReferencedDataOption) DataAssetOption {
	return func(opts *dataAssetOptions) {
		opts.referencedDataOptions = referencedDataOptions
	}
}

// WithReferencedDataOption specifies an option for validating ReferencedData within the Data Asset.
func WithReferencedDataOption(referencedDataOption ReferencedDataOption) DataAssetOption {
	return func(opts *dataAssetOptions) {
		opts.referencedDataOptions = append(opts.referencedDataOptions, referencedDataOption)
	}
}

// DataAsset validates a DataAsset.
func DataAsset(da *dapb.DataAsset, options ...DataAssetOption) error {
	opts := &dataAssetOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if da == nil {
		return fmt.Errorf("DataAsset must not be nil")
	}

	m := da.GetMetadata()
	if err := metadatautils.ValidateMetadata(m,
		metadatautils.WithAssetType(atpb.AssetType_ASSET_TYPE_DATA),
		metadatautils.WithInAssetOptions(),
	); err != nil {
		return fmt.Errorf("invalid DataAsset metadata: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(da.GetMetadata().GetIdVersion().GetId())

	if da.GetData() == nil {
		return fmt.Errorf("data payload must be specified for %q", id)
	}

	fds := da.GetFileDescriptorSet()
	if fds == nil {
		return fmt.Errorf("FileDescriptorSet must not be nil for %q", id)
	}
	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return fmt.Errorf("failed to populate registry for %q: %v", id, err)
	}
	if name := da.GetData().MessageName(); name == "" {
		return fmt.Errorf("data payload must not be an empty Any for %q", id)
	} else if _, err := files.FindDescriptorByName(name); err != nil {
		return fmt.Errorf("cannot find data payload message %q for %q: %w", name, id, err)
	}

	// Validate the payload's referenced data.
	if payload, err := utils.ExtractPayload(da); err != nil {
		return err
	} else if _, err := utils.WalkUniqueReferencedData(payload, func(ref *utils.ReferencedDataExt) error {
		if err := ReferencedData(ref, opts.referencedDataOptions...); err != nil {
			return fmt.Errorf("invalid ReferencedData for %q: %w", id, err)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

type referencedDataOptions struct {
	disallowFileReferences bool
}

// ReferencedDataOption is an option for validating a ReferencedData.
type ReferencedDataOption func(*referencedDataOptions)

// WithDisallowFileReferences specifies whether file references are disallowed in the
// ReferencedData.
func WithDisallowFileReferences(disallowFileReferences bool) ReferencedDataOption {
	return func(opts *referencedDataOptions) {
		opts.disallowFileReferences = disallowFileReferences
	}
}

// ReferencedData validates a ReferencedData.
//
// Validation includes:
// - If specified, compare the digest against the referenced data.
// - If the type is a file reference and file references are disallowed, return an error.
func ReferencedData(ref *utils.ReferencedDataExt, options ...ReferencedDataOption) error {
	opts := &referencedDataOptions{}
	for _, opt := range options {
		opt(opts)
	}

	// Type-specific validation.
	switch ref.Type() {
	case utils.FileReferenceType:
		if opts.disallowFileReferences {
			return fmt.Errorf("file references are not allowed (got: %q)", ref.Reference())
		}
	case utils.CASReferenceType:
	case utils.InlinedReferenceType:
		// Nothing to do.
	default:
		return fmt.Errorf("unknown reference type: %d", ref.Type())
	}

	// Validate the digest against the referenced data.
	if ref.Digest() != "" {
		// Get a reader for the data.
		var reader io.Reader
		switch ref.Type() {
		case utils.FileReferenceType:
			file, err := os.Open(ref.Reference())
			if err != nil {
				return fmt.Errorf("failed to open referenced file %q: %w", ref.Reference(), err)
			}
			defer file.Close()
			reader = file
		case utils.CASReferenceType:
			return fmt.Errorf("cannot validate digest of CAS reference %q", ref.Reference())
		case utils.InlinedReferenceType:
			reader = bytes.NewReader(ref.Inlined())
		default:
			return fmt.Errorf("unknown reference type: %d", ref.Type())
		}

		// Test the digest.
		if parsed, err := utils.ParseDigest(ref.Digest()); err != nil {
			return fmt.Errorf("failed to parse digest %q: %w", ref.Digest(), err)
		} else if gotDigest, err := utils.Digest(reader, utils.WithAlgorithm(parsed.Algorithm)); err != nil {
			return fmt.Errorf("failed to compute digest: %w", err)
		} else if gotDigest != parsed.Digest {
			return fmt.Errorf("digest mismatch: got %q, want %q", gotDigest, parsed.Digest)
		}
	}

	return nil
}
