// Copyright 2023 Intrinsic Innovation LLC

// Package entrypoint defines the `inbuild skill generate entrypoint` command.
package entrypoint

import (
	"fmt"
	slice "slices"

	"github.com/spf13/cobra"
	"intrinsic/skills/generator/gen"
	smpb "intrinsic/skills/proto/skill_manifest_go_proto"
	"intrinsic/util/proto/protoio"
)

func supportedLanguages() []string {
	return []string{"cpp", "python"}
}

var (
	flagLanguage string
	flagCcHeader string
	flagManifest string
	flagOutput   string
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
	EntryPointCmd.Flags().StringVar(&flagManifest, "manifest", "", "Path to a SkillManifest textproto file")
	EntryPointCmd.Flags().StringVar(&flagOutput, "output", "", "Path to write the entry point (Default: main.cc or main.py)")
}

func run(cmd *cobra.Command, args []string) error {
	// Validate flags.
	if !slice.Contains(supportedLanguages(), flagLanguage) {
		return fmt.Errorf("--language must be one of: %v", supportedLanguages())
	}
	if flagManifest == "" {
		return fmt.Errorf("--manifest is required")
	}
	if flagCcHeader == "" && flagLanguage == "cpp" {
		return fmt.Errorf("--cc_header is required when --language=cpp")
	} else if flagCcHeader != "" && flagLanguage != "cpp" {
		return fmt.Errorf("--cc_header is only supported when --language=cpp")
	}

	manifest := new(smpb.SkillManifest)
	if err := protoio.ReadTextProto(flagManifest, manifest); err != nil {
		return fmt.Errorf("unable to read manifest: %v", err)
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
