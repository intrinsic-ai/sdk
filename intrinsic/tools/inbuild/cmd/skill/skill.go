// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package skill defines the `inbuild skill` command.
package skill

import (
	"intrinsic/tools/inbuild/cmd/skill/bundle"
	"intrinsic/tools/inbuild/cmd/skill/generate/generate"
	"intrinsic/tools/inbuild/cmd/skill/manifest"

	"github.com/spf13/cobra"
)

// SkillCmd organizes commands for building skills.
var SkillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Commands for building skills",
	Long:  "Commands for building skills for Flowstate.",
}

// The init function adds subcommands to `inbuild skill`.
func init() {
	SkillCmd.AddCommand(bundle.BundleCmd)
	SkillCmd.AddCommand(generate.GenerateCmd)
	SkillCmd.AddCommand(manifest.ManifestCmd)
}
