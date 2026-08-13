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

// Package root contains the main entry point for the inbuild CLI.
package root

import (
	"flag"
	"os"

	"intrinsic/production/intrinsic"
	"intrinsic/tools/inbuild/cmd/data/data"
	"intrinsic/tools/inbuild/cmd/httpjson/httpjson"
	"intrinsic/tools/inbuild/cmd/service/service"
	"intrinsic/tools/inbuild/cmd/skill/skill"

	"github.com/spf13/cobra"
)

// RootCmd is the top level command of inbuild.
var RootCmd = &cobra.Command{
	Use:   "inbuild",
	Short: "inbuild builds assets",
	Long:  "inbuild builds assets for Flowstate.",
}

// Inbuild launches the main inbuild CLI.
func Inbuild() {
	intrinsic.Init()

	RootCmd.SetArgs(flag.Args())

	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

// The init function adds subcommands to `inbuild`.
func init() {
	RootCmd.AddCommand(data.DataCmd)
	RootCmd.AddCommand(httpjson.HttpJsonCmd)
	RootCmd.AddCommand(service.ServiceCmd)
	RootCmd.AddCommand(skill.SkillCmd)
}
