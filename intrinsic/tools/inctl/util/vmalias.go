// Copyright 2023 Intrinsic Innovation LLC

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
