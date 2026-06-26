// Copyright 2023 Intrinsic Innovation LLC

package pubsub

import (
	"context"
	"fmt"

	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keyUninstallServiceAsset = "uninstall-service-asset"
)

var flagHubServiceDeleteParams = viper.New()

// HubServiceDeleteCmdRunner handles execution of the hub-service-delete command.
// That subcommand deletes the relay service used for line orchestration.
type HubServiceDeleteCmdRunner struct {
	HubServiceCmdRunnerBase

	shouldUninstallServiceAsset bool
}

// run implements the core logic of the command:
//   - Deletes all instances of the relay service.
//   - Optionally uninstalls the service asset.
func (r *HubServiceDeleteCmdRunner) run(ctx context.Context) error {
	fmt.Fprintf(
		r.outputWriter,
		"Deleting existing instances of the relay service in the current solution.\n")
	numInstances, err := r.deleteExistingRelayServiceInstances(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete existing instances of the relay service: %w", err)
	}
	fmt.Fprintf(r.outputWriter, "%v instances have been deleted.\n", numInstances)

	if !r.shouldUninstallServiceAsset {
		fmt.Fprintf(
			r.outputWriter,
			"The %v option is disabled, won't try to uninstall the relay service asset.\n",
			keyUninstallServiceAsset)
		return nil
	}

	shouldUninstall := false
	if numInstances > 0 {
		// If there were instances of the relay service, then the relay service asset
		// must be installed. No need to check.
		shouldUninstall = true
	} else {
		// There were no instances of the relay service, but it is still possible that
		// the relay service asset is installed. Checking it here.
		fmt.Fprintf(
			r.outputWriter,
			"Checking if the relay service asset is installed.\n")
		currentVersion, err := r.getInstalledRelayServiceAssetVersion(ctx)
		if err != nil {
			return fmt.Errorf("failed to determine whether the relay service is installed: %w", err)
		}
		if len(currentVersion) != 0 {
			fmt.Fprintf(
				r.outputWriter,
				"Relay service, version %v, is installed. Will uninstall.\n", currentVersion)
			shouldUninstall = true
		} else {
			fmt.Fprintf(r.outputWriter, "Relay service asset is not installed, nothing else to do.\n")
		}
	}

	if shouldUninstall {
		fmt.Fprintf(r.outputWriter, "Uninstalling relay service asset.\n")
		if err = r.uninstallRelayServiceAsset(ctx); err != nil {
			return fmt.Errorf("failed to uninstall relay service asset: %w", err)
		}
		fmt.Fprintf(r.outputWriter, "The relay service asset has been uninstalled.\n")
	}

	return nil
}

// hubServiceDeleteCmdEnvironment is the execution environment for the
// hub-service-delete command. That environment contains command line
// flags and a connection to the gRPC service.
type hubServiceDeleteCmdEnvironment struct {
	cmdFlags *cmdutils.CmdFlags
}

// RunE sets up the execution environment and invokes HubServiceDeleteCmdRunner.run.
func (e *hubServiceDeleteCmdEnvironment) RunE(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	ctx, conn, _, err := clientutils.DialClusterFromInctl(ctx, e.cmdFlags)
	if err != nil {
		return err
	}
	defer conn.Close()

	runner := &HubServiceDeleteCmdRunner{
		HubServiceCmdRunnerBase:     *newHubServiceCmdRunnerBase(conn, cmd.OutOrStdout(), e.cmdFlags.GetString(cmdutils.KeyCluster)),
		shouldUninstallServiceAsset: e.cmdFlags.GetBool(keyUninstallServiceAsset),
	}

	return runner.run(ctx)
}

// NewHubServiceDeleteCmd returns the initialized cobra command for hub-service-delete.
func NewHubServiceDeleteCmd() *cobra.Command {
	flags := cmdutils.NewCmdFlagsWithViper(flagHubServiceDeleteParams)
	wrapper := hubServiceDeleteCmdEnvironment{cmdFlags: flags}

	cmd := &cobra.Command{
		Use:   "hub-service-delete",
		Short: "Deletes the PubSub Hub service from the currently running solution.",
		Args:  cobra.NoArgs,
		RunE:  wrapper.RunE,
	}

	flags.SetCommand(cmd)
	flags.AddFlagsProjectOrg()
	flags.AddFlagsAddressClusterSolution()

	// Can be useful during development, when the relay service built from source
	// is sideloaded into a solution. In this case, it may be better to keep it installed
	// instead of rebuilding it from source every time (it takes about 20 minutes).
	flags.OptionalBool(keyUninstallServiceAsset, true, "Whether to uninstall the service asset")

	return cmd
}

func init() {
	PubsubCmd.AddCommand(NewHubServiceDeleteCmd())
}
