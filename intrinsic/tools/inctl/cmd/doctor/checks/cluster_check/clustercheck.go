// Copyright 2023 Intrinsic Innovation LLC

// Package clustercheck implements a DiagnosticCheck that checks if the user's cluster is running.
package clustercheck

import (
	"fmt"

	"intrinsic/tools/inctl/cmd/doctor/api/api"

	"github.com/spf13/cobra"

	rpb "intrinsic/tools/inctl/cmd/doctor/proto/v1/report_go_proto"
)

const (
	checkName string = "cluster_check"
)

// Check is the DiagnosticCheck that checks if the user's cluster is running.
var Check = api.DiagnosticCheck{
	Name:                               checkName,
	Description:                        "Checks if the user's cluster is running",
	ExecuteCheck:                       checkE,
	CheckDependencyNames:               []string{"solution_check"},
	InformationReporterDependencyNames: []string{"cluster"},
}

func checkE(cmd *cobra.Command, args []string, report *rpb.Report) (*rpb.DiagnosticCheck, error) {
	result := &rpb.DiagnosticCheck{}
	localCheckName := checkName
	result.Name = &localCheckName

	// Get some information from the report.
	var clusterID string
	var clusterIssue string
	for _, entry := range report.GetEntries() {
		if entry.GetName() == "cluster_id" {
			clusterID = entry.GetValue()
		} else if entry.GetName() == "cluster_issue" {
			clusterIssue = entry.GetValue()
		}
	}

	// If the cluster id is not found, then skip the rest of the checks.
	if clusterID == "" {
		resultOutput := fmt.Sprintf(
			"The cluster id is not found in the report, skipping this check. " +
				"Use the --cluster or --solution flags to enable this check.",
		)
		result.Output = &resultOutput
		resultResult := rpb.DiagnosticCheckResult_DIAGNOSTIC_CHECK_RESULT_SKIPPED
		result.Result = &resultResult

		return result, nil
	}

	// If the cluster is not found, then return an error.
	if clusterIssue != "none" {
		resultOutput := fmt.Sprintf(
			"There was an issue getting the cluster details, "+
				"which could be due to the cluster not being found in the given org or "+
				"a different underlying issue: %s",
			clusterIssue,
		)
		result.Output = &resultOutput
		resultResult := rpb.DiagnosticCheckResult_DIAGNOSTIC_CHECK_RESULT_FAILED
		result.Result = &resultResult

		return result, nil
	}

	resultResult := rpb.DiagnosticCheckResult_DIAGNOSTIC_CHECK_RESULT_OK
	result.Result = &resultResult

	return result, nil
}
