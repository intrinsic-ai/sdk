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

var poweroffDesc = `
Power off a cluster

Example:
	inctl cluster poweroff --cluster <my-cluster> --org <my-org>

If the IPC is online, it will power off.
`

func poweroffCluster(ctx context.Context, conn *grpc.ClientConn, cluster string) error {
	client := clustermanagergrpcpb.NewClustersServiceClient(conn)
	if _, err := client.PoweroffCluster(
		ctx, &clustermanagerpb.PoweroffClusterRequest{ClusterName: cluster}); err != nil {
		return fmt.Errorf("request to poweroff cluster: %w", err)
	}

	return nil
}

var clusterPoweroffCmd = &cobra.Command{
	Use:   "poweroff --cluster <my-cluster> --org <my-org>",
	Short: "Power off an IPC",
	Long:  poweroffDesc,
	RunE: func(cmd *cobra.Command, argv []string) error {
		ctx := cmd.Context()
		conn, err := newCloudConn(ctx)
		if err != nil {
			return err
		}
		defer conn.Close()

		return poweroffCluster(ctx, conn, clusterName)
	},
}

func init() {
	ClusterCmd.AddCommand(clusterPoweroffCmd)
	clusterPoweroffCmd.PersistentFlags().StringVar(&clusterName, "cluster", "", "Name of the cluster to use")
}
