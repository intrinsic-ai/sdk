// Copyright 2023 Intrinsic Innovation LLC

package cluster

import (
	"context"
	"fmt"

	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/printer"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	workcellmodeservicegrpcpb "intrinsic/kubernetes/workcellmode/proto/workcellmode_service_go_grpc_proto"
	workcellmodeservicepb "intrinsic/kubernetes/workcellmode/proto/workcellmode_service_go_grpc_proto"
)

var (
	address string
	mode    string
)

var clusterModeCmd = &cobra.Command{
	Use:   "mode",
	Short: "Manage operational mode for a cluster",
	Long:  "Set or get mode for the cluster; choose between develop and operate",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("mode can only be used with a get or set subcommand")
	},
}

func getIPCGRPCClient(ctx context.Context) (*grpc.ClientConn, error) {
	if address != "" {
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}
		conn, err := grpc.NewClient(address, opts...)
		if err != nil {
			return nil, errors.Wrapf(err, "grpc.NewClient(%q) failed", address)
		}
		return conn, nil
	}

	conn, err := auth.NewCloudConnection(ctx, auth.WithFlagValues(ClusterCmdViper), auth.WithCluster(clusterName))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

var modeSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set mode for a cluster",
	Long:  "Set mode for the cluster; choose between develop and operate",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		prtr, err := printer.NewPrinter(root.FlagOutput)
		if err != nil {
			return fmt.Errorf("Internal error: %v", err)
		}

		conn, err := getIPCGRPCClient(ctx)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := workcellmodeservicegrpcpb.NewWorkcellModeClient(conn)
		modeEnum := workcellmodeservicepb.MODE_MODE_UNSPECIFIED
		if mode == "develop" {
			modeEnum = workcellmodeservicepb.MODE_MODE_DEVELOP
		} else if mode == "operate" {
			modeEnum = workcellmodeservicepb.MODE_MODE_OPERATE
		} else {
			return fmt.Errorf("Mode must be either develop or operate")
		}
		_, err = client.SetWorkcellMode(ctx, &workcellmodeservicepb.SetWorkcellModeRequest{Mode: modeEnum})
		if err != nil {
			return fmt.Errorf(
				"setting mode to %s failed: %v, please check that the cluster '%v' is running the latest "+
					"version of Intrinsic software", mode, err, clusterName)
		}
		prtr.PrintS(fmt.Sprintf("Workcell is switching to %s", mode))
		return nil
	},
}

var modeGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get mode for a cluster",
	Long:  "Get mode for the cluster; returns develop or operate",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		prtr, err := printer.NewPrinter(root.FlagOutput)
		if err != nil {
			return fmt.Errorf("Internal error: %v", err)
		}

		conn, err := getIPCGRPCClient(ctx)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := workcellmodeservicegrpcpb.NewWorkcellModeClient(conn)

		res, err := client.GetWorkcellMode(ctx, &workcellmodeservicepb.GetWorkcellModeRequest{})
		if err != nil {
			return fmt.Errorf(
				"getting mode failed: %v, please check that the cluster '%v' is running the latest "+
					"version of Intrinsic software", err, clusterName)
		}
		prtr.PrintS(fmt.Sprintf("Workcell is in %v", res.Mode))
		return nil
	},
}

func init() {
	ClusterCmd.AddCommand(clusterModeCmd)
	clusterModeCmd.PersistentFlags().StringVar(&clusterName, "cluster", "", "Name of the cluster to use")
	clusterModeCmd.PersistentFlags().StringVar(&address, "address", "", "Internal flag to directly set the API server address")
	clusterModeCmd.PersistentFlags().MarkHidden("address")

	clusterModeCmd.AddCommand(modeSetCmd)
	modeSetCmd.Flags().StringVar(&mode, "target", "", "Mode to set")

	clusterModeCmd.AddCommand(modeGetCmd)
}
