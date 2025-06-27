// Copyright 2023 Intrinsic Innovation LLC

package api

import (
	"github.com/spf13/cobra"
	reportpb "intrinsic/tools/inctl/cmd/doctor/proto/v1/report_go_proto"
)

// DiagnosticCheck is an interface for a check that can be run by the doctor command.
type DiagnosticCheck struct {
	Name         string
	Description  string
	ExecuteCheck func(
		cmd *cobra.Command,
		args []string,
		report *reportpb.Report,
	) (*reportpb.DiagnosticCheck, error)
	CheckDependencyNames               []string
	InformationReporterDependencyNames []string
}
