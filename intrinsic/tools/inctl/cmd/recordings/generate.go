// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"intrinsic/tools/inctl/util/orgutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	pb "intrinsic/logging/proto/bag_packager_service_go_grpc_proto"
)

// Number of times to status-check the recording generation after encountering a 504 timeout error
// on the initial GenerateBag call.
//
// Timeouts are from nginx, not client-side, as deadlines are infinite by default for gRPC clients:
// https://grpc.io/docs/guides/deadlines/#deadlines-on-the-client
//
// This is needed because the GenerateBag call can take a long time to complete.
const (
	maxPostTimeoutRetries = 10
	postTimeoutRetryDelay = 30 * time.Second
)

const (
	numBytesInMB           = 1024 * 1024
	largeRecordingByteSize = 50 * numBytesInMB
)

// GenerateCmdRunner manages dependencies for the generate command to allow for mocking in tests.
type GenerateCmdRunner struct {
	NewClient func(cmd *cobra.Command) (pb.BagPackagerClient, error)
}

// NewGenerateCmd creates a new cobra command for generating recordings.
func NewGenerateCmd(runner *GenerateCmdRunner) *cobra.Command {
	if runner == nil {
		runner = &GenerateCmdRunner{
			NewClient: func(cmd *cobra.Command) (pb.BagPackagerClient, error) {
				return newBagPackagerClient(cmd.Context(), generateParam)
			},
		}
	}

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generates an Intrinsic recording file for a given recording id",
		Long:  "Generates an Intrinsic recording file for a given recording id",
		Args:  cobra.NoArgs,
		RunE:  runner.RunE,
	}

	flags := generateCmd.Flags()
	flags.StringVar(&flagBagID, "recording_id", "", "The recording id to generate Intrinsic recording file for.")
	generateCmd.MarkFlagRequired("recording_id")

	return orgutil.WrapCmd(generateCmd, generateParam, orgutil.WithOrgExistsCheck(func() bool { return checkOrgExists }))
}

func (r *GenerateCmdRunner) RunE(cmd *cobra.Command, _ []string) error {
	client, err := r.NewClient(cmd)
	if err != nil {
		return err
	}

	// Fetch to validate.
	getReq := &pb.GetBagRequest{
		BagId:   flagBagID,
		WithUrl: false,
	}
	getResp, err := client.GetBag(cmd.Context(), getReq)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return fmt.Errorf("recording with id \"%s\" does not exist", flagBagID)
		}
		return err
	}
	if getResp.GetBag().GetBagFile() != nil {
		return fmt.Errorf("recording with id \"%s\" is already generated", flagBagID)
	}

	// Generate.
	recordingByteSize := getResp.GetBag().GetBagMetadata().GetTotalBytes()
	if recordingByteSize > largeRecordingByteSize {
		fmt.Fprintln(cmd.OutOrStdout(), "")
		fmt.Fprintln(cmd.OutOrStdout(), "WARNING:")
		fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf("  Recording with id \"%s\" is large (%d MB) and might take several minutes (usually up to ~15 minutes) to generate...", flagBagID, recordingByteSize/numBytesInMB))
		fmt.Fprintln(cmd.OutOrStdout(), "  Please wait and do NOT close this terminal or attempt to generate the recording again, the server will continue processing the request.")
		fmt.Fprintln(cmd.OutOrStdout(), "")
	}

	fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf("Starting generation of recording with id \"%s\"...", flagBagID))

	generateReq := &pb.GenerateBagRequest{
		Query: &pb.GenerateBagRequest_BagId{
			BagId: flagBagID,
		},
		OrganizationId: generateParam.GetString(orgutil.KeyOrganization),
	}
	_, err = client.GenerateBag(cmd.Context(), generateReq)
	if err != nil {
		// A server timeout is expected if the recording is large, this is usually not an error.
		//
		// It usually means that the server is still processing the request, so we should GetBag until
		// we see the file or timeout.
		if !strings.Contains(err.Error(), "504") {
			return err
		}

		for i := 0; i < maxPostTimeoutRetries; i++ {
			fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf("Still generating%s", strings.Repeat(".", i)))

			getResp, err := client.GetBag(cmd.Context(), getReq)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf("Failed to get recording with id \"%s\" to check generation status, server might still be processing: %v", flagBagID, err))
			}

			if getResp.GetBag().GetBagFile() != nil {
				break
			}

			if i == maxPostTimeoutRetries-1 {
				return fmt.Errorf("failed to generate recording with id \"%s\" after %d retries, try generating again or waiting longer for the recording to be generated", flagBagID, maxPostTimeoutRetries)
			}
			time.Sleep(postTimeoutRetryDelay + time.Duration(rand.Float32()*5.0)*time.Second)
		}
	}

	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf("Generated recording file for recording ID %s", flagBagID))
	fmt.Fprintln(cmd.OutOrStdout(), "")
	return nil
}

var generateParam = viper.New()

func init() {
	RecordingsCmd.AddCommand(NewGenerateCmd(nil))
}
