// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/prototext"
	pb "intrinsic/logging/proto/bag_packager_service_go_grpc_proto"
)

var (
	flagURL bool
)

var getRecordingE = func(cmd *cobra.Command, _ []string) error {
	client, err := newBagPackagerClient(cmd.Context())
	if err != nil {
		return err
	}
	req := &pb.GetBagRequest{
		BagId:   flagBagID,
		WithUrl: flagURL,
	}
	resp, err := client.GetBag(cmd.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return fmt.Errorf("recording with id %q does not exist", flagBagID)
		}
		return err
	}

	fmt.Print(prototext.Format(resp))
	return nil
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a ROS bag for a given recording id",
	Long:  "Gets a ROS bag for a given recording id",
	Args:  cobra.NoArgs,
	RunE:  getRecordingE,
}

func init() {
	recordingsCmd.AddCommand(getCmd)
	flags := getCmd.Flags()

	flags.StringVar(&flagBagID, "recording_id", "", "The recording id to get ROS bag for.")
	flags.BoolVar(&flagURL, "with_url", false, "If present, generates a signed url to download the bag with.")
}
