// Copyright 2023 Intrinsic Innovation LLC

// Package vm provides commands to get and manage VMs.
package vm

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"intrinsic/assets/cmdutils"
	"intrinsic/config/environments"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/cmd/vm/pool/pool"
	"intrinsic/tools/inctl/util/cobrautil"

	"github.com/spf13/viper"
	"google.golang.org/grpc"

	leaseapigrpcpb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_grpc_proto"
	vmpoolsgrpcpb "intrinsic/kubernetes/vmpool/service/api/v1/vmpool_api_go_grpc_proto"
)

var (
	viperLocal = viper.New()
	vmCmdFlags = cmdutils.NewCmdFlagsWithViper(viperLocal)
	vmCmd      = cobrautil.ParentOfNestedSubcommands("vm", "Administer and work with virtual machines")
)

const serviceTag string = "inctl"

var (
	flagUserEmail     string
	flagContextAlias  string
	flagAbortAfter    time.Duration
	flagReservationID string
	flagPool          string
	flagRuntime       string
	flagIntrinsicOS   string
	flagRetry         bool
	flagDuration      string
	flagSilent        bool
	flagSetContext    bool
	flagExtendOnly    bool
)

func init() {
	root.RootCmd.AddCommand(vmCmd)
	vmCmdFlags.SetCommand(vmCmd)
	vmCmdFlags.AddFlagsProjectOrg()

	vmCmd.AddCommand(pool.PoolCmd)

	vmLeaseCmd.PersistentFlags().StringVarP(&flagDuration, "duration", "d", "30m", "Desired duration of the lease. Optional.")
	vmLeaseCmd.PersistentFlags().StringVarP(&flagPool, "pool", "l", "", "The pool to use. Optional.")
	vmLeaseCmd.PersistentFlags().StringVarP(&flagRuntime, "runtime", "", "", "The intrinsic platform runtime version that the VM should run. Optional. If specified, a VM will be created with the specified runtime version. The runtime version should be in the form of 'intrinsic.platform.YYYYMMDD.RCXX'.")
	vmLeaseCmd.PersistentFlags().StringVarP(&flagIntrinsicOS, "intrinsic-os", "", "", "The IntinsicOS version to use for the pool's VMs. If not specified, the IntrinsicOS version at the time of execution will be set. The IntrinsicOS version should be in the form of 'YYYYMMDD.RCXX'.")
	vmLeaseCmd.MarkFlagsMutuallyExclusive("pool", "runtime")
	vmLeaseCmd.MarkFlagsMutuallyExclusive("pool", "intrinsic-os")
	vmLeaseCmd.PersistentFlags().BoolVarP(&flagSilent, "silent", "s", false, "Suppress output and only print the vm identifier for leases.")
	vmLeaseCmd.PersistentFlags().BoolVarP(&flagRetry, "retry", "r", true, "Retry a lease request until it succeeds.")
	vmLeaseCmd.PersistentFlags().DurationVarP(&flagAbortAfter, "timeout", "t", time.Minute*20, "Abort the lease operation (most useful when combined with --retry) after the given duration.")
	vmLeaseCmd.PersistentFlags().StringVar(&flagReservationID, "reservation-id", "", "A UUID to check/create a reservation on lease failure. If empty, a new UUID will be generated if there are no ready VMs.")
	vmCmd.AddCommand(vmLeaseCmd)

	vmCmd.AddCommand(vmReturnCmd)

	vmExpireInCmd.PersistentFlags().BoolVar(&flagExtendOnly, "extend-only", false, "Only update the lease expiration time if it is longer than the current one.")
	vmCmd.AddCommand(vmExpireInCmd)
}

func newConn(ctx context.Context) (*grpc.ClientConn, error) {
	// warn that those projects most probably have no VM pool
	noPools := []string{"intrinsic-portal", "intrinsic-assets", "intrinsic-accounts"}
	for _, p := range noPools {
		if strings.HasPrefix(vmCmdFlags.GetFlagProject(), p) {
			fmt.Fprintf(os.Stderr, "Warning: Project %q has most probably no VM pool. You probably meant to target a compute/backend project like intrinsic-prod-us instead.", vmCmdFlags.GetFlagProject())
		}
	}
	return auth.NewCloudConnection(ctx, auth.WithFlagValues(viperLocal))
}

func newVmpoolsClient(ctx context.Context) (vmpoolsgrpcpb.VMPoolServiceClient, error) {
	conn, err := newConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create VM pool client: %v", err)
	}
	return vmpoolsgrpcpb.NewVMPoolServiceClient(conn), nil
}

func newLeaseClient(ctx context.Context) (leaseapigrpcpb.VMPoolLeaseServiceClient, error) {
	conn, err := newConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create VM lease client: %v", err)
	}
	return leaseapigrpcpb.NewVMPoolLeaseServiceClient(conn), nil
}

func getPortalURL(project, cluster string) string {
	e := environments.FromComputeProject(project)
	d := environments.PortalDomain(e)
	return fmt.Sprintf("https://%s/solution-editor/%s/%s/", d, project, cluster)
}

func isPoolVM(name string) bool {
	return strings.HasPrefix(name, "vmp-")
}
