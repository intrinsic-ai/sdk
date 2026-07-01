// Copyright 2023 Intrinsic Innovation LLC

package pubsub

import (
	"context"
	"fmt"

	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	pb "intrinsic/platform/pubsub/connect/onprem/forwarding_service/forwarding_service_go_proto"

	"github.com/spf13/cobra"
)

const (
	keyForwardedTopics       = "topic"
	keyForwardedKvStorePaths = "kvstore-key"

	keyForwardingServiceVersion = "forwarding-service-version"
)

// StartForwardingCmdRunner handles execution of the start-forwarding command.
// That command installs or updates the forwarding service used for line orchestration.
type StartForwardingCmdRunner struct {
	ServiceInstallingCmdRunner

	topics       []string
	kvStorePaths []string
}

// makeConfig generates configuration of the forwarding service from command line flags.
func (r *StartForwardingCmdRunner) makeConfig() *pb.ForwardingServiceConfig {
	return &pb.ForwardingServiceConfig{
		Topics:      r.topics,
		KvStoreKeys: r.kvStorePaths,
	}
}

func (r *StartForwardingCmdRunner) run(ctx context.Context) error {
	if len(r.topics) == 0 && len(r.kvStorePaths) == 0 {
		return fmt.Errorf("no topics or KV store keys specified to forward")
	}

	return r.updateInstalledServiceInstances(ctx, r.makeConfig())
}

// StartForwardingCmdEnvironment is the execution environment for the
// start-forwarding command. That environment contains command line
// flags and a connection to the gRPC service.
type StartForwardingCmdEnvironment struct {
	cmdFlags *cmdutils.CmdFlags
}

// RunE sets up the execution environment and invokes StartForwardingCmdRunner.run.
func (e *StartForwardingCmdEnvironment) RunE(cmd *cobra.Command, _ []string) error {
	ctx, conn, _, err := clientutils.DialClusterFromInctl(cmd.Context(), e.cmdFlags)
	if err != nil {
		return err
	}
	defer conn.Close()

	runner := &StartForwardingCmdRunner{
		ServiceInstallingCmdRunner: ServiceInstallingCmdRunner{
			CmdRunnerBase: *newCmdRunnerBase(
				conn,
				cmd.OutOrStdout(),
				e.cmdFlags.GetString(cmdutils.KeyCluster),
				forwardingServicePackage,
				forwardingServiceName,
			),
			requestedVersion: e.cmdFlags.GetString(keyForwardingServiceVersion),
		},
		topics:       e.cmdFlags.GetStringSlice(keyForwardedTopics),
		kvStorePaths: e.cmdFlags.GetStringSlice(keyForwardedKvStorePaths),
	}
	return runner.run(ctx)
}

// NewStartForwardingCmd returns the initialized cobra command for start-forwarding.
func NewStartForwardingCmd() *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	commandWrapper := &StartForwardingCmdEnvironment{cmdFlags: flags}

	cmd := &cobra.Command{
		Use:   "start-forwarding",
		Short: "Starts forwarding of PubSub topics and KV store paths.",
		Args:  cobra.NoArgs,
		RunE:  commandWrapper.RunE,
	}

	flags.SetCommand(cmd)

	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagsProjectOrg()
	flags.StringSlice(
		keyForwardedTopics,
		[]string{},
		"List of PubSub topics to forward")
	flags.StringSlice(
		keyForwardedKvStorePaths,
		[]string{},
		"List of KV store paths to forward")
	flags.OptionalString(
		keyForwardingServiceVersion,
		defaultForwardingServiceVersion,
		fmt.Sprintf(
			"Version of the service asset to install. The default value is %v",
			defaultForwardingServiceVersion))

	return cmd
}

func init() {
	PubsubCmd.AddCommand(NewStartForwardingCmd())
}
