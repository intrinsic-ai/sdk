// Copyright 2023 Intrinsic Innovation LLC

// Package bundle contains the entry point for inbuild data bundle.
package bundle

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"intrinsic/assets/data/build_defs/datagen"
	dapb "intrinsic/assets/data/proto/v1/data_asset_go_proto"
	dmpb "intrinsic/assets/data/proto/v1/data_manifest_go_proto"
	rdspb "intrinsic/assets/data/proto/v1/referenced_data_struct_go_proto"
	"intrinsic/assets/data/utils"
	"intrinsic/assets/referenceddata"
	"intrinsic/util/proto/descriptor"
	"intrinsic/util/proto/registryutil"

	"github.com/spf13/cobra"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
)

var (
	flagFileDescriptorSets     []string
	flagManifest               string
	flagOutput                 string
	flagReferenceToPath        []string
	flagReplaceWithExternalRef []string
)

// BundleCmd creates data bundles
var BundleCmd *cobra.Command

// Reset global variables so unit tests don't interfere with each other.
func resetBundleCommand() {
	BundleCmd = &cobra.Command{
		Use:   "bundle",
		Short: "Creates data bundles",
		Long:  "Creates data bundles for Flowstate.",
		RunE:  run,
	}

	BundleCmd.Flags().StringVar(&flagManifest, "manifest", "", "Path to a DataManifest textproto file.")
	BundleCmd.Flags().StringArrayVar(&flagFileDescriptorSets, "file_descriptor_set", nil, "Path to binary file descriptor set protos used to resolve custom data payload types. Optional when using standard payloads like ReferencedDataStruct.")
	BundleCmd.Flags().StringVar(&flagOutput, "output", "data.bundle.tar", "Path to write data bundle to.")
	BundleCmd.Flags().StringArrayVar(&flagReferenceToPath, "reference_to_path", nil, "Map a file reference in a manifest to a file path relative to the current working directory 'file_reference=path'. Can be repeated. Optional when the reference is a relative path and that path exists relative to the manifest.")
	BundleCmd.Flags().StringArrayVar(&flagReplaceWithExternalRef, "replace_with_external_reference", nil, "Replace all references to a file on disk with an external reference in the payload as 'path=external_ref'. Can be repeated.")
}

func makeAbsolutePaths(paths []string) ([]string, error) {
	var absolutePaths []string
	for _, path := range paths {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		absolutePaths = append(absolutePaths, absPath)
	}
	return absolutePaths, nil
}

func defaultMergeKeys(n int) []string {
	keys := []string{"ReferencedDataStruct builtin", "DataManifest builtin"}
	for i := 2; i < n; i++ {
		keys = append(keys, fmt.Sprintf("FileDescriptorSet %d", i-1))
	}
	return keys
}

// mergeWithBuiltinDescriptors merges a user file descriptor set with built-in descriptors
// (ReferencedDataStruct and DataManifest) into a single FileDescriptorSet in memory.
func mergeWithBuiltinDescriptors(userFds *descriptorpb.FileDescriptorSet) (*descriptorpb.FileDescriptorSet, error) {
	builtinFds := descriptor.FileDescriptorSetFrom(&rdspb.ReferencedDataStruct{})
	manifestFds := descriptor.FileDescriptorSetFrom(&dmpb.DataManifest{})
	fdsToMerge := []*descriptorpb.FileDescriptorSet{builtinFds, manifestFds}
	if userFds != nil {
		fdsToMerge = append(fdsToMerge, userFds)
	}
	fds, err := descriptor.MergeFileDescriptorSets(fdsToMerge, descriptor.WithKeys(defaultMergeKeys(len(fdsToMerge))))
	if err != nil {
		return nil, fmt.Errorf("failed to merge FileDescriptorSets: %w", err)
	}
	return fds, nil
}

// findExternalReferences extracts the data payload from manifest and returns all unique
// file references found within it using the provided FileDescriptorSet to resolve types.
func findExternalReferences(manifest *dmpb.DataManifest, fds *descriptorpb.FileDescriptorSet) ([]string, error) {
	da := &dapb.DataAsset{
		Data:              manifest.GetData(),
		FileDescriptorSet: fds,
	}
	payload, err := utils.ExtractPayload(da)
	if err != nil {
		return nil, fmt.Errorf("failed to extract payload from manifest: %w", err)
	}
	var refs []string
	if _, err := referenceddata.WalkUnique(payload, func(ref *referenceddata.ReferencedData) error {
		if ref.Type() == referenceddata.FileReferenceType {
			refs = append(refs, ref.Reference())
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to walk referenced data: %w", err)
	}
	return refs, nil
}

func run(cmd *cobra.Command, args []string) error {
	if flagManifest == "" {
		return fmt.Errorf("--manifest is required")
	}
	if flagOutput == "" {
		return fmt.Errorf("--output must be a valid writable path")
	}

	externalReferencedFilePaths := make(map[string]string, len(flagReplaceWithExternalRef))
	for _, entry := range flagReplaceWithExternalRef {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --replace_with_external_reference flag %q: expected path=external_ref", entry)
		}
		externalReferencedFilePaths[parts[0]] = parts[1]
	}

	referenceToPathMap := make(map[string]string, len(flagReferenceToPath))
	for _, entry := range flagReferenceToPath {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --reference_to_path flag %q: expected file_reference=path", entry)
		}
		referenceToPathMap[parts[0]] = parts[1]
	}

	fileDescriptorSets, err := makeAbsolutePaths(flagFileDescriptorSets)
	if err != nil {
		return err
	}

	userFds, err := registryutil.LoadFileDescriptorSets(fileDescriptorSets)
	if err != nil {
		return fmt.Errorf("failed to load FileDescriptorSets: %w", err)
	}

	// This lets users omit --file_descriptor_sets when they use a known
	// message like ReferencedDataStruct instead of a custom proto.
	fds, err := mergeWithBuiltinDescriptors(userFds)
	if err != nil {
		return err
	}

	manifest, err := datagen.ReadDataAssetManifest(flagManifest, fds)
	if err != nil {
		return err
	}

	// If a Data Asset manifest references a file by relative path,
	// check if that file exists relative to the manifest. If so,
	// automatically add it to referenceToPathMap as a convenience
	// so that users may omit --reference_to_path for it.
	refs, err := findExternalReferences(manifest, fds)
	if err != nil {
		return err
	}
	manifestDir := filepath.Dir(flagManifest)
	for _, ref := range refs {
		if _, ok := referenceToPathMap[ref]; !ok {
			candidate := filepath.Join(manifestDir, ref)
			if _, err := os.Stat(candidate); err == nil {
				referenceToPathMap[ref] = candidate
			}
		}
	}

	return datagen.CreateDataBundle(cmd.Context(), &datagen.CreateDataBundleOptions{
		Manifest:                    manifest,
		ReferenceToPath:             referenceToPathMap,
		ExternalReferencedFilePaths: externalReferencedFilePaths,
		FileDescriptorSet:           fds,
		OutputBundlePath:            flagOutput,
	})
}

// The init function establishes command line flags for `inbuild data bundle`
func init() {
	resetBundleCommand()
}
