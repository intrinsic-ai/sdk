// Copyright 2023 Intrinsic Innovation LLC

// Package reporters contains all the reporters for the doctor command.
package reporters

import (
	"intrinsic/tools/inctl/cmd/doctor/api/api"
	"intrinsic/tools/inctl/cmd/doctor/reporters/cluster/cluster"
	"intrinsic/tools/inctl/cmd/doctor/reporters/env_vars/envvars"
	"intrinsic/tools/inctl/cmd/doctor/reporters/inctl_details/inctldetails"
	"intrinsic/tools/inctl/cmd/doctor/reporters/org/org"
	"intrinsic/tools/inctl/cmd/doctor/reporters/solution/solution"
)

// Reporters is a list of all the reporters, but note that maps are not sorted by insert order, so
// the keys are manually sorted alphabetically wherever used.
var Reporters = map[string]*api.DiagnosticInformationReporter{
	envvars.ReporterName:      &envvars.Reporter,
	inctldetails.ReporterName: &inctldetails.Reporter,
	org.ReporterName:          &org.Reporter,
	cluster.ReporterName:      &cluster.Reporter,
	solution.ReporterName:     &solution.Reporter,
}
