// Copyright 2023 Intrinsic Innovation LLC

// Package datagen implements creation of the Data Asset bundle.
package datagen

import (
	"fmt"
	"path/filepath"

	"intrinsic/assets/data/databundle"
	"intrinsic/assets/data/datavalidate"
	"intrinsic/assets/data/utils"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/proto/registryutil"

	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	dmpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	mpb "intrinsic/assets/proto/metadata_go_proto"

	"google.golang.org/protobuf/reflect/protodesc"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

// CreateDataBundleOptions provides the data needed to create a Data Asset bundle.
type CreateDataBundleOptions struct {
	// ManifestPath is the path to a DataManifest .textproto file.
	ManifestPath string
	// ReferencedFilePaths is a list of paths to files that are referenced in the Data Asset.
	ReferencedFilePaths []string
	// ExternalReferencedFilePaths is a map specifying the referenced files to exclude from the .tar
	// bundle and the paths to which to remap those references in the payload.
	//
	// Keys are paths to referenced files to exclude.
	// Values are remapped paths for references in the output payload.
	//
	// Excluded files are left out of the .tar bundle and are kept as external references in the
	// payload along with digests to ensure the data are not modified after bundle creation.
	ExternalReferencedFilePaths map[string]string
	// FileDescriptorSetPaths is the paths to binary file descriptor set protos to be used to resolve
	// the data payload message.
	FileDescriptorSetPaths []string
	// OutputBundlePath is the output path for the tar bundle.
	OutputBundlePath string
}

// CreateDataBundle creates a Data Asset bundle on disk.
func CreateDataBundle(opts *CreateDataBundleOptions) error {
	fds, err := registryutil.LoadFileDescriptorSets(opts.FileDescriptorSetPaths)
	if err != nil {
		return fmt.Errorf("failed to load FileDescriptorSets: %w", err)
	}

	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return fmt.Errorf("failed to populate registry: %w", err)
	}
	types, err := registryutil.NewTypesFromFileDescriptorSet(fds)
	if err != nil {
		return fmt.Errorf("failed to populate registry types: %w", err)
	}

	m := &dmpb.DataManifest{}
	if err := protoio.ReadTextProto(opts.ManifestPath, m, protoio.WithResolver(types)); err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}
	if err := datavalidate.DataManifest(m,
		datavalidate.WithFiles(files),
		datavalidate.WithAllowDataManifestRuntimeAssetID(),
	); err != nil {
		return fmt.Errorf("invalid DataManifest: %w", err)
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

	// Resolve relative file path references (that should be relative to the manifest directory).
	manifestDir := filepath.Dir(opts.ManifestPath)
	if payload, err := utils.ExtractPayload(da); err != nil {
		return fmt.Errorf("failed to extract data payload: %w", err)
	} else if payloadOut, err := utils.WalkUniqueReferencedData(payload, func(ref *utils.ReferencedDataExt) error {
		if ref.Type() == utils.FileReferenceType && !filepath.IsAbs(ref.Reference()) {
			ref.SetReference(filepath.Join(manifestDir, ref.Reference()))
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to walk data payload when creating Data Asset bundle at %q: %w", opts.OutputBundlePath, err)
	} else if payloadOutAny, err := anypb.New(payloadOut); err != nil {
		return fmt.Errorf("failed to create Any proto for data payload: %w", err)
	} else {
		da.Data = payloadOutAny
	}

	if err := databundle.Write(da, opts.OutputBundlePath,
		databundle.WithExternalReferencedFilePaths(opts.ExternalReferencedFilePaths),
	); err != nil {
		return fmt.Errorf("failed to write Data Asset bundle: %w", err)
	}

	return nil
}
