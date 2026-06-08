// Copyright 2023 Intrinsic Innovation LLC

// Package entrypoint defines the `inbuild skill generate entrypoint` command.
package entrypoint

import (
	_ "embed"
	"fmt"
	slice "slices"

	"intrinsic/skills/generator/gen"
	"intrinsic/skills/skillmanifest"
	"intrinsic/util/proto/descriptor"
	"intrinsic/util/proto/protoio"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
)

//go:embed skill_services_provided_to_platform_transitive_set_sci.proto.bin
var providedToPlatformFDSBytes []byte

var serviceVersionsProvidedToPlatform = []smpb.SkillServicesConfig_ServiceVersion{
	smpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_PROJECTOR,
	smpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_EXECUTOR,
	smpb.SkillServicesConfig_INTRINSIC_PROTO_SKILLS_SKILL_INFORMATION,
}

// populateServiceVersions adds the services the skill provides to the platform.
func populateServiceVersions(m *smpb.SkillManifest) error {
	// We do not add anything if the manifest already contains service versions.
	if len(m.GetOptions().GetSkillServicesConfig().GetServiceVersions()) != 0 {
		return nil
	}
	if m.GetOptions() == nil {
		m.Options = &smpb.Options{}
	}
	if m.GetOptions().GetSkillServicesConfig() == nil {
		m.Options.SkillServicesConfig = &smpb.SkillServicesConfig{}
	}
	config := m.GetOptions().GetSkillServicesConfig()
	for _, sv := range serviceVersionsProvidedToPlatform {
		config.ServiceVersions = append(config.ServiceVersions, sv)
	}
	return nil
}

func supportedLanguages() []string {
	return []string{"cpp", "python"}
}

var (
	flagLanguage             string
	flagCcHeader             string
	flagManifest             string
	flagOutput               string
	flagFileDescriptorSet    string
	flagManifestOut          string
	flagFileDescriptorSetOut string
)

// EntryPointCmd creates skill bundles
var EntryPointCmd *cobra.Command

// Reset global variables so unit tests don't interfere with each other.
func resetEntryPointCommand() {
	EntryPointCmd = &cobra.Command{
		Use:   "entrypoint",
		Short: "Generates a skill's entry point",
		Long:  "Generates the main entry point for a skill for Flowstate.",
		RunE:  run,
	}

	EntryPointCmd.Flags().StringVar(&flagLanguage, "language", "", fmt.Sprintf("Language to generate the entry point for. Must be one of: %v", supportedLanguages()))
	EntryPointCmd.Flags().StringVar(&flagCcHeader, "cc_header", "", "Path to header file containing the 'create skill' method referenced in the SkillManifest. Only valid for C++ skills.")
	EntryPointCmd.Flags().StringVar(&flagManifest, "manifest", "", "Path to a SkillManifest binary or text proto file")
	EntryPointCmd.Flags().StringVar(&flagOutput, "output", "", "Path to write the entry point (Default: main.cc or main.py)")
	EntryPointCmd.Flags().StringVar(&flagFileDescriptorSet, "file_descriptor_set", "", "Path to the binary file descriptor set")
	EntryPointCmd.Flags().StringVar(&flagManifestOut, "augmented_manifest_out", "", "Path to write augmented skill manifest binary proto")
	EntryPointCmd.Flags().StringVar(&flagFileDescriptorSetOut, "augmented_file_descriptor_set_out", "", "Path to write augmented file descriptor set binary proto")
}

func run(cmd *cobra.Command, args []string) error {
	// Validate flags.
	if !slice.Contains(supportedLanguages(), flagLanguage) {
		return fmt.Errorf("--language must be one of: %v", supportedLanguages())
	}
	if flagManifest == "" {
		return fmt.Errorf("--manifest is required")
	}
	if flagFileDescriptorSet == "" {
		return fmt.Errorf("--file_descriptor_set is required")
	}
	if flagManifestOut == "" {
		return fmt.Errorf("--augmented_manifest_out is required")
	}
	if flagFileDescriptorSetOut == "" {
		return fmt.Errorf("--augmented_file_descriptor_set_out is required")
	}
	if flagCcHeader == "" && flagLanguage == "cpp" {
		return fmt.Errorf("--cc_header is required when --language=cpp")
	} else if flagCcHeader != "" && flagLanguage != "cpp" {
		return fmt.Errorf("--cc_header is only supported when --language=cpp")
	}

	manifest := new(smpb.SkillManifest)
	if err := protoio.ReadBinaryProto(flagManifest, manifest); err != nil {
		if err := protoio.ReadTextProto(flagManifest, manifest); err != nil {
			return fmt.Errorf("unable to read manifest as binary or text: %v", err)
		}
	}

	fds := &descriptorpb.FileDescriptorSet{}
	if err := protoio.ReadBinaryProto(flagFileDescriptorSet, fds); err != nil {
		return fmt.Errorf("failed to read file descriptor set: %v", err)
	}

	providedToPlatformFDS := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(providedToPlatformFDSBytes, providedToPlatformFDS); err != nil {
		return fmt.Errorf("failed to unmarshal provided to platform file descriptor set: %v", err)
	}
	augmentedFDS, err := descriptor.MergeFileDescriptorSets([]*descriptorpb.FileDescriptorSet{fds, providedToPlatformFDS})
	if err != nil {
		return fmt.Errorf("failed to merge file descriptor sets: %v", err)
	}
	fds = augmentedFDS

	populateServiceVersions(manifest)
	if err := skillmanifest.PruneSourceCodeInfo(manifest, fds); err != nil {
		return fmt.Errorf("failed to prune source code info: %v", err)
	}

	if err := protoio.WriteBinaryProto(flagManifestOut, manifest); err != nil {
		return fmt.Errorf("failed to write augmented skill manifest: %v", err)
	}
	if err := protoio.WriteBinaryProto(flagFileDescriptorSetOut, fds); err != nil {
		return fmt.Errorf("failed to write augmented file descriptor set: %v", err)
	}

	if flagLanguage == "cpp" {
		if flagOutput == "" {
			flagOutput = "main.cc"
		}

		return gen.WriteSkillServiceCC(manifest, []string{flagCcHeader}, flagOutput)
	} else if flagLanguage == "python" {
		if flagOutput == "" {
			flagOutput = "main.py"
		}
		return gen.WriteSkillServicePy(manifest, flagOutput)
	}
	// Unreachable.
	return fmt.Errorf("unreachable return reached; language: %v", flagLanguage)
}

// The init function establishes command line flags for `inbuild skill bundle`
func init() {
	resetEntryPointCommand()
}
