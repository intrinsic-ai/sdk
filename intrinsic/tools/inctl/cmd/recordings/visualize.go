// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"context"
	"fmt"
	"io"
	"time"

	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/util/cobrautil"
	"intrinsic/tools/inctl/util/color"
	"intrinsic/tools/inctl/util/orgutil"

	"github.com/pborman/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	leaseapigrpcpb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_proto"
	leasepb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_proto"
	replaygrpcpb "intrinsic/logging/proto/replay_service_go_proto"
	replaypb "intrinsic/logging/proto/replay_service_go_proto"

	tpb "google.golang.org/protobuf/types/known/timestamppb"
)

var visualizeCmd = cobrautil.ParentOfNestedSubcommands("visualize", "Visualize Intrinsic recordings")

const (
	serviceTag         string = "inctl"
	leaseRetryInterval        = 10 * time.Second
)

var (
	flagRecordingID string
	flagDuration    string
)

// VisualizeCmdRunner contains the business logic for the visualize command.
// It is a separate struct to allow for dependency injection and easier testing.
type VisualizeCmdRunner struct {
	// A factory function to create a ReplayClient.
	// This is a factory because the client can only be created after a VM is leased.
	ReplayClientFactory func(ctx context.Context, v *viper.Viper, clusterName string) (replaygrpcpb.ReplayClient, io.Closer, error)
	NewLeaseClient      func(cmd *cobra.Command) (leaseapigrpcpb.VMPoolLeaseServiceClient, error)
}

func (r *VisualizeCmdRunner) leaseVM(cmd *cobra.Command, v *viper.Viper, duration time.Duration, leaseClient leaseapigrpcpb.VMPoolLeaseServiceClient) (*leasepb.Lease, error) {
	reservationUUID := uuid.New()
	for { // retry until lease successful.
		expires := time.Now().Add(duration)
		leaseReq := &leasepb.LeaseRequest{Pool: "", Expires: tpb.New(expires), ServiceTag: serviceTag}
		leaseReq.ReservationId = &reservationUUID

		leaseResp, err := leaseClient.Lease(cmd.Context(), leaseReq)
		if err != nil {
			if status.Code(err) == codes.PermissionDenied {
				return nil, fmt.Errorf("visualization host create request failed: %v\n. Your api-key might have expired, run `inctl auth login` to refresh it and retry", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "visualization host create request failed, retrying soon: %v\n", err)
			time.Sleep(leaseRetryInterval)
			continue
		}

		lease := leaseResp.GetLease()
		fmt.Fprintf(cmd.OutOrStdout(), "Visualization host started successfully: %s\n", lease.GetInstance())
		return lease, nil
	}
}

// RunE executes the visualize command.
func (r *VisualizeCmdRunner) RunE(cmd *cobra.Command, _ []string) error {
	v := visualizeCreateParam
	out := cmd.OutOrStdout()

	duration, err := time.ParseDuration(flagDuration)
	if err != nil {
		return fmt.Errorf("Duration '%v' entered is not valid, use something like '30m' or '1h': %v", flagDuration, err)
	}
	leaseClient, err := r.NewLeaseClient(cmd)
	if err != nil {
		return fmt.Errorf("could not create visualization host client: %v", err)
	}

	lease, err := r.leaseVM(cmd, v, duration, leaseClient)
	if err != nil {
		return fmt.Errorf("visualization host creation failed: %v", err)
	}
	clusterName := lease.GetInstance()

	replayClient, closer, err := r.ReplayClientFactory(cmd.Context(), v, clusterName)
	if err != nil {
		return err
	}
	if closer != nil {
		defer closer.Close()
	}

	req := &replaypb.VisualizeRecordingRequest{
		RecordingId: flagRecordingID,
	}

	resp, err := replayClient.VisualizeRecording(cmd.Context(), req)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return fmt.Errorf("%w", err)
		}
		return fmt.Errorf("failed to visualize recording, did you generate it first with `inctl recordings generate`? Error: %w", err)
	}

	fmt.Fprintln(out, "")
	fmt.Fprintln(out, fmt.Sprintf("Visualization created successfully for recording %s", flagRecordingID))
	fmt.Fprintf(out, "- Visualization valid for %s, expires at %s\n", time.Until(lease.GetExpires().AsTime()), lease.GetExpires().AsTime().Format(time.RFC3339))
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Data will load into the visualization over the next few minutes. You will know it is done when data stops appearing in the timeline.")
	color.C.BlueBackground().White().Fprintf(out, "\nLink to visualization: %s", resp.GetUrl())
	fmt.Fprintln(out, "")

	return nil
}

var visualizeCreateParam = viper.New()

// NewVisualizeCmd creates a new cobra command for the visualize create command.
func NewVisualizeCmd(runner *VisualizeCmdRunner) *cobra.Command {
	if runner == nil {
		runner = &VisualizeCmdRunner{
			ReplayClientFactory: func(ctx context.Context, v *viper.Viper, clusterName string) (replaygrpcpb.ReplayClient, io.Closer, error) {
				conn, err := auth.NewCloudConnection(ctx, auth.WithFlagValues(v), auth.WithCluster(clusterName))
				if err != nil {
					return nil, nil, err
				}
				return replaygrpcpb.NewReplayClient(conn), conn, nil
			},
			NewLeaseClient: func(cmd *cobra.Command) (leaseapigrpcpb.VMPoolLeaseServiceClient, error) {
				conn, err := auth.NewCloudConnection(cmd.Context(), auth.WithFlagValues(visualizeCreateParam))
				if err != nil {
					return nil, err
				}
				return leaseapigrpcpb.NewVMPoolLeaseServiceClient(conn), nil
			},
		}
	}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a visualization of a recording in a hosted version of Rerun.io",
		Long:  "Creates a visualization of a recording in a hosted version of Rerun.io",
		Args:  cobra.NoArgs,
		RunE:  runner.RunE,
	}
	cmd.Flags().StringVar(&flagRecordingID, "recording_id", "", "The recording id to visualize.")
	cmd.Flags().StringVarP(&flagDuration, "duration", "d", "", "Desired duration for the visualization to be accessible.")
	cmd.MarkFlagRequired("recording_id")
	cmd.MarkFlagRequired("duration")
	return orgutil.WrapCmd(cmd, visualizeCreateParam, orgutil.WithOrgExistsCheck(func() bool { return checkOrgExists }))
}

var visualizeCreateCmd = NewVisualizeCmd(nil)

func init() {
	RecordingsCmd.AddCommand(visualizeCmd)
	visualizeCmd.AddCommand(visualizeCreateCmd)
	//
	// Until then, create is the only command. It might seem redundant, but it will allow us to avoid
	// having users change what command they call down the line.
}
