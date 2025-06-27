// Copyright 2023 Intrinsic Innovation LLC

package api

import (
	"github.com/spf13/cobra"
	reportpb "intrinsic/tools/inctl/cmd/doctor/proto/v1/report_go_proto"
)

// DiagnosticInformationReporter is an interface for a class that generates information entry for the report.
type DiagnosticInformationReporter struct {
	Name                string
	Description         string
	GenerateInformation func(
		cmd *cobra.Command,
		args []string,
		report *reportpb.Report,
	) (*[]*reportpb.DiagnosticInformationEntry, error)
	InformationReporterDependencyNames []string
}
