// Copyright 2023 Intrinsic Innovation LLC

package cluster

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	fmpb "google.golang.org/protobuf/types/known/fieldmaskpb"
	"intrinsic/assets/baseclientutils"
	clustermanagergrpcpb "intrinsic/frontend/cloud/api/v1/clustermanager_api_go_grpc_proto"
	clustermanagerpb "intrinsic/frontend/cloud/api/v1/clustermanager_api_go_grpc_proto"
	inversiongrpcpb "intrinsic/kubernetes/inversion/v1/inversion_go_grpc_proto"
	inversionpb "intrinsic/kubernetes/inversion/v1/inversion_go_grpc_proto"
	"intrinsic/skills/tools/skill/cmd/dialerutil"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/util/orgutil"
)

var (
	clusterName  string
	rollbackFlag bool
)

// client helps run auth'ed requests for a specific cluster
type client struct {
	tokenSource *auth.ProjectToken
	cluster     string
	project     string
	org         string
	grpcConn    *grpc.ClientConn
	grpcClient  clustermanagergrpcpb.ClustersServiceClient
}

type clusterInfo struct {
	rollback    bool
	mode        string
	state       string
	currentBase string
	currentOS   string
}

// status queries the update status of a cluster
func (c *client) status(ctx context.Context) (*clusterInfo, error) {
	req := clustermanagerpb.GetClusterRequest{
		Project:   c.project,
		Org:       c.org,
		ClusterId: c.cluster,
	}
	cluster, err := c.grpcClient.GetCluster(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("cluster status: %w", err)
	}
	var cp *clustermanagerpb.IPCNode
	for _, n := range cluster.GetIpcNodes() {
		if n.GetIsControlPlane() {
			cp = n
			break
		}
	}
	if cp == nil {
		return nil, fmt.Errorf("control plane not found in cluster list: %q", c.cluster)
	}
	info := &clusterInfo{
		rollback:    cluster.GetRollbackAvailable(),
		mode:        decodeUpdateMode(cluster.GetUpdateMode()),
		state:       decodeUpdateState(cluster.GetUpdateState()),
		currentBase: cluster.GetPlatformVersion(),
		currentOS:   cp.GetOsVersion(),
	}
	return info, nil
}

// setMode runs a request to set the update mode
func (c *client) setMode(ctx context.Context, mode string) error {
	pbm := encodeUpdateMode(mode)
	if pbm == clustermanagerpb.PlatformUpdateMode_PLATFORM_UPDATE_MODE_UNSPECIFIED {
		return fmt.Errorf("invalid mode: %s", mode)
	}
	req := clustermanagerpb.UpdateClusterRequest{
		Project: c.project,
		Org:     c.org,
		Cluster: &clustermanagerpb.Cluster{
			ClusterName: c.cluster,
			UpdateMode:  pbm,
		},
		UpdateMask: &fmpb.FieldMask{Paths: []string{"update_mode"}},
	}
	_, err := c.grpcClient.UpdateCluster(ctx, &req)
	if err != nil {
		return fmt.Errorf("update cluster: %w", err)
	}
	return nil
}

// This is copied from clustermanager.go, but we could diverge from the strings used by
// Inversion if we prefer a different UX.
var updateModeMap = map[string]clustermanagerpb.PlatformUpdateMode{
	"off":              clustermanagerpb.PlatformUpdateMode_PLATFORM_UPDATE_MODE_OFF,
	"on":               clustermanagerpb.PlatformUpdateMode_PLATFORM_UPDATE_MODE_ON,
	"automatic":        clustermanagerpb.PlatformUpdateMode_PLATFORM_UPDATE_MODE_AUTOMATIC,
	"on+accept":        clustermanagerpb.PlatformUpdateMode_PLATFORM_UPDATE_MODE_MANUAL_WITH_ACCEPT,
	"automatic+accept": clustermanagerpb.PlatformUpdateMode_PLATFORM_UPDATE_MODE_AUTOMATIC_WITH_ACCEPT,
}

