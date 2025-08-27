// Copyright 2023 Intrinsic Innovation LLC

package vm

import (
	"context"
	"fmt"
	"time"

	log "github.com/golang/glog"
	"github.com/spf13/cobra"
	"go.opencensus.io/trace"

	tpb "google.golang.org/protobuf/types/known/timestamppb"
	leaseapigrpcpb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_grpc_proto"
	leasepb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_grpc_proto"
)

var extendDesc = `
Extend the expiration time of a lease by a duration relative to now.

Use the time units "m" and "h". Specify --extend-only to only update the lease expiration time if it
is longer than the current one.

Example:
	inctl vm expire-in vmp-3f30-x9t7q72u 1h --org <my-org>
` +
	``

var vmExpireInCmd = &cobra.Command{
	Use:   "expire-in",
	Short: "Extend the expiration time of a leased VM by a duration relative to now.",
	Long:  extendDesc,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, span := trace.StartSpan(cmd.Context(), "inctl.vm.expire-in", trace.WithSampler(trace.AlwaysSample()))
		span.AddAttributes(trace.StringAttribute("vm", args[0]))
		span.AddAttributes(trace.StringAttribute("org", orgID))
		defer span.End()
		cl, err := newLeaseClient(ctx)
		if err != nil {
			log.ExitContextf(ctx, "could not create lease client: %v", err)
		}
		return ExpireIn(ctx, cl, args[0], args[1], flagProject, flagExtendOnly)
	},
	PreRunE: checkParams,
}

// ExpireIn extends the expiration time of a lease by a duration relative to now.
func ExpireIn(ctx context.Context, cl leaseapigrpcpb.VMPoolLeaseServiceClient, vmArg, byStr, project string, extendOnly bool) error {
	vm := resolveVM(vmArg, project)
	byDur, err := time.ParseDuration(byStr)
	if err != nil {
		log.ExitContextf(ctx, "%v is not valid for time.ParseDuration: %v", byStr, err)
	}
	to := time.Now().Add(byDur)

	r, err := cl.ExtendTo(ctx, &leasepb.ExtendToRequest{
		Instance: vm, To: tpb.New(to), ExtendOnly: extendOnly, ServiceTag: serviceTag})
	if err != nil {
		return fmt.Errorf("extending lease failed with: %v", err)
	}
	var expires time.Time = r.GetLease().Expires.AsTime()
	fmt.Printf("Lease extended to %s (in %s)\n", expires.Format(time.RFC3339), time.Until(expires).Round(time.Second))
	if expires.Before(to) {
		fmt.Printf("Warning: This is less than what you expected, you wanted %v\n", byStr)
	}
	if extendOnly && expires.After(to) {
		fmt.Printf("Warning: This is longer than requested due to --extend-only, you wanted %v\n", byStr)
	}
	return nil
}
