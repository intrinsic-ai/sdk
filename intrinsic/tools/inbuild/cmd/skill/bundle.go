// Copyright 2023 Intrinsic Innovation LLC

// Package bundle contains the entry point for inbuild skill bundle.
package bundle

import (
	"fmt"

	"intrinsic/skills/skillbundle"
	"intrinsic/skills/skillfix"
	"intrinsic/tools/inbuild/util/skillmanifest"

	"github.com/spf13/cobra"
)

var (
	flagFileDescriptorSet                        string
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

	BundleCmd.Flags().StringVar(&flagFileDescriptorSet, "augmented_file_descriptor_set", "", "Path to an augmented file descriptor set binary proto used to resolve messages referenced in the skill manifest")
	BundleCmd.Flags().StringVar(&flagOciImage, "oci_image", "", "Path to tar archive of an OCI image")
	BundleCmd.Flags().StringVar(&flagManifest, "augmented_manifest", "", "Path to an augmented SkillManifest binary proto")
	BundleCmd.Flags().BoolVar(&flagIncompatibleDisallowManifestDependencies, "incompatible_disallow_manifest_dependencies", false, "Whether to prevent this manifest from declaring dependencies")
	BundleCmd.Flags().StringVar(&flagOutput, "output", "skill.bundle.tar", "Path to write skill bundle to")
}

func run(cmd *cobra.Command, args []string) error {
	// Validate flags.
	if flagManifest == "" {
		return fmt.Errorf("--augmented_manifest is required")
	}
	if flagFileDescriptorSet == "" {
		return fmt.Errorf("--augmented_file_descriptor_set is required")
	}
	if flagOciImage == "" {
		return fmt.Errorf("--oci_image is required")
	}
	if flagOutput == "" {
		return fmt.Errorf("--output must be a valid writable path")
	}

	// Prep the manifest and file descriptor set
	m, fds, err := skillmanifest.LoadManifestAndFileDescriptorSets(cmd.Context(), flagManifest, []string{flagFileDescriptorSet}, flagIncompatibleDisallowManifestDependencies)
	if err != nil {
		return fmt.Errorf("unable to load manifest and file descriptor sets: %v", err)
	}
	if err := skillfix.Manifest(m, skillfix.WithPopulateOldFields(true)); err != nil {
		return fmt.Errorf("unable to make manifest compatible with the latest version of the platform: %v", err)
	}

	// Actually create the skill bundle
	return skillbundle.Write(cmd.Context(), m, flagOutput,
		skillbundle.WithFileDescriptorSet(fds),
		skillbundle.WithImageTarPath(flagOciImage),
	)
}

// The init function establishes command line flags for `inbuild skill bundle`
func init() {
	resetBundleCommand()
}
