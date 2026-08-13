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

// Package generate defines the `inbuild skill generate` command.
package generate

import (
	"intrinsic/tools/inbuild/cmd/skill/generate/config"
	"intrinsic/tools/inbuild/cmd/skill/generate/entrypoint"

	"github.com/spf13/cobra"
)

// GenerateCmd organizes commands for generating code for skills.
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Commands for generating code for skills",
	Long:  "Commands for generating code for skills for Flowstate.",
}

// The init function adds subcommands to `inbuild skill generate`.
func init() {
	GenerateCmd.AddCommand(entrypoint.EntryPointCmd)
	GenerateCmd.AddCommand(config.ConfigCmd)
}
