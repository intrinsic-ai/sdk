// Copyright 2023 Intrinsic Innovation LLC

// Package manifest contains the entry point for inbuild skill manifest.
package manifest

import (
	"fmt"

	"intrinsic/tools/inbuild/util/skillmanifest"
	"intrinsic/util/proto/protoio"

	"github.com/spf13/cobra"
)

var (
	flagManifest                                 string
	flagOutput                                   string
	flagFileDescriptorSetOut                     string
	flagFileDescriptorSets                       []string
	flagIncompatibleDisallowManifestDependencies bool
)

// ManifestCmd creates skill manifest binaries from text protos.
var ManifestCmd *cobra.Command

// Reset global variables so unit tests don't interfere with each other.
func resetManifestCommand() {
	ManifestCmd = &cobra.Command{
		Use:   "manifest",
		Short: "Creates skill manifest binaries from text protos",
		Long:  "Creates skill manifest binaries from text protos for Flowstate.",
		RunE:  run,
	}

	ManifestCmd.Flags().StringVar(&flagManifest, "manifest", "", "Path to a SkillManifest textproto file")
	ManifestCmd.Flags().StringVar(&flagOutput, "output", "", "Output path for the skill manifest binary")
	ManifestCmd.Flags().StringVar(&flagFileDescriptorSetOut, "file_descriptor_set_out", "", "Output path for a single file descriptor set")
	ManifestCmd.Flags().StringSliceVar(&flagFileDescriptorSets, "file_descriptor_sets", nil, "Comma separated paths to binary file descriptor set protos")
	ManifestCmd.Flags().BoolVar(&flagIncompatibleDisallowManifestDependencies, "incompatible_disallow_manifest_dependencies", false, "Whether to prevent the skill from declaring dependencies in the manifest")
}

func run(cmd *cobra.Command, args []string) error {
	if flagManifest == "" {
		return fmt.Errorf("--manifest is required")
	}
	if flagOutput == "" {
		return fmt.Errorf("--output is required")
	}
	if flagFileDescriptorSetOut == "" {
		return fmt.Errorf("--file_descriptor_set_out is required")
	}
	if len(flagFileDescriptorSets) == 0 {
		return fmt.Errorf("--file_descriptor_sets is required")
	}

	m, fds, err := skillmanifest.LoadManifestAndFileDescriptorSets(flagManifest, flagFileDescriptorSets, flagIncompatibleDisallowManifestDependencies)
	if err != nil {
		return fmt.Errorf("unable to load manifest and file descriptor sets: %v", err)
	}

	if err := protoio.WriteBinaryProto(flagOutput, m, protoio.WithDeterministic(true)); err != nil {
		return fmt.Errorf("could not write skill manifest proto: %v", err)
	}

	if err := protoio.WriteBinaryProto(flagFileDescriptorSetOut, fds, protoio.WithDeterministic(true)); err != nil {
		return fmt.Errorf("could not write file descriptor set proto: %v", err)
	}

	return nil
}

func init() {
	resetManifestCommand()
}
