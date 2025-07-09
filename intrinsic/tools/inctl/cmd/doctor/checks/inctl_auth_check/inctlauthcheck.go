// Copyright 2023 Intrinsic Innovation LLC

// Package inctlauthcheck implements a DiagnosticCheck that checks if the user has authenticated
// with inctl.
package inctlauthcheck

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"intrinsic/tools/inctl/cmd/doctor/api/api"

	rpb "intrinsic/tools/inctl/cmd/doctor/proto/v1/report_go_proto"
)

const (
	checkName string = "inctl_auth"
)

// Check is the DiagnosticCheck that checks if the user has authenticated with inctl.
var Check = api.DiagnosticCheck{
	Name:                               checkName,
	Description:                        "Checks if the user has authenticated with inctl",
	ExecuteCheck:                       checkE,
	CheckDependencyNames:               []string{},
	InformationReporterDependencyNames: []string{},
}

func checkE(cmd *cobra.Command, args []string, report *rpb.Report) (*rpb.DiagnosticCheck, error) {
	result := &rpb.DiagnosticCheck{}
	localCheckName := checkName
	result.Name = &localCheckName
	var inctlPath string
	for _, entry := range report.GetEntries() {
		if entry.GetName() == "inctl_path" {
			inctlPath = entry.GetValue()
			break
		}
	}
	if inctlPath == "" {
		return nil, fmt.Errorf("inctl_path not found in report")
	}
	cmdArgs := exec.Command(inctlPath, "auth", "list")
	out, err := cmdArgs.CombinedOutput()
	outStr := string(out)
	if err != nil {
		resultOutput := fmt.Errorf("failed to run 'inctl auth list': %w: %s", err, out).Error()
		result.Output = &resultOutput
		resultResult := rpb.DiagnosticCheckResult_DIAGNOSTIC_CHECK_RESULT_FAILED
		result.Result = &resultResult
	} else {
		result.Output = &outStr
		if strings.Contains(outStr, "The following organizations can be used:") {
			resultResult := rpb.DiagnosticCheckResult_DIAGNOSTIC_CHECK_RESULT_OK
			result.Result = &resultResult
		} else {
			resultResult := rpb.DiagnosticCheckResult_DIAGNOSTIC_CHECK_RESULT_WARNING
			result.Result = &resultResult
			detailName := "warning_reason"
			detailValue := "output of inctl auth list was empty, which may be a problem"
			detail := &rpb.DiagnosticCheckDetail{
				Name:  &detailName,
				Value: &detailValue,
			}
			result.Details = append(result.GetDetails(), detail)
		}
	}
	return result, nil
}
