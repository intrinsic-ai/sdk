// Copyright 2023 Intrinsic Innovation LLC

// Package cluster envvars provides a reporter for the cluster details.
package cluster

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"intrinsic/tools/inctl/cmd/doctor/api/api"

	rpb "intrinsic/tools/inctl/cmd/doctor/proto/v1/report_go_proto"
)

const (
	// ReporterName is the name of the reporter.
	ReporterName string = "cluster"
	// ReportPrefix is the prefix for the report entries.
	ReportPrefix string = "cluster_"
)

var (
	// Reporter is the DiagnosticInformationReporter that reports cluster details.
	Reporter = api.DiagnosticInformationReporter{
		Name:                               ReporterName,
		Description:                        "Reports the cluster details if a cluster id is provided.",
		GenerateInformation:                generateInformation,
		InformationReporterDependencyNames: []string{"inctl_details", "org", "solution"},
	}
)

func generateInformation(cmd *cobra.Command, args []string, report *rpb.Report) (*[]*rpb.DiagnosticInformationEntry, error) {
	var entries []*rpb.DiagnosticInformationEntry

	// Check if the user has provided a solution id.
	var solutionID string
	for _, entry := range report.GetEntries() {
		if entry.GetName() == "solution_id" {
			solutionID = entry.GetValue()
		}
	}

	clusterID := ""
	if solutionID != "" {
		// If there is a solution id, then use the solution id to get the cluster id.
		for _, entry := range report.GetEntries() {
			if entry.GetName() == "solution_cluster_id" {
				clusterID = entry.GetValue()
			}
		}
	} else {
		// If there is no solution id, check if the user has provided a cluster id.
		clusterIDFromArgs, err := cmd.Flags().GetString("cluster")
		if err != nil {
			return nil, fmt.Errorf("failed to get cluster id: %w", err)
		}
		clusterID = clusterIDFromArgs
	}

	// Either way, add the cluster id to the report, even if it is empty.
	clusterIDName := ReportPrefix + "id"
	clusterIDEntry := &rpb.DiagnosticInformationEntry{
		Name:  &clusterIDName,
		Value: &clusterID,
	}
	entries = append(entries, clusterIDEntry)

	clusterIssueName := ReportPrefix + "issue"
	clusterIssue := "none"

	// If no cluster id is provided, then just return the id, but not the cluster details.
	if clusterID == "" {
		clusterIssue = "no cluster id provided"
		clusterIssueEntry := &rpb.DiagnosticInformationEntry{
			Name:  &clusterIssueName,
			Value: &clusterIssue,
		}
		entries = append(entries, clusterIssueEntry)
		return &entries, nil
	}

	// Get the inctl path from the report.
	var inctlPath string
	for _, entry := range report.GetEntries() {
		if entry.GetName() == "inctl_path" {
			inctlPath = entry.GetValue()
		}
	}
	if inctlPath == "" {
		return nil, fmt.Errorf("inctl_path not found in report")
	}
	// Get the org name from the report.
	var orgName string
	for _, entry := range report.GetEntries() {
		if entry.GetName() == "org_name" {
			orgName = entry.GetValue()
		}
	}
	// If no org name is provided, then just return the id, but not the cluster details.
	if orgName == "" {
		clusterIssue = "no org name provided, cluster details can not be retrieved"
		clusterIssueEntry := &rpb.DiagnosticInformationEntry{
			Name:  &clusterIssueName,
			Value: &clusterIssue,
		}
		entries = append(entries, clusterIssueEntry)
		return &entries, nil
	}

	// Capture the cluster display name and region using the inctl cluster list command.
	clusterListCmd := exec.Command(inctlPath, "cluster", "--org="+orgName, "list", "--output=json")
	out, err := clusterListCmd.CombinedOutput()
	if err != nil {
		clusterIssue = fmt.Sprintf("failed to run 'inctl cluster list' command: %v: %s", err, out)
		clusterIssueEntry := &rpb.DiagnosticInformationEntry{
			Name:  &clusterIssueName,
			Value: &clusterIssue,
		}
		entries = append(entries, clusterIssueEntry)
		return &entries, nil
	}
	var clusterListJSON map[string]any
	if err := json.Unmarshal(out, &clusterListJSON); err != nil {
		clusterIssue = fmt.Sprintf("failed to unmarshal cluster list json: %v", err)
		clusterIssueEntry := &rpb.DiagnosticInformationEntry{
			Name:  &clusterIssueName,
			Value: &clusterIssue,
		}
		entries = append(entries, clusterIssueEntry)
		return &entries, nil
	}
	clusters := clusterListJSON["clusters"].([]any)
	clusterDisplayName := ""
	clusterRegion := ""
	clusterIssue = "not found in cluster list"
	for _, clusterInfoJSON := range clusters {
		clusterInfo := clusterInfoJSON.(map[string]any)
		if clusterInfo["clusterName"] == clusterID {
			clusterDisplayName = clusterInfo["displayName"].(string)
			clusterRegion = clusterInfo["region"].(string)
			clusterIssue = "none"
		}
	}

	clusterDisplayNameName := ReportPrefix + "display_name"
	clusterDisplayNameEntry := &rpb.DiagnosticInformationEntry{
		Name:  &clusterDisplayNameName,
		Value: &clusterDisplayName,
	}
	entries = append(entries, clusterDisplayNameEntry)
	clusterRegionName := ReportPrefix + "region"
	clusterRegionEntry := &rpb.DiagnosticInformationEntry{
		Name:  &clusterRegionName,
		Value: &clusterRegion,
	}
	entries = append(entries, clusterRegionEntry)
	clusterIssueEntry := &rpb.DiagnosticInformationEntry{
		Name:  &clusterIssueName,
		Value: &clusterIssue,
	}
	entries = append(entries, clusterIssueEntry)

	return &entries, nil
}
