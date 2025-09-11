// Copyright 2023 Intrinsic Innovation LLC

// Package datagen implements creation of the Data asset bundle.
package datagen

import (
	"fmt"
	"path/filepath"

	"intrinsic/assets/bundleio"
	"intrinsic/assets/data/datamanifest"
	"intrinsic/assets/data/utils"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/proto/registryutil"

	anypb "google.golang.org/protobuf/types/known/anypb"
	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	dmpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	mpb "intrinsic/assets/proto/metadata_go_proto"
)

// CreateDataAssetBundleOptions contains the data needed to create a Data asset bundle.
type CreateDataAssetBundleOptions struct {
	// ExcludedReferencedFilePaths is a list of paths to files that should not be included in the tar
	// bundle.
	//
	// Relative paths must be relative to the same base as the output bundle path.
	//
	// These files are left as is and referenced by the Data asset along with a digest to ensure the
	// data are not modified.
	ExcludedReferencedFilePaths []string
	// ExpectedReferencedFilePaths is a list of paths to files that are expected to be referenced in
	// the Data asset.
	//
	// Relative paths must be relative to the same base as the output bundle path.
	ExpectedReferencedFilePaths []string
	// FileDescriptorSetPaths is the paths to binary file descriptor set protos to be used to resolve
	// the data payload message.
	FileDescriptorSetPaths []string
	// Manifest is the path to a DataManifest .textproto file.
	ManifestPath string
	// OutputBundlePath is the output path for the tar bundle.
	OutputBundlePath string
}

// CreateDataAssetBundle creates a Data asset bundle on disk.
func CreateDataAssetBundle(opts CreateDataAssetBundleOptions) error {
	fds, err := registryutil.LoadFileDescriptorSets(opts.FileDescriptorSetPaths)
	if err != nil {
		return fmt.Errorf("cannot build FileDescriptorSet: %w", err)
	}

	types, err := registryutil.NewTypesFromFileDescriptorSet(fds)
	if err != nil {
		return fmt.Errorf("cannot populate registry types: %w", err)
	}

	m := &dmpb.DataManifest{}
	if err := protoio.ReadTextProto(opts.ManifestPath, m, protoio.WithResolver(types)); err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}
	if err := datamanifest.ValidateDataManifest(m,
		datamanifest.WithTypes(types),
	); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}

	da := &dapb.DataAsset{
		Data:              m.GetData(),
		FileDescriptorSet: fds,
		Metadata: &mpb.Metadata{
			AssetType:   atypepb.AssetType_ASSET_TYPE_DATA,
			DisplayName: m.GetMetadata().GetDisplayName(),
			IdVersion: &idpb.IdVersion{
				Id: m.GetMetadata().GetId(),
			},
			Documentation: m.GetMetadata().GetDocumentation(),
			Vendor:        m.GetMetadata().GetVendor(),
		},
	}

	manifestDir := filepath.Dir(opts.ManifestPath)
	outputBundleDir := filepath.Dir(opts.OutputBundlePath)

	// Change relative file path references to be based on the output bundle directory.
	if payload, err := utils.ExtractPayload(da); err != nil {
		return fmt.Errorf("cannot extract data payload: %w", err)
	} else if payloadOut, err := utils.WalkUniqueReferencedData(payload, func(ref *utils.ReferencedDataExt) error {
		if ref.Type() == utils.FileReferenceType && !filepath.IsAbs(ref.Reference()) {
			pathFromManifest := filepath.Join(manifestDir, ref.Reference())
			pathFromOutputBundle, err := filepath.Rel(outputBundleDir, pathFromManifest)
			if err != nil {
				return fmt.Errorf("cannot get relative path: %w", err)
			}
			ref.SetReference(pathFromOutputBundle)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("cannot walk data payload when creating Data asset bundle at %q: %w", opts.OutputBundlePath, err)
	} else if payloadOutAny, err := anypb.New(payloadOut); err != nil {
		return fmt.Errorf("cannot create Any proto for data payload: %w", err)
	} else {
		da.Data = payloadOutAny
	}

	// Get the excluded and expected referenced file paths relative to the output bundle directory.
	excludedReferencedFilePaths := make([]string, len(opts.ExcludedReferencedFilePaths))
	for i, path := range opts.ExcludedReferencedFilePaths {
		if !filepath.IsAbs(path) {
			if excludedReferencedFilePaths[i], err = filepath.Rel(outputBundleDir, path); err != nil {
				return fmt.Errorf("cannot get relative path: %w", err)
			}
		}
	}
	expectedReferencedFilePaths := make([]string, len(opts.ExpectedReferencedFilePaths))
	for i, path := range opts.ExpectedReferencedFilePaths {
		if !filepath.IsAbs(path) {
			if expectedReferencedFilePaths[i], err = filepath.Rel(outputBundleDir, path); err != nil {
				return fmt.Errorf("cannot get relative path: %w", err)
			}
		}
	}

	if err := bundleio.WriteDataAsset(
		da,
		opts.OutputBundlePath,
		bundleio.WithExcludedReferencedFilePaths(excludedReferencedFilePaths),
		bundleio.WithExpectedReferencedFilePaths(expectedReferencedFilePaths),
	); err != nil {
		return fmt.Errorf("cannot write Data asset bundle: %w", err)
	}

	return nil
}
