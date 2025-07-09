// Copyright 2023 Intrinsic Innovation LLC

// Package inctldetails provides a reporter for the inctl details.
package inctldetails

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"intrinsic/tools/inctl/cmd/doctor/api/api"
	"intrinsic/tools/inctl/cmd/doctor/reporters/env_vars/envvars"

	rpb "intrinsic/tools/inctl/cmd/doctor/proto/v1/report_go_proto"
)

var (
	// ReporterName is the name of the reporter.
	ReporterName string = "inctl_details"
	// ReportPrefix is the prefix for the report entries.
	ReportPrefix string = "inctl_"
)

var (
	// Reporter is the DiagnosticInformationReporter that reports the environment variables.
	Reporter = api.DiagnosticInformationReporter{
		Name:                               ReporterName,
		Description:                        "Reports details about inctl.",
		GenerateInformation:                generateInformation,
		InformationReporterDependencyNames: []string{envvars.ReporterName},
	}
)

func generateInformation(cmd *cobra.Command, args []string, report *rpb.Report) (*[]*rpb.DiagnosticInformationEntry, error) {
	var entries []*rpb.DiagnosticInformationEntry

	// check for the location of inctl on the $PATH, a la `which inctl`
	whichName := ReportPrefix + "which"
	inctlWhich, err := exec.LookPath("inctl")
	var found bool = false
	var whichValue string
	if err != nil {
		whichValue = err.Error()
	} else {
		whichValue = inctlWhich
		found = true
	}
	which := &rpb.DiagnosticInformationEntry{
		Name:  &whichName,
		Value: &whichValue,
	}
	entries = append(entries, which)

	// check the version of inctl from which, if found
	versionWhichName := ReportPrefix + "which_version"
	var versionWhichValue string
	if found {
		cmd := exec.Command(inctlWhich, "version")
		out, err := cmd.CombinedOutput()
		if err != nil {
			versionWhichValue = fmt.Errorf("failed to run inctl version: %w: %s", err, out).Error()
		}
		versionWhichValue = string(out)
	} else {
		versionWhichValue = "not found"
	}
	versionWhich := &rpb.DiagnosticInformationEntry{
		Name:  &versionWhichName,
		Value: &versionWhichValue,
	}
	entries = append(entries, versionWhich)

	// check for the path of the inctl binary for this execution (how did you run `inctl doctor`?)
	pathName := ReportPrefix + "path"
	inctlPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	pathValue := inctlPath
	path := &rpb.DiagnosticInformationEntry{
		Name:  &pathName,
		Value: &pathValue,
	}
	entries = append(entries, path)

	// get the version of inctl from the binary
	versionName := ReportPrefix + "version"
	var versionValue string
	{
		cmd := exec.Command(inctlPath, "version")
		out, err := cmd.CombinedOutput()
		if err != nil {
			versionValue = fmt.Errorf("failed to run inctl version: %w: %s", err, out).Error()
		}
		versionValue = string(out)
	}
	version := &rpb.DiagnosticInformationEntry{
		Name:  &versionName,
		Value: &versionValue,
	}
	entries = append(entries, version)

	return &entries, nil
}
