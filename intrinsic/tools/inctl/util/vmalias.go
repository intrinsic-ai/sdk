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

// Package vmalias provides utilities for working with VM aliases.
package vmalias

import (
	"strings"

)

// IsPoolVM returns true if the given name looks like a pool VM.
func IsPoolVM(name string) bool {
	return strings.HasPrefix(name, "vmp-")
}

// ResolveResult contains the result of a VM or alias resolution.
type ResolveResult struct {
	VM    string
	Alias string
}

// Resolve resolves a VM or alias. Returns the (project, VM) name.
func Resolve(vmOrAlias string) (ResolveResult, error) {
	ret := ResolveResult{
		VM:    vmOrAlias,
		Alias: vmOrAlias,
	}

	return ret, nil
}

// ResolvePrint resolves a VM or alias and prints a warning if the alias does not resolve to the
// expected project. expectedProject may be empty, in which case no warning is printed.
func ResolvePrint(vmOrAlias, expectedProject string) string {
	if IsPoolVM(vmOrAlias) {
		return vmOrAlias
	}
	retVM := vmOrAlias
	return retVM
}
