// Copyright 2023 Intrinsic Innovation LLC

// Package databundle provides utils for working with Process bundles.
package processbundle

import (
	"context"
	"fmt"
	"io"
	"os"

	"intrinsic/assets/ioutils"
	"intrinsic/assets/processes/processmanifest"
	"intrinsic/util/archive/tartooling"

	"github.com/google/safearchive/tar"

	processassetpb "intrinsic/assets/processes/proto/process_asset_go_proto"
	processmanifestpb "intrinsic/assets/processes/proto/process_manifest_go_proto"
	assettypepb "intrinsic/assets/proto/asset_type_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	metadatapb "intrinsic/assets/proto/metadata_go_proto"
)

const (
	processManifestFileName = "process_manifest.binpb"
)

type writeOptions struct {
	writer io.Writer
}

// WriteOption is a functional option for Write.
type WriteOption func(*writeOptions)

// WithWriter specifies the Writer to use instead of creating one for the specified path.
func WithWriter(w io.Writer) WriteOption {
	return func(opts *writeOptions) {
		opts.writer = w
	}
}

// Write writes a Process Asset .tar bundle.
func Write(manifest *processmanifestpb.ProcessManifest, path string, options ...WriteOption) error {
	opts := &writeOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if manifest == nil {
		return fmt.Errorf("ProcessManifest must not be nil")
	}
	err := processmanifest.ValidateProcessManifest(manifest)
	if err != nil {
		return fmt.Errorf("invalid ProcessManifest: %w", err)
	}

	writer := opts.writer
	if writer == nil {
		if path == "" {
			return fmt.Errorf("path must not be empty if a writer is not specified")
		}
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			return fmt.Errorf("failed to open %q for writing: %w", path, err)
		}
		defer f.Close()
		writer = f
	}

	tw := tar.NewWriter(writer)

	if err := tartooling.AddBinaryProto(manifest, tw, processManifestFileName); err != nil {
		return fmt.Errorf("failed write ProcessManifest to bundle: %w", err)
	}

	if err := tw.Close(); err != nil {
		return fmt.Errorf("failed to close tar writer: %w", err)
	}

	return nil
}

// WriteFromAsset writes a Process Asset .tar bundle, given a ProcessAsset.
func WriteFromAsset(pa *processassetpb.ProcessAsset, path string, options ...WriteOption) error {
	if pa == nil {
		return fmt.Errorf("ProcessAsset must not be nil")
	}

	manifest := &processmanifestpb.ProcessManifest{
		Metadata: &processmanifestpb.ProcessMetadata{
			Id:            pa.GetMetadata().GetIdVersion().GetId(),
			DisplayName:   pa.GetMetadata().GetDisplayName(),
			Documentation: pa.GetMetadata().GetDocumentation(),
			Vendor:        pa.GetMetadata().GetVendor(),
			AssetTag:      pa.GetMetadata().GetAssetTag(),
		},
		BehaviorTree: pa.GetBehaviorTree(),
	}

	// Clear the ID version from the Skill metadata in the BehaviorTree. The manifest does not contain
	// a version and the behavior tree on it should not be referencing one either for consistency.
	// This can be seen as the counterpart of [processmanifest.FillInSkillMetadataFromAssetMetadata]
	// in [Process]. The remaining fields of the Skill metadata are assumed to be valid/consistent.
	skill := manifest.GetBehaviorTree().GetDescription()
	if skill != nil {
		skill.IdVersion = ""
	}

	return Write(manifest, path, options...)
}

// ProcessBundle represents a Process Asset bundle.
type ProcessBundle struct {
	Manifest *processmanifestpb.ProcessManifest
}

type readOptions struct {
	reader io.Reader
}

// ReadOption is a functional option for Read.
type ReadOption func(*readOptions)

// WithReader specifies the Reader to use instead of creating one for the specified path.
func WithReader(r io.Reader) ReadOption {
	return func(opts *readOptions) {
		opts.reader = r
	}
}

// Read reads a Process Asset bundle (see Write).
func Read(ctx context.Context, path string, options ...ReadOption) (*ProcessBundle, error) {
	opts := &readOptions{}
	for _, opt := range options {
		opt(opts)
	}

	reader := opts.reader
	if reader == nil {
		if path == "" {
			return nil, fmt.Errorf("path must not be empty if a reader is not specified")
		}
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open %q for reading: %w", path, err)
		}
		defer f.Close()
		reader = f
	}

	// Read single file from the bundle.
	tr := tar.NewReader(reader)
	header, err := tr.Next()
	if err != nil {
		return nil, fmt.Errorf("failed to read first entry of Process bundle: %w", err)
	}
	if header.Typeflag != tar.TypeReg {
		return nil, fmt.Errorf("unexpected entry type in Process bundle: %v", header.Typeflag)
	}
	if header.Name != processManifestFileName {
		return nil, fmt.Errorf("unexpected file in Process bundle: %v", header.Name)
	}

	manifest := &processmanifestpb.ProcessManifest{}
	if err = ioutils.ReadBinaryProto(tr, manifest); err != nil {
		return nil, fmt.Errorf("failed to read ProcessManifest proto in bundle: %w", err)
	}

	// Fail if there are other files in the bundle.
	header, err = tr.Next()
	if err != io.EOF {
		if err != nil {
			return nil, fmt.Errorf("failed to read second entry from Process bundle: %w", err)
		}
		return nil, fmt.Errorf("unexpected second entry in Process bundle: %v", header.Name)
	}

	return &ProcessBundle{
		Manifest: manifest,
	}, nil
}

type processOptions struct {
	readOptions []ReadOption
}

// ProcessOption is a functional option for Process.
type ProcessOption func(*processOptions)

// WithReadOptions provides options to pass to Read.
func WithReadOptions(options ...ReadOption) ProcessOption {
	return func(opts *processOptions) {
		opts.readOptions = options
	}
}

// Process creates a processed Process from a bundle.
func Process(ctx context.Context, path string, options ...ProcessOption) (*processassetpb.ProcessAsset, error) {
	opts := &processOptions{}
	for _, opt := range options {
		opt(opts)
	}

	bundle, err := Read(ctx, path, opts.readOptions...)
	if err != nil {
		return nil, fmt.Errorf("failed to read Process bundle: %w", err)
	}
	manifest := bundle.Manifest

	asset := &processassetpb.ProcessAsset{
		Metadata: &metadatapb.Metadata{
			IdVersion: &idpb.IdVersion{
				Id: manifest.GetMetadata().GetId(),
			},
			DisplayName:   manifest.GetMetadata().GetDisplayName(),
			Documentation: manifest.GetMetadata().GetDocumentation(),
			Vendor:        manifest.GetMetadata().GetVendor(),
			AssetType:     assettypepb.AssetType_ASSET_TYPE_PROCESS,
			AssetTag:      manifest.GetMetadata().GetAssetTag(),
		},
		BehaviorTree: manifest.GetBehaviorTree(),
	}

	// Update the Skill metadata in the BehaviorTree to match the Process asset's metadata. In the
	// manifest the affected fields in the Skill metadata are allowed to be empty but need to be
	// filled in the processed Asset.
	processmanifest.FillInSkillMetadataFromAssetMetadata(
		asset.GetBehaviorTree(), asset.GetMetadata(),
	)

	return asset, nil
}
