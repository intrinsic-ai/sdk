// Copyright 2023 Intrinsic Innovation LLC

// Package org provides a reporter for the organization provided by the user.
package org

import (
	"strings"

	"github.com/spf13/cobra"
	"intrinsic/tools/inctl/cmd/doctor/api/api"
	"intrinsic/tools/inctl/util/orgutil"

	rpb "intrinsic/tools/inctl/cmd/doctor/proto/v1/report_go_proto"
)

const (
	// ReporterName is the name of the reporter.
	ReporterName string = "org"
	// ReportPrefix is the prefix for the report entries.
	ReportPrefix string = "org_"
)

var (
	// Reporter is the DiagnosticInformationReporter that reports the organization details.
	Reporter = api.DiagnosticInformationReporter{
		Name:                               ReporterName,
		Description:                        "Reports the organization details.",
		GenerateInformation:                generateInformation,
		InformationReporterDependencyNames: []string{},
	}
)

func generateInformation(cmd *cobra.Command, args []string, report *rpb.Report) (*[]*rpb.DiagnosticInformationEntry, error) {
	var entries []*rpb.DiagnosticInformationEntry

	orgNameName := ReportPrefix + "name"
	orgNameValue := orgutil.QualifiedOrg(api.CmdFlags.GetFlagProject(), api.CmdFlags.GetFlagOrganization())
	if strings.HasPrefix(orgNameValue, "@") || strings.HasSuffix(orgNameValue, "@") {
		// If either the org or project are missing, make the orgNameValue empty.
		orgNameValue = ""
	}
	orgName := &rpb.DiagnosticInformationEntry{
		Name:  &orgNameName,
		Value: &orgNameValue,
	}
	entries = append(entries, orgName)

	return &entries, nil
}
