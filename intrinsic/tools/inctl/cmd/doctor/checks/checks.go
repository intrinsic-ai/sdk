// Copyright 2023 Intrinsic Innovation LLC

// Package checks contains the implementation of checks as well as the interface for checks.
package checks

import (
	"intrinsic/tools/inctl/cmd/doctor/api/api"
	"intrinsic/tools/inctl/cmd/doctor/checks/cluster_check/clustercheck"
	"intrinsic/tools/inctl/cmd/doctor/checks/inctl_auth_check/inctlauthcheck"
	"intrinsic/tools/inctl/cmd/doctor/checks/solution_check/solutioncheck"
)

// Checks is a list of all the checks, but note that maps are not sorted by insert order, so
// the keys are manually sorted alphabetically wherever used.
var Checks = map[string]*api.DiagnosticCheck{
	inctlauthcheck.Check.Name: &inctlauthcheck.Check,
	solutioncheck.Check.Name:  &solutioncheck.Check,
	clustercheck.Check.Name:   &clustercheck.Check,
}
