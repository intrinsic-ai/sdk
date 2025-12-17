// Copyright 2023 Intrinsic Innovation LLC

// Package bundle contains the entry point for inbuild service bundle.
package bundle

import (
	"fmt"
	"path/filepath"

	"intrinsic/assets/services/build_defs/servicegen"

	"github.com/spf13/cobra"
)

var (
	flagDefaultConfig      string
	flagFileDescriptorSets []string
	flagOciImages          []string
	flagManifest           string
	flagOutput             string
)

// BundleCmd creates service bundles
var BundleCmd *cobra.Command

// Reset global variables so unit tests don't interfere with each other.
func resetBundleCommand() {
	BundleCmd = &cobra.Command{
		Use:   "bundle",
		Short: "Creates service bundles",
		Long:  "Creates service bundles for Flowstate.",
		RunE:  run,
	}

	BundleCmd.Flags().StringVar(&flagDefaultConfig, "default_config", "", "Optional path to default config proto.")
	BundleCmd.Flags().StringArrayVar(&flagFileDescriptorSets, "file_descriptor_set", nil, "Path to binary file descriptor set protos to be used to resolve the configuration message.")
	BundleCmd.Flags().StringArrayVar(&flagOciImages, "oci_image", nil, "Path to tar archive of an OCI image. Must be given at least once, and no more than twice.")
	BundleCmd.Flags().StringVar(&flagManifest, "manifest", "", "Path to a ServiceManifest textproto file.")
	BundleCmd.Flags().StringVar(&flagOutput, "output", "service.bundle.tar", "Path to write service bundle to")
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

func run(cmd *cobra.Command, args []string) error {
	// Validate flags.
	if flagManifest == "" {
		return fmt.Errorf("--manifest is required")
	}
	if flagDefaultConfig != "" && len(flagFileDescriptorSets) == 0 {
		return fmt.Errorf("at least one --file_descriptor_set is required when --default_config is used")
	}
	if len(flagOciImages) == 0 {
		return fmt.Errorf("at least one --oci_image is required")
	}
	if flagOutput == "" {
		return fmt.Errorf("--output must be a valid writable path")
	}

	// Convert file descriptor set paths to absolute paths.
	fileDescriptorSets, err := makeAbsolutePaths(flagFileDescriptorSets)
	if err != nil {
		return err
	}

	// Convert oci image paths to absolute paths.
	ociImages, err := makeAbsolutePaths(flagOciImages)
	if err != nil {
		return err
	}

	return servicegen.CreateServiceBundle(&servicegen.CreateServiceBundleOptions{
		DefaultConfigPath:      flagDefaultConfig,
		FileDescriptorSetPaths: fileDescriptorSets,
		ImageTarPaths:          ociImages,
		ManifestPath:           flagManifest,
		OutputBundlePath:       flagOutput,
	})
}

// The init function establishes command line flags for `inbuild service bundle`
func init() {
	resetBundleCommand()
}
