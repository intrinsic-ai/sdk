// Copyright 2023 Intrinsic Innovation LLC

package cluster

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	clustermanagergrpcpb "intrinsic/frontend/cloud/api/v1/clustermanager_api_go_proto"
	clustermanagerpb "intrinsic/frontend/cloud/api/v1/clustermanager_api_go_proto"
)

var (
	deleteDesc = `
Delete an IPC.

Example:
	inctl cluster delete <cluster-name> --org <my-org>

If the IPC is online, it will be reset back to the unregistered state. If it
is offline, it will be removed from your organization, but will need to be
reset or reinstalled before it is registered to another project.
`
	require_reachable = false
)

func deleteCluster(ctx context.Context, conn *grpc.ClientConn, cluster string) error {
	client := clustermanagergrpcpb.NewClustersServiceClient(conn)
	if _, err := client.DeleteCluster(
		ctx, &clustermanagerpb.DeleteClusterRequest{ClusterName: cluster, RequireReachable: require_reachable}); err != nil {
		return fmt.Errorf("request to delete cluster: %w", err)
	}

	return nil
}

var clusterDeleteCmd = &cobra.Command{
	Use:   "delete my-cluster --org <my-org>",
	Short: "Delete an IPC",
	Long:  deleteDesc,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, argv []string) error {
		ctx := cmd.Context()
		conn, err := newCloudConn(ctx)
		if err != nil {
			return err
		}
		defer conn.Close()

		return deleteCluster(ctx, conn, argv[0])
	},
}

func init() {
	ClusterCmd.AddCommand(clusterDeleteCmd)

	clusterDeleteCmd.Flags().BoolVar(&require_reachable, "require_reachable", false, "Fail the deletion command when the cluster was not reachable")
}
