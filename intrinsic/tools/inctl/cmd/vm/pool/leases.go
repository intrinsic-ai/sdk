// Copyright 2023 Intrinsic Innovation LLC

package pool

import (
	"github.com/spf13/cobra"
	"go.opencensus.io/trace"
	"intrinsic/tools/inctl/util/cobrautil"
	"intrinsic/tools/inctl/util/printer"

	vmpoolspb "intrinsic/kubernetes/vmpool/service/api/v1/vmpool_api_go_grpc_proto"
)

var vmpoolsLeasesCmd = cobrautil.ParentOfNestedSubcommands("leases", "Administer leases on your VM pools")

var leasesListDesc = `
List all VM leases for a VM pool.

Example:
	inctl vm pool leases list --pool my-pool --org <my-org>
`

type leasesRow struct {
	Idx     uint16 `json:"idx"`
	Name    string `json:"name"`
	Pool    string `json:"pool"`
	Expires string `json:"expires"`
}

func getLeasesRowCommandPrinter(cmd *cobra.Command) printer.CommandPrinter {
	ot := printer.GetFlagOutputType(cmd)
	if ot == printer.OutputTypeText {
		ot = printer.OutputTypeTAB
	}
	cp, err := printer.NewPrinterOfType(
		ot,
		cmd,
		printer.WithDefaultsFromValue(&leasesRow{}, func(columns []string) []string {
			return []string{"idx", "name", "pool", "expires"}
		}),
	)
	if err != nil {
		cmd.PrintErrf("Error setting up output: %v\n", err)
		cp = printer.GetDefaultPrinter(cmd)
	}
	return cp
}

func asLeasesRow(l *vmpoolspb.Lease, idx int) *leasesRow {
	return &leasesRow{
		Idx:     uint16(idx),
		Name:    l.GetName(),
		Pool:    l.GetPoolName(),
		Expires: l.GetExpirationTime().AsTime().String(),
	}
}

var vmpoolsLeasesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VM leases for a VM pool.",
	Long:  leasesListDesc,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ctx, span := trace.StartSpan(ctx, "inctl.vmpools.leases.list")
		defer span.End()
		prtr := getLeasesRowCommandPrinter(cmd)
		cl, err := newVmpoolsClient(ctx)
		if err != nil {
			return err
		}
		resp, err := cl.ListLeases(ctx, &vmpoolspb.ListLeasesRequest{PoolName: flagPool})
		if err != nil {
			return err
		}
		var view printer.View = nil // this is to reuse reflectors in default views
		for i, l := range resp.GetLeases() {
			view = printer.NextView(asLeasesRow(l, i), view)
			prtr.Println(view)
		}
		printer.Flush(prtr)
		return nil
	},
	PreRunE: checkParams,
}

var leasesStopDesc = `
Stop a VM lease.

Example:
	# find the lease that you want to stop
	inctl vm pool leases list --pool my-pool --org <my-org>
	# stop the lease
	inctl vm pool leases stop --lease vmp-my-lease --org <my-org>
`

var vmpoolsLeasesStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a VM lease.",
	Long:  leasesStopDesc,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ctx, span := trace.StartSpan(ctx, "inctl.vmpools.leases.stop")
		defer span.End()
		prtr := printer.GetDefaultPrinter(cmd)
		cl, err := newVmpoolsClient(ctx)
		if err != nil {
			return err
		}
		_, err = cl.StopLease(ctx, &vmpoolspb.StopLeaseRequest{LeaseName: flagLease})
		if err != nil {
			return err
		}
		prtr.Printf("VM lease %s will stop.\n", flagLease)
		return nil
	},
	PreRunE: checkParams,
}