// encodeUpdateMode encodes a mode string to a proto definition
func encodeUpdateMode(mode string) clustermanagerpb.PlatformUpdateMode {
	return updateModeMap[mode]
}

var updateModeReverseMap map[clustermanagerpb.PlatformUpdateMode]string

func init() {
	updateModeReverseMap = make(map[clustermanagerpb.PlatformUpdateMode]string, len(updateModeMap))
	for k, v := range updateModeMap {
		updateModeReverseMap[v] = k
	}
}

// decodeUpdateMode decodes a mode proto definition into a string
func decodeUpdateMode(mode clustermanagerpb.PlatformUpdateMode) string {
	if m, ok := updateModeReverseMap[mode]; ok {
		return m
	}
	return "unknown"
}

func decodeUpdateState(state clustermanagerpb.UpdateState) string {
	switch state {
	case clustermanagerpb.UpdateState_UPDATE_STATE_UPDATING:
		return "Updating"
	case clustermanagerpb.UpdateState_UPDATE_STATE_PENDING:
		// While we handle this UpdateState it is not actually returned by the backend.
		// It gets translated to UPDATE_STATE_DEPLOYED.
		return "Pending"
	case clustermanagerpb.UpdateState_UPDATE_STATE_FAULT:
		return "Fault"
	case clustermanagerpb.UpdateState_UPDATE_STATE_DEPLOYED:
		return "Deployed"
	// We no longer expose the "Blocked" state to the user.
	// It gets translated to UPDATE_STATE_DEPLOYED.
	default:
		return "Unknown"
	}
}

// getMode runs a request to read the update mode
func (c *client) getMode(ctx context.Context) (string, error) {
	req := clustermanagerpb.GetClusterRequest{
		Project:   c.project,
		Org:       c.org,
		ClusterId: c.cluster,
	}
	cluster, err := c.grpcClient.GetCluster(ctx, &req)
	if err != nil {
		return "", fmt.Errorf("cluster status: %w", err)
	}
	mode := cluster.GetUpdateMode()
	return decodeUpdateMode(mode), nil
}

// run runs an update if one is pending
func (c *client) run(ctx context.Context, rollback bool) error {
	req := clustermanagerpb.SchedulePlatformUpdateRequest{
		Project:    c.project,
		Org:        c.org,
		ClusterId:  c.cluster,
		UpdateType: clustermanagerpb.SchedulePlatformUpdateRequest_UPDATE_TYPE_FORWARD,
	}
	if rollback {
		req.UpdateType = clustermanagerpb.SchedulePlatformUpdateRequest_UPDATE_TYPE_ROLLBACK
	}
	_, err := c.grpcClient.SchedulePlatformUpdate(ctx, &req)
	if err != nil {
		return fmt.Errorf("cluster upgrade run: %w", err)
	}
	return nil
}

func (c *client) close() error {
	if c.grpcConn != nil {
		return c.grpcConn.Close()
	}
	return nil
}

func newTokenSource(project string) (*auth.ProjectToken, error) {
	configuration, err := auth.NewStore().GetConfiguration(project)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, &dialerutil.ErrCredentialsNotFound{
				CredentialName: project,
				Err:            err,
			}
		}
		return nil, fmt.Errorf("get configuration for project %q: %w", project, err)
	}
	token, err := configuration.GetDefaultCredentials()
	if err != nil {
		return nil, fmt.Errorf("get default credentials for project %q: %w", project, err)
	}
	return token, nil
}

func newClient(ctx context.Context, org, project, cluster string) (context.Context, client, error) {
	ts, err := newTokenSource(project)
	if err != nil {
		return nil, client{}, err
	}
	params := dialerutil.DialInfoParams{
		Cluster:  cluster,
		CredName: project,
		CredOrg:  org,
	}
	ctx, conn, err := dialerutil.DialConnectionCtx(ctx, params)
	if err != nil {
		return nil, client{}, fmt.Errorf("create grpc client: %w", err)
	}
	return ctx, client{
		tokenSource: ts,
		cluster:     cluster,
		project:     project,
		org:         org,
		grpcConn:    conn,
		grpcClient:  clustermanagergrpcpb.NewClustersServiceClient(conn),
	}, nil
}

