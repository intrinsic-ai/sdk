// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stopforwarding

import (
	"fmt"

	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/tools/inctl/util/agents"

	"github.com/spf13/cobra"

	"intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common"
	servicedeletionutils "intrinsic/tools/inctl/cmd/pubsub/line_orchestration/common/service_deletion_utils"
	pubsubcmd "intrinsic/tools/inctl/cmd/pubsub/pubsub_cmd"
)

// stopForwardingCmdEnvironment is the execution environment for
// the stop-forwarding command. That command disables PubSub and
// key-value forwarding that was previously enabled with the
// start-forwarding command.
type stopForwardingCmdEnvironment struct {
	cmdFlags *cmdutils.CmdFlags
}

// RunE sets up the execution environment and initiates deletion of the forwarding
// service from the onprem cluster.
func (e *stopForwardingCmdEnvironment) RunE(cmd *cobra.Command, _ []string) error {
	if err := agents.Check(cmd); err != nil {
		return err
	}
	ctx := cmd.Context()

	ctx, conn, _, err := clientutils.DialClusterFromInctl(ctx, e.cmdFlags)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, cluster, _, err := e.cmdFlags.GetFlagsAddressClusterSolution()
	if err != nil {
		return fmt.Errorf("could not get flags: %w", err)
	}

	runner := &servicedeletionutils.ServiceDeleteCmdRunner{
		CmdRunnerBase: *common.NewCmdRunnerBase(
			conn,
			cmd.OutOrStdout(),
			cluster,
			common.ForwardingServicePackage,
			common.ForwardingServiceName),
		ShouldRetainServiceAsset: e.cmdFlags.GetBool(common.KeyRetainServiceAsset),
	}

	return runner.Run(ctx)
}

// NewStopForwardingCmd returns the initialized cobra command for service deletion
// of the forwarding service from an onprem cluster.
func NewStopForwardingCmd(use, short string) *cobra.Command {
	flags := cmdutils.NewCmdFlags()
	env := stopForwardingCmdEnvironment{
		cmdFlags: flags,
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
	flags.OptionalBool(common.KeyRetainServiceAsset, false, "Whether to retain the service asset")

	return cmd
}

func init() {
	pubsubcmd.PubsubCmd.AddCommand(
		NewStopForwardingCmd(
			"stop-forwarding",
			"Stops forwarding of PubSub topics and KV store paths."))
}
