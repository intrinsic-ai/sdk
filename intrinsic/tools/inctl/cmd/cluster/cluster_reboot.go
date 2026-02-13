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

var rebootDesc = `
Reboot a cluster

Example:
	inctl cluster reboot --cluster <my-cluster> --org <my-org>

If the IPC is online, it will reboot.
`

func rebootCluster(ctx context.Context, conn *grpc.ClientConn, cluster string) error {
	client := clustermanagergrpcpb.NewClustersServiceClient(conn)
	if _, err := client.RebootCluster(
		ctx, &clustermanagerpb.RebootClusterRequest{ClusterName: cluster}); err != nil {
		return fmt.Errorf("request to reboot cluster: %w", err)
	}

	return nil
}

var clusterRebootCmd = &cobra.Command{
	Use:   "reboot --cluster <my-cluster> --org <my-org>",
	Short: "Reboot an IPC",
	Long:  rebootDesc,
	RunE: func(cmd *cobra.Command, argv []string) error {
		ctx := cmd.Context()
		conn, err := newCloudConn(ctx)
		if err != nil {
			return err
		}
		defer conn.Close()

		return rebootCluster(ctx, conn, clusterName)
	},
}

func init() {
	ClusterCmd.AddCommand(clusterRebootCmd)
	clusterRebootCmd.PersistentFlags().StringVar(&clusterName, "cluster", "", "Name of the cluster to use")
}
