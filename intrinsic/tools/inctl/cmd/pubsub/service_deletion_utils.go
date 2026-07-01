// Copyright 2023 Intrinsic Innovation LLC

package pubsub

import (
	"context"
	"fmt"

	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"

	"github.com/spf13/cobra"
)

const (
	keyUninstallServiceAsset = "uninstall-service-asset"
)

// ServiceDeleteCmdRunner handles execution of commands that delete service assets
// used for line orchestration.
type ServiceDeleteCmdRunner struct {
	CmdRunnerBase

	shouldUninstallServiceAsset bool
}

// run implements the core logic of the command:
//   - Deletes all instances of the service.
//   - Optionally uninstalls the service asset.
func (r *ServiceDeleteCmdRunner) run(ctx context.Context) error {
	fmt.Fprintf(
		r.outputWriter,
		"Deleting existing instances of the %v service in the current solution.\n",
		r.serviceName)
	numInstances, err := r.deleteExistingServiceInstances(ctx)
	if err != nil {
		return fmt.Errorf(
			"failed to delete existing instances of the %v service: %w",
			r.serviceName, err)
	}
	fmt.Fprintf(r.outputWriter, "%v instances have been deleted.\n", numInstances)

	if !r.shouldUninstallServiceAsset {
		fmt.Fprintf(
			r.outputWriter,
			"The %v option is disabled, won't try to uninstall the %v service asset.\n",
			keyUninstallServiceAsset,
			r.serviceName)
		return nil
	}

	shouldUninstall := false
	if numInstances > 0 {
		// If there were instances of the service, then the service asset
		// must be installed. No need to check.
		shouldUninstall = true
	} else {
		// There were no instances of the service, but it is still possible that
		// the service asset is installed. Checking it here.
		fmt.Fprintf(
			r.outputWriter,
			"Checking if the %v service asset is installed.\n",
			r.serviceName)
		currentVersion, err := r.getInstalledServiceAssetVersion(ctx)
		if err != nil {
			return fmt.Errorf(
				"failed to determine whether the %v service is installed: %w",
				r.serviceName,
				err)
		}
		if len(currentVersion) != 0 {
			fmt.Fprintf(
				r.outputWriter,
				"%v service, version %v, is installed. Will uninstall.\n", r.serviceName, currentVersion)
			shouldUninstall = true
		} else {
			fmt.Fprintf(
				r.outputWriter,
				"%v service asset is not installed, nothing else to do.\n",
				r.serviceName)
		}
	}

	if shouldUninstall {
		fmt.Fprintf(r.outputWriter, "Uninstalling %v service asset.\n", r.serviceName)
		if err = r.uninstallServiceAsset(ctx); err != nil {
			return fmt.Errorf("failed to uninstall %v service asset: %w", r.serviceName, err)
		}
		fmt.Fprintf(
			r.outputWriter,
			"The %v service asset has been uninstalled.\n",
			r.serviceName)
	}

	return nil
}

// serviceDeleteCmdEnvironment is the execution environment for commands
// that delete services used for line orchestration.
// That environment contains command line
// flags and a connection to the gRPC service.
type serviceDeleteCmdEnvironment struct {
	cmdFlags *cmdutils.CmdFlags

	packageName string
	serviceName string
}

// RunE sets up the execution environment and invokes ServiceDeleteCmdRunner.run.
func (e *serviceDeleteCmdEnvironment) RunE(cmd *cobra.Command, _ []string) error {
	ctx := cmd.Context()

	ctx, conn, _, err := clientutils.DialClusterFromInctl(ctx, e.cmdFlags)
	if err != nil {
		return err
	}
	defer conn.Close()

	runner := &ServiceDeleteCmdRunner{
		CmdRunnerBase: *newCmdRunnerBase(
			conn,
			cmd.OutOrStdout(),
			e.cmdFlags.GetString(cmdutils.KeyCluster),
			e.packageName,
			e.serviceName),
		shouldUninstallServiceAsset: e.cmdFlags.GetBool(keyUninstallServiceAsset),
	}

	return runner.run(ctx)
}

// NewServiceDeleteCmd returns the initialized cobra command for service deletion.
func NewServiceDeleteCmd(use, short, packageName, serviceName string) *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	env := serviceDeleteCmdEnvironment{
		cmdFlags:    flags,
		packageName: packageName,
		serviceName: serviceName,
	}

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.NoArgs,
		RunE:  env.RunE,
	}

	flags.SetCommand(cmd)
	flags.AddFlagsProjectOrg()
	flags.AddFlagsAddressClusterSolution()

	// Can be useful during development, when a service built from source is sideloaded
	// into a solution. In this case, it may be better to keep it installed instead of
	// rebuilding it from source every time (it takes about 20 minutes).
	flags.OptionalBool(keyUninstallServiceAsset, true, "Whether to uninstall the service asset")

	return cmd
}
