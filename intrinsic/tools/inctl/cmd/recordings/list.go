// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"intrinsic/tools/inctl/util/orgutil"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	bmpb "intrinsic/logging/proto/bag_metadata_go_proto"
	grpcpb "intrinsic/logging/proto/bag_packager_service_go_grpc_proto"
	pb "intrinsic/logging/proto/bag_packager_service_go_grpc_proto"
)

var (
	flagStartTimestamp string
	flagEndTimestamp   string
	flagMaxNumResults  uint32
	flagWorkcellName   string
)

var (
	bagStatusToString = map[bmpb.BagStatus_BagStatusEnum]string{
		bmpb.BagStatus_UNSET:                   "0: Status unset",
		bmpb.BagStatus_UPLOAD_PENDING:          "1: Upload pending",
		bmpb.BagStatus_UPLOADING:               "2: Uploading...",
		bmpb.BagStatus_UPLOADED:                "3: Uploaded. Recording file generation pending",
		bmpb.BagStatus_UNCOMPLETABLE:           "4: Uploaded. Recording file generation pending with dropped data",
		bmpb.BagStatus_COMPLETED:               "5: Uploaded. Generated recording file",
		bmpb.BagStatus_UNCOMPLETABLE_COMPLETED: "6: Uploaded. Generated recording file with dropped data",
		bmpb.BagStatus_FAILED:                  "7: Failed",
	}
)

// executeAndPrintListBagsResponse executes the ListBags RPC and prints the response.
func executeAndPrintListBagsResponse(ctx context.Context, client grpcpb.BagPackagerClient, req *pb.ListBagsRequest) (numLines uint32, nextPageCursor []byte, err error) {
	// Execute.
	resp, err := client.ListBags(ctx, req)
	if err != nil {
		return 0, []byte{}, err
	}
	if len(resp.GetBags()) == 0 {
		fmt.Println("No recordings found")
		return 0, []byte{}, nil
	}

	// Print.
	const formatString = "%-22s %-22s %-40s %-55s %-40s"
	lines := []string{
		fmt.Sprintf(formatString, "Start Time", "End Time", "ID", "Status", "Description"),
	}
	numLines = uint32(0)
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
		lines = append(lines,
			fmt.Sprintf(formatString,
				bag.GetBagMetadata().GetStartTime().AsTime().Format(time.RFC3339),
				bag.GetBagMetadata().GetEndTime().AsTime().Format(time.RFC3339),
				bag.GetBagMetadata().GetBagId(),
				statusString,
				description,
			))
		numLines++
	}

	fmt.Println()
	fmt.Println(strings.Join(lines, "\n"))
	return numLines, resp.GetNextPageCursor(), nil
}

var listRecordingsE = func(cmd *cobra.Command, _ []string) error {
	client, err := newBagPackagerClient(cmd.Context())
	if err != nil {
		return err
	}
	var startTime, endTime time.Time
	if flagStartTimestamp != "" {
		startTime, err = time.Parse(time.RFC3339, flagStartTimestamp)
		if err != nil {
			return errors.Wrapf(err, "invalid start timestamp: %s", flagEndTimestamp)
		}
	} else {
		startTime = time.Now().Add(-1000000 * time.Hour)
	}
	if flagEndTimestamp != "" {
		endTime, err = time.Parse(time.RFC3339, flagEndTimestamp)
		if err != nil {
			return errors.Wrapf(err, "invalid end timestamp: %s", flagEndTimestamp)
		}
	} else {
		endTime = time.Now()
	}

	req := &pb.ListBagsRequest{
		OrganizationId: cmdFlags.GetString(orgutil.KeyOrganization),
		MaxNumResults:  &flagMaxNumResults,
		Query: &pb.ListBagsRequest_ListQuery{
			ListQuery: &pb.ListBagsRequest_Query{
				WorkcellName: flagWorkcellName,
				StartTime:    timestamppb.New(startTime),
				EndTime:      timestamppb.New(endTime),
			},
		},
	}
	numLines, nextPageCursor, err := executeAndPrintListBagsResponse(cmd.Context(), client, req)
	if err != nil {
		return err
	}
	if numLines == 0 {
		return nil
	}

	// Handle pagination.
	page := 0
	seenRecordings := uint32(0)
	if len(nextPageCursor) == 0 {
		seenRecordings += numLines
		fmt.Println()
		fmt.Println(fmt.Sprintf("Num recordings: %d", seenRecordings))
		return nil
	}

	for len(nextPageCursor) > 0 {
		page++
		seenRecordings += numLines

		fmt.Println()
		fmt.Println(fmt.Sprintf("Seen pages: %d | Seen recordings: %d", page, seenRecordings))
		fmt.Println()
		fmt.Printf("More results further into the past are available, continue? [Y/n] ")

		// Prompt user if they want to continue.
		var input string
		fmt.Scanln(&input)
		input = strings.ToLower(input)
		if input != "y" && input != "n" && input != "" {
			fmt.Printf("More results further into the past are available, continue? [Y/n] ")
			continue // Invalid input, prompt again.
		}
		if input == "n" {
			return nil
		}

		req := &pb.ListBagsRequest{
			OrganizationId: cmdFlags.GetString(orgutil.KeyOrganization),
			MaxNumResults:  &flagMaxNumResults,
			Query: &pb.ListBagsRequest_Cursor{
				Cursor: nextPageCursor,
			},
		}
		numLines, nextPageCursor, err = executeAndPrintListBagsResponse(cmd.Context(), client, req)
		if err != nil {
			return err
		}
		if numLines == 0 {
			return nil
		}
	}

	return nil
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Lists available recordings for a given workcell",
	Long:    "Lists available recordings for a given workcell",
	Args:    cobra.NoArgs,
	RunE:    listRecordingsE,
}

func init() {
	recordingsCmd.AddCommand(listCmd)
	flags := listCmd.Flags()

	flags.StringVar(&flagWorkcellName, "workcell", "", "The Kubernetes cluster to use.")
	flags.StringVar(&flagStartTimestamp, "start_timestamp", "", "Start timestamp in RFC3339 format for fetching recordings. eg. 2024-08-20T12:00:00Z")
	flags.StringVar(&flagEndTimestamp, "end_timestamp", "", "End timestamp in RFC3339 format for fetching recordings. eg. 2024-08-20T12:00:00Z")
	flags.Uint32Var(&flagMaxNumResults, "max_num_results", 10, "The maximum number of recordings to list per page.")
	listCmd.MarkFlagRequired("workcell")
}
