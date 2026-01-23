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
	disallowFileReferences bool
}

// DataAssetOption is an option for validating a DataAsset.
type DataAssetOption func(*dataAssetOptions)

// WithDisallowFileReferences specifies whether the Data Asset must not contain ReferencedData with
// file references.
func WithDisallowFileReferences(disallowFileReferences bool) DataAssetOption {
	return func(opts *dataAssetOptions) {
		opts.disallowFileReferences = disallowFileReferences
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

	if opts.disallowFileReferences {
		if payload, err := utils.ExtractPayload(da); err != nil {
			return err
		} else if _, err := utils.WalkUniqueReferencedData(payload, func(ref *utils.ReferencedDataExt) error {
			if ref.Type() == utils.FileReferenceType {
				return fmt.Errorf("file references are not allowed for %q, (got: %q)", id, ref.Reference())
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

// ReferencedData validates a ReferencedData.
//
// Validation includes:
// - If specified, compare the digest against the referenced data.
func ReferencedData(ref *utils.ReferencedDataExt) error {
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
