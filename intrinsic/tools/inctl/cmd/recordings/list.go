// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"intrinsic/tools/inctl/util/orgutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/timestamppb"

	bmpb "intrinsic/logging/proto/bag_metadata_go_proto"
	pb "intrinsic/logging/proto/bag_packager_service_go_grpc_proto"
)

var (
	flagStartTimestamp string
	flagEndTimestamp   string
	flagMaxNumResults  uint32
	flagWorkcellName   string
)

var bagStatusToString = map[bmpb.BagStatus_BagStatusEnum]string{
	bmpb.BagStatus_UNSET:                   "0: Status unset",
	bmpb.BagStatus_UPLOAD_PENDING:          "1: Upload pending",
	bmpb.BagStatus_UPLOADING:               "2: Uploading...",
	bmpb.BagStatus_UPLOADED:                "3: Uploaded. Recording file generation pending",
	bmpb.BagStatus_UNCOMPLETABLE:           "4: Uploaded. Recording file generation pending with dropped data",
	bmpb.BagStatus_COMPLETED:               "5: Uploaded. Generated recording file",
	bmpb.BagStatus_UNCOMPLETABLE_COMPLETED: "6: Uploaded. Generated recording file with dropped data",
	bmpb.BagStatus_FAILED:                  "7: Failed",
}

// ListCmdRunner manages dependencies for the list command to allow for mocking in tests.
type ListCmdRunner struct {
	NewClient      func(cmd *cobra.Command) (pb.BagPackagerClient, error)
	PromptContinue func(cmd *cobra.Command) (bool, error)
}

// NewListCmd creates a new cobra command for listing recordings.
func NewListCmd(runner *ListCmdRunner) *cobra.Command {
	if runner == nil {
		runner = &ListCmdRunner{
			NewClient: func(cmd *cobra.Command) (pb.BagPackagerClient, error) {
				return newBagPackagerClient(cmd.Context(), listParams)
			},
			PromptContinue: func(cmd *cobra.Command) (bool, error) {
				var input string
				if _, err := fmt.Fscanln(cmd.InOrStdin(), &input); err != nil && err != io.EOF {
					return false, err
				}
				input = strings.ToLower(input)
				return input == "y" || input == "", nil
			},
		}
	}

	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Lists available recordings for a given workcell",
		Long:    "Lists available recordings for a given workcell",
		Args:    cobra.NoArgs,
		RunE:    runner.RunE,
	}

	flags := listCmd.Flags()
	flags.StringVar(&flagWorkcellName, "workcell", "", "The Kubernetes cluster to use.")
	flags.StringVar(&flagStartTimestamp, "start_timestamp", "", "Start timestamp in RFC3339 format for fetching recordings. eg. 2024-08-20T12:00:00Z")
	flags.StringVar(&flagEndTimestamp, "end_timestamp", "", "End timestamp in RFC3339 format for fetching recordings. eg. 2024-08-20T12:00:00Z")
	flags.Uint32Var(&flagMaxNumResults, "max_num_results", 10, "The maximum number of recordings to list per page.")
	listCmd.MarkFlagRequired("workcell")

	return orgutil.WrapCmd(listCmd, listParams, orgutil.WithOrgExistsCheck(func() bool { return checkOrgExists }))
}

func (r *ListCmdRunner) RunE(cmd *cobra.Command, _ []string) error {
	client, err := r.NewClient(cmd)
	if err != nil {
		return err
	}

	startTime, endTime, err := parseTimeFlags()
	if err != nil {
		return err
	}

	req := newListBagsRequest(startTime, endTime)
	return r.executeAndPaginate(cmd, client, req)
}

