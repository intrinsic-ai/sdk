// Copyright 2023 Intrinsic Innovation LLC

// Package config defines the `inbuild skill generate config` command.
package config

import (
	"fmt"

	"github.com/spf13/cobra"
	"intrinsic/skills/build_defs/skillserviceconfiggen"
	"intrinsic/tools/inbuild/util/skillmanifest"
	"intrinsic/util/proto/protoio"
)

var (
	flagFileDescriptorSets []string
	flagManifest           string
	flagOutput             string
)

// ConfigCmd creates skill bundles
var ConfigCmd *cobra.Command

// Reset global variables so unit tests don't interfere with each other.
func resetConfigCommand() {
	ConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "Generates a skill's config file",
		Long:  "Generates the configuration file used by the entry point for a skill.",
		RunE:  run,
	}

	ConfigCmd.Flags().StringArrayVar(&flagFileDescriptorSets, "file_descriptor_set", nil, "Path to binary file descriptor set protos to be used to resolve messages referenced by the skill manifest")
	ConfigCmd.Flags().StringVar(&flagManifest, "manifest", "", "Path to a SkillManifest textproto file")
	ConfigCmd.Flags().StringVar(&flagOutput, "output", "config.pbbin", "Path to write skill service")
}

func run(cmd *cobra.Command, args []string) error {
	// Validate flags.
	if len(flagFileDescriptorSets) == 0 {
		return fmt.Errorf("at least one --file_descriptor_set is required")
	}
	if flagManifest == "" {
		return fmt.Errorf("--manifest is required")
	}
	if flagOutput == "" {
		return fmt.Errorf("--output must be a valid writable path")
	}

	// Prep the manifest and file descriptor set
	m, fds, err := skillmanifest.LoadManifestAndFileDescriptorSets(flagManifest, flagFileDescriptorSets)
	if err != nil {
		return fmt.Errorf("unable to load manifest and file descriptor sets: %v", err)
	}

	config, err := skillserviceconfiggen.ExtractSkillServiceConfigFromManifest(m, fds)
	if err != nil {
		return fmt.Errorf("unable to generate skill service config: %v", err)
	}

	if err := protoio.WriteBinaryProto(flagOutput, config, protoio.WithDeterministic(true)); err != nil {
		return fmt.Errorf("unable to write skill service config: %v", err)
	}

	return nil
}

// The init function establishes command line flags for `inbuild skill bundle`
func init() {
	resetConfigCommand()
}