const modeCmdDesc = `
Read/Write the current update mechanism mode

There are 3 modes on the system:

- 'off': no updates can run
- 'on': updates go to the IPC when triggered with inctl or the IPC manager
- 'automatic': updates go to the IPC as soon as they are available

You can add the "+accept" suffix to require acceptance of the update on the
IPC. Acceptance is normally performed through the HMI, although for testing
you can also use "inctl cluster upgrade accept".
`

var modeCmd = &cobra.Command{
	Use:   "mode",
	Short: "Read/Write the current update mechanism mode",
	Long:  modeCmdDesc,
	// at most one arg, the mode
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		projectName := ClusterCmdViper.GetString(orgutil.KeyProject)
		orgName := ClusterCmdViper.GetString(orgutil.KeyOrganization)
		ctx, c, err := newClient(ctx, orgName, projectName, clusterName)
		if err != nil {
			return fmt.Errorf("cluster upgrade client: %w", err)
		}
		defer c.close()
		switch len(args) {
		case 0:
			mode, err := c.getMode(ctx)
			if err != nil {
				return fmt.Errorf("get cluster upgrade mode:\n%w", err)
			}
			fmt.Printf("update mechanism mode: %s\n", mode)
			return nil
		case 1:
			if err := c.setMode(ctx, args[0]); err != nil {
				return fmt.Errorf("set cluster upgrade mode:\n%w", err)
			}
			return nil
		default:
			return fmt.Errorf("invalid number of arguments. At most 1: %d", len(args))
		}
	},
}

const runCmdDesc = `
Run an upgrade of the specified cluster, if new software is available.

This command will execute right away. Please make sure the cluster is safe
and ready to upgrade. It might reboot in the process.
`

// runCmd is the command to execute an update if available
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run an upgrade if available.",
	Long:  runCmdDesc,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		projectName := ClusterCmdViper.GetString(orgutil.KeyProject)
		orgName := ClusterCmdViper.GetString(orgutil.KeyOrganization)
		qOrgName := orgutil.QualifiedOrg(projectName, orgName)
		ctx, c, err := newClient(ctx, orgName, projectName, clusterName)
		if err != nil {
			return fmt.Errorf("cluster upgrade client:\n%w", err)
		}
		defer c.close()
		err = c.run(ctx, rollbackFlag)
		if err != nil {
			return fmt.Errorf("cluster upgrade run:\n%w", err)
		}

		fmt.Printf("update for cluster %q in %q kicked off successfully.\n", clusterName, qOrgName)
		fmt.Printf("monitor running `inctl cluster upgrade --org %s --cluster %s\n`", qOrgName, clusterName)
		return nil
	},
}

// clusterUpgradeCmd is the base command to query the upgrade state
var clusterUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade Intrinsic software on target cluster",
	Long:  "Upgrade Intrinsic software (OS and intrinsic-base) on target cluster.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()

		projectName := ClusterCmdViper.GetString(orgutil.KeyProject)
		orgName := ClusterCmdViper.GetString(orgutil.KeyOrganization)
		ctx, c, err := newClient(ctx, orgName, projectName, clusterName)
		if err != nil {
			return fmt.Errorf("cluster upgrade client:\n%w", err)
		}
		defer c.close()
		ui, err := c.status(ctx)
		if err != nil {
			return fmt.Errorf("cluster status:\n%w", err)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintf(w, "project\tcluster\tmode\tstate\trollback available\tflowstate\tos\n")
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%v\t%s\t%s\n", projectName, clusterName, ui.mode, ui.state, ui.rollback, ui.currentBase, ui.currentOS)
		w.Flush()
		return nil
	},
}