func parseTimeFlags() (time.Time, time.Time, error) {
	var startTime, endTime time.Time
	var err error

	if flagStartTimestamp != "" {
		startTime, err = time.Parse(time.RFC3339, flagStartTimestamp)
		if err != nil {
			return time.Time{}, time.Time{}, errors.Wrapf(err, "invalid start timestamp: %s", flagStartTimestamp)
		}
	} else {
		startTime = time.Now().Add(-1000000 * time.Hour)
	}

	if flagEndTimestamp != "" {
		endTime, err = time.Parse(time.RFC3339, flagEndTimestamp)
		if err != nil {
			return time.Time{}, time.Time{}, errors.Wrapf(err, "invalid end timestamp: %s", flagEndTimestamp)
		}
	} else {
		endTime = time.Now()
	}

	return startTime, endTime, nil
}

func newListBagsRequest(startTime, endTime time.Time) *pb.ListBagsRequest {
	return &pb.ListBagsRequest{
		OrganizationId: listParams.GetString(orgutil.KeyOrganization),
		MaxNumResults:  &flagMaxNumResults,
		Query: &pb.ListBagsRequest_ListQuery{
			ListQuery: &pb.ListBagsRequest_Query{
				WorkcellName: flagWorkcellName,
				StartTime:    timestamppb.New(startTime),
				EndTime:      timestamppb.New(endTime),
			},
		},
	}
}

func (r *ListCmdRunner) executeAndPaginate(cmd *cobra.Command, client pb.BagPackagerClient, req *pb.ListBagsRequest) error {
	var nextPageCursor []byte
	var err error
	var numLines uint32

	numLines, nextPageCursor, err = r.executeAndPrintListBagsResponse(cmd, client, req)
	if err != nil {
		return err
	}

	page := 0
	seenRecordings := numLines
	for len(nextPageCursor) > 0 {
		page++
		fmt.Fprintf(cmd.OutOrStdout(), "\nSeen pages: %d | Seen recordings: %d\n", page, seenRecordings)
		fmt.Fprint(cmd.OutOrStdout(), "\nMore results further into the past are available, continue? [Y/n] ")

		cont, err := r.PromptContinue(cmd)
		if err != nil {
			return err
		}
		if !cont {
			break
		}

		req.Query = &pb.ListBagsRequest_Cursor{Cursor: nextPageCursor}
		numLines, nextPageCursor, err = r.executeAndPrintListBagsResponse(cmd, client, req)
		if err != nil {
			return err
		}
		seenRecordings += numLines
	}

	fmt.Fprintf(cmd.OutOrStdout(), "\nNum recordings: %d\n", seenRecordings)

	return nil
}

func (r *ListCmdRunner) executeAndPrintListBagsResponse(cmd *cobra.Command, client pb.BagPackagerClient, req *pb.ListBagsRequest) (uint32, []byte, error) {
	resp, err := client.ListBags(cmd.Context(), req)
	if err != nil {
		return 0, nil, err
	}
	if len(resp.GetBags()) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No recordings found")
		return 0, nil, nil
	}

	var out bytes.Buffer
	const formatString = "%-22s %-22s %-40s %-55s %-40s"
	fmt.Fprintln(&out, "")
	fmt.Fprintf(&out, formatString, "Start Time", "End Time", "ID", "Status", "Description")
	fmt.Fprintln(&out)

	for _, bag := range resp.GetBags() {
		description := bag.GetBagMetadata().GetDescription()
		if description == "" {
			description = "<NO-DESCRIPTION>"
		}
		status := bag.GetBagMetadata().GetStatus().GetStatus()
		statusString, ok := bagStatusToString[status]
		if !ok {
			statusString = status.String()
		}
		fmt.Fprintf(&out, formatString,
			bag.GetBagMetadata().GetStartTime().AsTime().Format(time.RFC3339),
			bag.GetBagMetadata().GetEndTime().AsTime().Format(time.RFC3339),
			bag.GetBagMetadata().GetBagId(),
			statusString,
			description,
		)
		fmt.Fprintln(&out)
	}

	fmt.Fprint(cmd.OutOrStdout(), out.String())
	return uint32(len(resp.GetBags())), resp.GetNextPageCursor(), nil
}

var listParams = viper.New()

func init() {
	listCmd := NewListCmd(nil)
	RecordingsCmd.AddCommand(listCmd)
}
