// Copyright 2023 Intrinsic Innovation LLC

package pool

import (
	"intrinsic/tools/inctl/util/printer"

	"github.com/spf13/cobra"
	"go.opencensus.io/trace"

	vmpoolspb "intrinsic/kubernetes/vmpool/service/api/v1/vmpool_api_go_proto"
)

var deleteDesc = `
Delete a VM pool.

The VMPool will converge to a state gracefully where it can be deleted.

New leases will be rejected. The existing leases will be allowed to finish until expiation. Ready VMs will be discarded.
Use 'inctl vm pool leases stop' if you want to stop the leases faster. When no VMs exist anymore, the pool will be deleted.

Example:
	inctl vm pool delete --pool vmpool-my-pool --org <my-org>
`

var resumeDesc = `
Resume a VM pool.

Allow new leases to be granted for the pool. Fill the pool with ready VMs.
Will have no sideeffects if the pool is already in state running.

Example:
	inctl vm pool resume --pool vmpool-my-pool --org <my-org>
`

var stopDesc = `
Stop a VM pool.

No new leases will be granted for the pool. The existing leases will be allowed to finish until expiation. Ready VMs will be discarded.
The pool achieves state stopped when there are no more VMs.

Example:
	inctl vm pool stop --pool vmpool-my-pool --org <my-org>
`

var pauseDesc = `
Pause a VM pool.

No new leases will be granted for the pool. Ready VMs will be discarded.
The pool achieves state paused when there are no more ready VMs.

Example:
	inctl vm pool pause --pool vmpool-my-pool --org <my-org>
`

var vmpoolsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a VM pool.",
	Long:  deleteDesc,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ctx, span := trace.StartSpan(ctx, "inctl.vmpools.delete")
		defer span.End()
		prtr := printer.GetDefaultPrinter(cmd)
		cl, err := newVmpoolsClient(ctx)
		if err != nil {
			return err
		}
		_, err = cl.DeletePool(ctx, &vmpoolspb.DeletePoolRequest{Name: flagPool})
		if err != nil {
			return err
		}
		prtr.Println("VM pool will converge to deletion.")
		return nil
	},
}

var vmpoolsResumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a VM pool.",
	Long:  resumeDesc,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ctx, span := trace.StartSpan(ctx, "inctl.vmpools.resume")
		defer span.End()
		prtr := printer.GetDefaultPrinter(cmd)
		cl, err := newVmpoolsClient(ctx)
		if err != nil {
			return err
		}
		_, err = cl.ResumePool(ctx, &vmpoolspb.ResumePoolRequest{Name: flagPool})
		if err != nil {
			return err
		}
		prtr.Println("VM pool resumed.")
		return nil
	},
}

var vmpoolsStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a VM pool.",
	Long:  stopDesc,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ctx, span := trace.StartSpan(ctx, "inctl.vmpools.stop")
		defer span.End()
		prtr := printer.GetDefaultPrinter(cmd)
		cl, err := newVmpoolsClient(ctx)
		if err != nil {
			return err
		}
		_, err = cl.StopPool(ctx, &vmpoolspb.StopPoolRequest{Name: flagPool})
		if err != nil {
			return err
		}
		prtr.Println("VM pool will stop.")
		return nil
	},
}

var vmpoolsPauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause a VM pool.",
	Long:  pauseDesc,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ctx, span := trace.StartSpan(ctx, "inctl.vmpools.pause")
		defer span.End()
		prtr := printer.GetDefaultPrinter(cmd)
		cl, err := newVmpoolsClient(ctx)
		if err != nil {
			return err
		}
		_, err = cl.PausePool(ctx, &vmpoolspb.PausePoolRequest{Name: flagPool})
		if err != nil {
			return err
		}
		prtr.Println("VM pool will pause.")
		return nil
	},
}