// acceptCmd is the command to accept an update on the IPC in the '+accept' modes.
var acceptCmd = &cobra.Command{
	Use:   "accept",
	Short: "Accept an upgraded Intrinsic software on target cluster",
	Long:  "Accept an upgraded Intrinsic software (OS and intrinsic-base) on target cluster.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()

		consoleIO := bufio.NewReadWriter(
			bufio.NewReader(cmd.InOrStdin()),
			bufio.NewWriter(cmd.OutOrStdout()))

		projectName := ClusterCmdViper.GetString(orgutil.KeyProject)
		orgName := ClusterCmdViper.GetString(orgutil.KeyOrganization)

		ctx, conn, err := newIPCGRPCClient(ctx, projectName, orgName, clusterName)
		if err != nil {
			return fmt.Errorf("cluster upgrade client:\n%w", err)
		}
		defer conn.Close()

		client := inversiongrpcpb.NewIpcUpdaterClient(conn)
		uir, err := client.ReportUpdateInfo(ctx, &inversionpb.GetUpdateInfoRequest{})
		if err != nil {
			return fmt.Errorf("update info request: %w", err)
		}
		if uir.GetState() != inversionpb.UpdateInfo_STATE_UPDATE_AVAILABLE {
			return fmt.Errorf("update not available")
		}

		fmt.Fprintf(consoleIO,
			"Update from %s to %s is available.\nAre you sure you want to accept the update? [y/n] ",
			uir.GetCurrent().GetVersionId(), uir.GetAvailable().GetVersionId())
		consoleIO.Flush()
		response, err := consoleIO.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read response: %w", err)
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" {
			return fmt.Errorf("user did not confirm: %q", response)
		}

		if _, err := client.ApproveUpdate(ctx, &inversionpb.ApproveUpdateRequest{
			Approved: &inversionpb.IntrinsicVersion{
				VersionId: uir.GetAvailable().GetVersionId(),
			},
		}); err != nil {
			return fmt.Errorf("accept update: %w", err)
		}
		return nil
	},
}

func newIPCGRPCClient(ctx context.Context, projectName, orgName, clusterName string) (context.Context, *grpc.ClientConn, error) {
	address := fmt.Sprintf("dns:///www.endpoints.%s.cloud.goog:443", projectName)
	configuration, err := auth.NewStore().GetConfiguration(projectName)
	if err != nil {
		return ctx, nil, fmt.Errorf("credentials not found: %w", err)
	}
	creds, err := configuration.GetDefaultCredentials()
	if err != nil {
		return ctx, nil, fmt.Errorf("get default credentials: %w", err)
	}
	tcOption, err := baseclientutils.GetTransportCredentialsDialOption()
	if err != nil {
		return ctx, nil, fmt.Errorf("cannot retrieve transport credentials: %w", err)
	}
	dialerOpts := append(baseclientutils.BaseDialOptions(),
		grpc.WithPerRPCCredentials(creds),
		tcOption,
	)
	conn, err := grpc.NewClient(address, dialerOpts...)
	if err != nil {
		return ctx, nil, fmt.Errorf("dialing context: %w", err)
	}
	ctx = metadata.AppendToOutgoingContext(ctx, auth.OrgIDHeader, orgName)
	ctx = metadata.AppendToOutgoingContext(ctx, "x-server-name", clusterName)
	return ctx, conn, nil
}

func init() {
	ClusterCmd.AddCommand(clusterUpgradeCmd)
	clusterUpgradeCmd.PersistentFlags().StringVar(&clusterName, "cluster", "", "Name of cluster to upgrade.")
	clusterUpgradeCmd.MarkPersistentFlagRequired("cluster")
	clusterUpgradeCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().BoolVar(&rollbackFlag, "rollback", false, "Whether to trigger a rollback update instead")
	clusterUpgradeCmd.AddCommand(modeCmd)
	clusterUpgradeCmd.AddCommand(acceptCmd)
}
