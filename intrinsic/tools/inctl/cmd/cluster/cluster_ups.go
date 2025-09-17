// Copyright 2023 Intrinsic Innovation LLC

package cluster

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"intrinsic/tools/inctl/util/orgutil"

	clustermanageralphapb "intrinsic/frontend/cloud/api/v1alpha1/clustermanager_api_go_grpc_proto"
)

var (
	upsFlags struct {
		driver string
		port   string
	}
	experimentalFlag bool
)

var upsCmd = &cobra.Command{
	Use:   "ups",
	Short: "[EXPERIMENTAL] Show the current UPS configuration.",
	Long:  "[EXPERIMENTAL] Show the current UPS configuration of the cluster.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()

		projectName := ClusterCmdViper.GetString(orgutil.KeyProject)
		orgName := ClusterCmdViper.GetString(orgutil.KeyOrganization)

		ctx, c, err := newClient(ctx, orgName, projectName, clusterName)
		if err != nil {
			return fmt.Errorf("device manager client:\n%w", err)
		}
		defer c.close()

		req := clustermanageralphapb.GetConfigRequest{
			Project:   projectName,
			Org:       orgName,
			ClusterId: clusterName,
		}
		resp, err := c.alphaClient.GetConfig(ctx, &req)
		if err != nil {
			return err
		}
		ups := resp.GetUps()
		if ups == nil {
			return fmt.Errorf("this IPC doesn't support UPS monitoring: please update the OS")
		}
		if ups.GetDriver() == "" {
			fmt.Println("No UPS configured.")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintf(w, "driver\tport\n")
		fmt.Fprintf(w, "%s\t%s\n", ups.GetDriver(), ups.GetPort())
		w.Flush()
		return nil
	},
}

var upsEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "[EXPERIMENTAL] Enable a UPS.",
	Long: `[EXPERIMENTAL] Configure a UPS (uninterruptible power supply) attached to the cluster.

It will monitor the UPS and shut the IPC down when running from battery. This only works for
single-node clusters (a single IPC, as opposed to an older Tier 3 system).

To select the driver, see https://networkupstools.org/stable-hcl.html. To select the port, see the
driver-specific documentation.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		projectName := ClusterCmdViper.GetString(orgutil.KeyProject)
		orgName := ClusterCmdViper.GetString(orgutil.KeyOrganization)
		ctx, c, err := newClient(ctx, orgName, projectName, clusterName)
		if err != nil {
			return fmt.Errorf("cluster upgrade client:\n%w", err)
		}
		defer c.close()
		_, err = c.alphaClient.UpdateConfig(ctx, &clustermanageralphapb.UpdateConfigRequest{
			Project:   projectName,
			Org:       orgName,
			ClusterId: clusterName,
			Config: &clustermanageralphapb.Config{
				Ups: &clustermanageralphapb.UPS{
					Driver: upsFlags.driver,
					Port:   upsFlags.port,
				},
			},
		})
		return err
	},
}
var upsDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "[EXPERIMENTAL] Disable a UPS.",
	Long: `[EXPERIMENTAL] Configure the cluster to not connect to a UPS (uninterruptible power supply).

You can still plug it into a UPS, but it won't shut down when running from battery.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		projectName := ClusterCmdViper.GetString(orgutil.KeyProject)
		orgName := ClusterCmdViper.GetString(orgutil.KeyOrganization)
		ctx, c, err := newClient(ctx, orgName, projectName, clusterName)
		if err != nil {
			return fmt.Errorf("cluster upgrade client:\n%w", err)
		}
		defer c.close()
		_, err = c.alphaClient.UpdateConfig(ctx, &clustermanageralphapb.UpdateConfigRequest{
			Project:   projectName,
			Org:       orgName,
			ClusterId: clusterName,
			Config: &clustermanageralphapb.Config{
				Ups: &clustermanageralphapb.UPS{
					// Empty UPS config disables the UPS.
				},
			},
		})
		return err
	},
}

func init() {
	ClusterCmd.AddCommand(upsCmd)
	upsCmd.PersistentFlags().StringVar(&clusterName, "cluster", "", "Name of cluster to configure.")
	upsCmd.MarkPersistentFlagRequired("cluster")
	upsCmd.PersistentFlags().BoolVar(&experimentalFlag, "enable-experimental", false, "Enable experimental features.")
	upsCmd.MarkPersistentFlagRequired("enable-experimental")
	upsCmd.AddCommand(upsDisableCmd)
	upsCmd.AddCommand(upsEnableCmd)
	upsEnableCmd.Flags().StringVar(&upsFlags.driver, "driver", "", "Driver to use for the UPS, eg \"usbhid-ups\" or \"snmp-ups\".")
	upsEnableCmd.MarkFlagRequired("driver")
	upsEnableCmd.Flags().StringVar(&upsFlags.port, "port", "auto", "Driver-dependent identifier for the UPS, eg \"auto\" for USB or \"192.168.1.123:161\" for SNMP")
}
