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

package common

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrConfigStorageNotImplemented = errors.New("line configuration storage is not implemented")
	ErrLineConfigNotFound          = errors.New("line configuration not found")
)

func ExplainConfigStorageNotImplementedError(out io.Writer, command, flag, action string) {
	fmt.Fprintf(
		out,
		`
------------------------------------------------------------------
Got the 'Unimplemented' error from the configuration storage.

It is possible that you are running the %v command in the project where
the configuration storage is not supported yet.

Try running the command with the '--%v' flag.
Note, however, that when this flag is used, the command will not be able
to %v networks that include VMs. It will only be able to %v networks that
consist entirely of workcells.
------------------------------------------------------------------
`,
		command,
		flag,
		action,
		action)
}
