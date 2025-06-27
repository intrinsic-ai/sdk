// Copyright 2023 Intrinsic Innovation LLC

// Package envvars provides a reporter for the environment variables.
package envvars

import (
	"os"

	"github.com/spf13/cobra"
	"intrinsic/tools/inctl/cmd/doctor/api/api"
	rpb "intrinsic/tools/inctl/cmd/doctor/proto/v1/report_go_proto"
)

const (
	// ReporterName is the name of the reporter.
	ReporterName string = "env_vars"
	// ReportPrefix is the prefix for the report entries.
	ReportPrefix string = "env_"
)

var (
	// Reporter is the DiagnosticInformationReporter that reports the environment variables.
	Reporter = api.DiagnosticInformationReporter{
		Name:                               ReporterName,
		Description:                        "Reports the environment variables.",
		GenerateInformation:                generateInformation,
		InformationReporterDependencyNames: []string{},
	}
)

func generateInformation(cmd *cobra.Command, args []string, report *rpb.Report) (*[]*rpb.DiagnosticInformationEntry, error) {
	var entries []*rpb.DiagnosticInformationEntry
	envVars := []string{
		"INTRINSIC_ORG",
		"INTRINSIC_PROJECT",
	}
	for _, env := range envVars {
		name := ReportPrefix + env
		value := os.Getenv(env)
		entry := &rpb.DiagnosticInformationEntry{
			Name:  &name,
			Value: &value,
		}
		entries = append(entries, entry)
	}
	return &entries, nil
}
