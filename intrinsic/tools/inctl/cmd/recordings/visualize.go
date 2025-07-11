// Copyright 2023 Intrinsic Innovation LLC

package recordings

import (
	"context"
	"fmt"
	"time"

	"github.com/pborman/uuid"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"intrinsic/kubernetes/acl/identity"
	"intrinsic/tools/inctl/util/cobrautil"
	"intrinsic/tools/inctl/util/color"
	utilgrpc "intrinsic/tools/inctl/util/grpc"
	"intrinsic/tools/inctl/util/orgutil"

	tpb "google.golang.org/protobuf/types/known/timestamppb"
	leaseapigrpcpb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_grpc_proto"
	leasepb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_grpc_proto"
	replaygrpcpb "intrinsic/logging/proto/replay_service_go_grpc_proto"
	replaypb "intrinsic/logging/proto/replay_service_go_grpc_proto"
)

var visualizeCmd = cobrautil.ParentOfNestedSubcommands("visualize", "Visualize Intrinsic recordings")

const serviceTag string = "inctl"
const leaseRetryInterval = 10 * time.Second

var (
	flagProject string
	flagOrgID   string

	flagRecordingID string
	flagDuration    string
)

// Checks and populates the project and org flags.
func checkParams(_ *cobra.Command, _ []string) error {
	flagProject = localViper.GetString(orgutil.KeyProject)
	flagOrgID = localViper.GetString(orgutil.KeyOrganization)
	if flagOrgID == "" {
		return fmt.Errorf("--org is required")
	}
	return nil
}

func newLeaseClient(ctx context.Context) (leaseapigrpcpb.VMPoolLeaseServiceClient, error) {
	conn, err := newConn(ctx)
	if err != nil {
		return nil, err
	}
	return leaseapigrpcpb.NewVMPoolLeaseServiceClient(conn), nil
}

func leaseVM(ctx context.Context, project string, orgID string, duration time.Duration) (*leasepb.Lease, error) {
	leaseClient, err := newLeaseClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create visualization host client: %v", err)
	}

	reservationUUID := uuid.New()
	for { // retry until lease successful.
		expires := time.Now().Add(duration)
		leaseReq := &leasepb.LeaseRequest{Pool: "", Expires: tpb.New(expires), ServiceTag: serviceTag}
		leaseReq.ReservationId = &reservationUUID

		leaseResp, err := leaseClient.Lease(ctx, leaseReq)
		if err != nil {
			if status.Code(err) == codes.PermissionDenied {
				return nil, fmt.Errorf("visualization host create request failed: %v\n. Your api-key might have expired, run `inctl auth login` to refresh it and retry", err)
			}
			fmt.Printf("visualization host create request failed, retrying soon: %v\n", err)
			time.Sleep(leaseRetryInterval)
			continue
		}

		lease := leaseResp.GetLease()
		fmt.Printf("Visualization host started successfully: %s\n", lease.GetInstance())
		return lease, nil
	}
}

var visualizeCreateE = func(cmd *cobra.Command, _ []string) error {
	cmdCtx, err := identity.OrgToContext(cmd.Context(), flagOrgID)
	if err != nil {
		return fmt.Errorf("failed to add org info to context: %v", err)
	}
	duration, err := time.ParseDuration(flagDuration)
	if err != nil {
		return fmt.Errorf("Duration '%v' entered is not valid, use something like '30m' or '1h': %v", flagDuration, err)
	}

	lease, err := leaseVM(cmdCtx, flagProject, flagOrgID, duration)
	if err != nil {
		return fmt.Errorf("visualization host creation failed: %v", err)
	}
	clusterName := lease.GetInstance()

	ctx, conn, err := utilgrpc.NewIPCGRPCClient(cmdCtx, flagProject, flagOrgID, clusterName)
	if err != nil {
		return fmt.Errorf("failed to get client for visualization host:\n%w", err)
	}
	defer conn.Close()

	replayClient := replaygrpcpb.NewReplayClient(conn)
	req := &replaypb.VisualizeRecordingRequest{
		RecordingId: flagRecordingID,
	}

	resp, err := replayClient.VisualizeRecording(ctx, req)
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return fmt.Errorf("%w", err)
		}
		return fmt.Errorf("failed to visualize recording, did you generate it first with `inctl recordings generate`? Error: %w", err)
	}

	fmt.Println("")
	fmt.Println(fmt.Sprintf("Visualization created successfully for recording %s", flagRecordingID))
	fmt.Printf("- Visualization valid for %s, expires at %s\n", time.Until(lease.GetExpires().AsTime()), lease.GetExpires().AsTime().Format(time.RFC3339))
	fmt.Println("")
	fmt.Println("Data will load into the visualization over the next few minutes. You will know it is done when data stops appearing in the timeline.")
	color.C.BlueBackground().White().Printf("\nLink to visualization: %s", resp.GetUrl())
	fmt.Println("")

	return nil
}

var visualizeCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Creates a visualization of a recording in a hosted version of Rerun.io",
	Long:    "Creates a visualization of a recording in a hosted version of Rerun.io",
	Args:    cobra.NoArgs,
	RunE:    visualizeCreateE,
	PreRunE: checkParams,
}

func init() {
	recordingsCmd.AddCommand(visualizeCmd)

	visualizeCmd.AddCommand(visualizeCreateCmd)
	visualizeCreateCmd.Flags().StringVar(&flagRecordingID, "recording_id", "", "The recording id to visualize.")
	visualizeCreateCmd.Flags().StringVarP(&flagDuration, "duration", "d", "", "Desired duration for the visualization to be accessible.")
	visualizeCreateCmd.MarkFlagRequired("recording_id")
	visualizeCreateCmd.MarkFlagRequired("duration")
	//
	// Until then, create is the only command. It might seem redundant, but it will allow us to avoid
	// having users change what command they call down the line.
}
