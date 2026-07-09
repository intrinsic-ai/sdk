// Copyright 2023 Intrinsic Innovation LLC

// Package datagen implements creation of the Data Asset bundle.
package datagen

import (
	"context"
	"fmt"
	"path/filepath"

	"intrinsic/assets/data/databundle"
	"intrinsic/assets/data/datavalidate"
	"intrinsic/assets/data/utils"
	"intrinsic/assets/referenceddata"
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
	// ReferenceToPath is a map of file references as specified in a manifest to file paths
	// relative to the current working directory. Data Asset manifests specify file references
	// as paths relative to the manifest. However, the intrinsic_data Bazel rule invokes the
	// datagen tool from a Bazel action. Bazel will run the tool in that action's execroot,
	// and it has its own conventions about where in the execroot it places files. This mapping
	// allows people to specify paths in the manifest that make sense in the source tree, while
	// letting intrinsic_data tell datagen where Bazel actually put the files when the Bazel
	// action runs this tool.
	ReferenceToPath map[string]string
	// ExternalReferencedFilePaths is a map specifying files on disk to exclude from
	// the .tar bundle, and external references to replace those files with.
	//
	// Keys are paths to files relative to the current working directory.
	// Values are external references to put in the Data Asset instead of those files.
	//
	// Files excluded here are left out of the .tar bundle and are kept as external references in the
	// payload along with digests to ensure the data are not modified after bundle creation.
	ExternalReferencedFilePaths map[string]string
	// FileDescriptorSetPaths is the paths to binary file descriptor set protos to be used to resolve
	// the data payload message.
	FileDescriptorSetPaths []string
	// OutputBundlePath is the output path for the tar bundle.
	OutputBundlePath string
}

// CreateDataBundle creates a Data Asset bundle on disk.
func CreateDataBundle(ctx context.Context, opts *CreateDataBundleOptions) error {
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
	if err := datavalidate.DataManifest(ctx, m,
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

	if payload, err := utils.ExtractPayload(da); err != nil {
		return fmt.Errorf("failed to extract data payload: %w", err)
	} else if payloadOut, err := referenceddata.WalkUnique(payload, func(ref *referenceddata.ReferencedData) error {
		if ref.Type() == referenceddata.FileReferenceType && !filepath.IsAbs(ref.Reference()) {
			// Translate file references to their actual locations on disk
			execPath, ok := opts.ReferenceToPath[ref.Reference()]
			if !ok {
				return fmt.Errorf("referenced file %q not found in ReferenceToPath; check that the file is included in the data attribute of intrinsic_data", ref.Reference())
			}
			ref.SetReference(execPath)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to walk data payload when creating Data Asset bundle at %q: %w", opts.OutputBundlePath, err)
	} else if payloadOutAny, err := anypb.New(payloadOut); err != nil {
		return fmt.Errorf("failed to create Any proto for data payload: %w", err)
	} else {
		da.Data = payloadOutAny
	}

	if err := databundle.WriteFile(ctx, da, opts.OutputBundlePath,
		databundle.WithExternalReferencedFilePaths(opts.ExternalReferencedFilePaths),
	); err != nil {
		return fmt.Errorf("failed to write Data Asset bundle: %w", err)
	}

	return nil
}
