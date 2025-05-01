// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	pb "intrinsic/logging/proto/bag_packager_service_go_grpc_proto"
	"intrinsic/tools/inctl/util/orgutil"
)

var (
	flagStartTimestamp string
	flagEndTimestamp   string
	flagWorkcellName   string
)

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
		WorkcellName:   flagWorkcellName,
		StartTime:      timestamppb.New(startTime),
		EndTime:        timestamppb.New(endTime),
	}
	resp, err := client.ListBags(cmd.Context(), req)
	if err != nil {
		return err
	}
	if len(resp.GetBags()) == 0 {
		return nil
	}
	const formatString = "%-40s %-30s %-45s %-45s %-45s"
	lines := []string{
		fmt.Sprintf(formatString, "Description", "Status", "ID", "Start Time", "End Time"),
	}

	// Sort the response by start time.
	sort.Slice(resp.GetBags(), func(i, j int) bool {
		return resp.GetBags()[i].GetBagMetadata().GetStartTime().AsTime().Before(resp.GetBags()[j].GetBagMetadata().GetStartTime().AsTime())
	})
	for _, bag := range resp.GetBags() {
		description := bag.GetBagMetadata().GetDescription()
		if description == "" {
			description = "<NO-DESCRIPTION>"
		}
		lines = append(lines,
			fmt.Sprintf(formatString,
				description,
				bag.GetBagMetadata().GetStatus().GetStatus().String(),
				bag.GetBagMetadata().GetBagId(),
				bag.GetBagMetadata().GetStartTime().AsTime(),
				bag.GetBagMetadata().GetEndTime().AsTime()))
	}
	fmt.Println(strings.Join(lines, "\n"))

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
	listCmd.MarkFlagRequired("workcell")
}
