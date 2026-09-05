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

// Package httpjson defines the `inbuild httpjson` command.
package httpjson

import (
	"intrinsic/tools/inbuild/cmd/httpjson/generatemain"

	"github.com/spf13/cobra"
)

// HttpJsonCmd organizes commands for building HTTP / JSON bridges.
var HttpJsonCmd = &cobra.Command{
	Use:   "httpjson",
	Short: "Commands for building HTTP / JSON bridges to gRPC services.",
	Long:  "Commands for building HTTP / JSON bridges to gRPC services.",
}

// The init function adds subcommands to `inbuild httpjson`.
func init() {
	HttpJsonCmd.AddCommand(generatemain.GenerateMainCmd)
}
