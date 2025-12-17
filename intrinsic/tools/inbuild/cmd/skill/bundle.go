// Copyright 2023 Intrinsic Innovation LLC

// Package bundle contains the entry point for inbuild skill bundle.
package bundle

import (
	"fmt"

	"intrinsic/assets/bundleio"
	"intrinsic/tools/inbuild/util/skillmanifest"

	"github.com/spf13/cobra"
)

var (
	flagFileDescriptorSets                       []string
	flagOciImage                                 string
	flagManifest                                 string
	flagIncompatibleDisallowManifestDependencies bool
	flagOutput                                   string
)

// BundleCmd creates skill bundles
var BundleCmd *cobra.Command

// Reset global variables so unit tests don't interfere with each other.
func resetBundleCommand() {
	BundleCmd = &cobra.Command{
		Use:   "bundle",
		Short: "Creates skill bundles",
		Long:  "Creates skill bundles for Flowstate.",
		RunE:  run,
	}

	BundleCmd.Flags().StringArrayVar(&flagFileDescriptorSets, "file_descriptor_set", nil, "Path to binary file descriptor set protos to be used to resolve messages referenced by the skill manifest")
	BundleCmd.Flags().StringVar(&flagOciImage, "oci_image", "", "Path to tar archive of an OCI image")
	BundleCmd.Flags().StringVar(&flagManifest, "manifest", "", "Path to a SkillManifest textproto file")
	BundleCmd.Flags().BoolVar(&flagIncompatibleDisallowManifestDependencies, "incompatible_disallow_manifest_dependencies", false, "Whether to prevent this manifest from declaring dependencies")
	BundleCmd.Flags().StringVar(&flagOutput, "output", "skill.bundle.tar", "Path to write skill bundle to")
}

func run(cmd *cobra.Command, args []string) error {
	// Validate flags.
	if flagManifest == "" {
		return fmt.Errorf("--manifest is required")
	}
	if len(flagFileDescriptorSets) == 0 {
		return fmt.Errorf("at least one --file_descriptor_set is required")
	}
	if flagOciImage == "" {
		return fmt.Errorf("--oci_image is required")
	}
	if flagOutput == "" {
		return fmt.Errorf("--output must be a valid writable path")
	}

	// Prep the manifest and file descriptor set
	m, fds, err := skillmanifest.LoadManifestAndFileDescriptorSets(flagManifest, flagFileDescriptorSets, flagIncompatibleDisallowManifestDependencies)
	if err != nil {
		return fmt.Errorf("unable to load manifest and file descriptor sets: %v", err)
	}

	// Actually create the skill bundle
	return bundleio.WriteSkillBundle(m, flagOutput, &bundleio.WriteSkillBundleOptions{
		Descriptors:  fds,
		ImageTarPath: flagOciImage,
	})
}

// The init function establishes command line flags for `inbuild skill bundle`
func init() {
	resetBundleCommand()
}
