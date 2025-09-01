// Copyright 2023 Intrinsic Innovation LLC

package cluster

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	clustermanagergrpcpb "intrinsic/frontend/cloud/api/v1/clustermanager_api_go_grpc_proto"
	clustermanagerpb "intrinsic/frontend/cloud/api/v1/clustermanager_api_go_grpc_proto"
)

func deleteCluster(ctx context.Context, conn *grpc.ClientConn, cluster string) error {
	client := clustermanagergrpcpb.NewClustersServiceClient(conn)
	if _, err := client.DeleteCluster(
		ctx, &clustermanagerpb.DeleteClusterRequest{ClusterName: cluster}); err != nil {
		return fmt.Errorf("request to delete cluster: %w", err)
	}

	return nil
}

var clusterDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a cluster in a project",
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
}
