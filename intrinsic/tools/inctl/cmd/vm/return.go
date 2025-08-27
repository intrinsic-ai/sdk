// Copyright 2023 Intrinsic Innovation LLC

package vm

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go.opencensus.io/trace"

	leaseapigrpcpb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_grpc_proto"
	leasepb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_grpc_proto"
)

var returnDesc = `
Return a leased VM back to the pool.

All data on the VM will be lost.

Example:
	inctl vm return vmp-3f30-x9t7q72u --org <my-org>
` +
	``

var vmReturnCmd = &cobra.Command{
	Use:   "return",
	Short: "Return a leased VM back to the pool.",
	Long:  returnDesc,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, span := trace.StartSpan(cmd.Context(), "inctl.vm.return", trace.WithSampler(trace.AlwaysSample()))
		span.AddAttributes(trace.StringAttribute("vm", args[0]))
		span.AddAttributes(trace.StringAttribute("org", orgID))
		defer span.End()
		cl, err := newLeaseClient(ctx)
		if err != nil {
			return err
		}

		return Return(ctx, cl, args[0], flagProject)
	},
	PreRunE: checkParams,
}

// Return returns a leased VM back to the pool.
func Return(ctx context.Context, cl leaseapigrpcpb.VMPoolLeaseServiceClient, vmArg, project string) error {
	vmID := resolveVM(vmArg, project)
	if _, err := cl.Return(ctx, &leasepb.ReturnRequest{Instance: vmID, ServiceTag: serviceTag}); err != nil {
		return fmt.Errorf("return failed with: %v", err)
	}
	fmt.Printf("VM %s returned. Thank you.\n", vmID)
	return nil
}
