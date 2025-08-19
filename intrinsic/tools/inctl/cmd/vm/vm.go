// Copyright 2023 Intrinsic Innovation LLC

// Package vm provides commands to get and manage VMs.
package vm

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"intrinsic/config/environments"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/cmd/vm/pool/pool"
	"intrinsic/tools/inctl/util/cobrautil"
	"intrinsic/tools/inctl/util/orgutil"

	leaseapigrpcpb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_grpc_proto"
	vmpoolsgrpcpb "intrinsic/kubernetes/vmpool/service/api/v1/vmpool_api_go_grpc_proto"
)

var viperLocal = viper.New()

var vmCmd = orgutil.WrapCmd(cobrautil.ParentOfNestedSubcommands("vm", "Administer and work with virtual machines"), viperLocal)

const serviceTag string = "inctl"

var (
	flagProject       string
	flagUserEmail     string
	orgID             string
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

	vmCmd.AddCommand(pool.PoolCmd)

	vmLeaseCmd.PersistentFlags().StringVarP(&flagDuration, "duration", "d", "30m", "Desired duration of the lease. Optional.")
	vmLeaseCmd.PersistentFlags().StringVarP(&flagPool, "pool", "l", "", "The pool to use. Optional.")
	vmLeaseCmd.PersistentFlags().StringVarP(&flagRuntime, "runtime", "", "", "The intrinsic platform runtime version that the VM should run. Optional. If specified, a VM will be created with the specified runtime version. The runtime version should be in the form of 'intrinsic.platform.YYYYMMDD.RCXX'.")
	vmLeaseCmd.PersistentFlags().StringVarP(&flagIntrinsicOS, "intrinsic-os", "", "", "The IntinsicOS version to use for the pool's VMs. If not specified, the IntrinsicOS version at the time of execution will be set. The IntrinsicOS version should be in the form of 'YYYYMMDD.RCXX'.")
	vmLeaseCmd.MarkFlagsMutuallyExclusive("pool", "runtime")
	vmLeaseCmd.MarkFlagsMutuallyExclusive("pool", "intrinsic-os")
	vmLeaseCmd.PersistentFlags().BoolVarP(&flagSilent, "silent", "s", false, "Suppress output and only print the vm identifier for leases.")
	vmLeaseCmd.PersistentFlags().BoolVarP(&flagRetry, "retry", "r", true, "Retry a lease request until it succeeds.")
	vmLeaseCmd.PersistentFlags().DurationVarP(&flagAbortAfter, "timeout", "t", time.Minute*30, "Abort the lease operation (most useful when combined with --retry) after the given duration.")
	vmLeaseCmd.PersistentFlags().StringVar(&flagReservationID, "reservation-id", "", "A UUID to check/create a reservation on lease failure. If empty, a new UUID will be generated if there are no ready VMs.")
	vmCmd.AddCommand(vmLeaseCmd)

	vmCmd.AddCommand(vmReturnCmd)

	vmExpireInCmd.PersistentFlags().BoolVar(&flagExtendOnly, "extend-only", false, "Only update the lease expiration time if it is longer than the current one.")
	vmCmd.AddCommand(vmExpireInCmd)
}

func checkParams(_ *cobra.Command, _ []string) error {
	flagProject = viperLocal.GetString(orgutil.KeyProject)
	orgID = viperLocal.GetString(orgutil.KeyOrganization)
	if orgID == "" {
		return fmt.Errorf("--org is required")
	}
	return nil
}

func newConn(ctx context.Context) (*grpc.ClientConn, error) {
	// warn that those projects most probably have no VM pool
	noPools := []string{"intrinsic-portal", "intrinsic-assets", "intrinsic-accounts"}
	for _, p := range noPools {
		if strings.HasPrefix(flagProject, p) {
			fmt.Fprintf(os.Stderr, "Warning: Project %q has most probably no VM pool. You probably meant to target a compute/backend project like intrinsic-prod-us instead.", flagProject)
		}
	}

	addr := "www.endpoints." + flagProject + ".cloud.goog:443"

	cfg, err := auth.NewStore().GetConfiguration(flagProject)
	if err != nil {
		return nil, err
	}
	creds, err := cfg.GetDefaultCredentials()
	if err != nil {
		return nil, err
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(creds),
		grpc.WithStatsHandler(new(ocgrpc.ClientHandler)),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
	}
	conn, err := grpc.NewClient(addr, grpcOpts...)
	if err != nil {
		return nil, errors.Wrapf(err, "grpc.NewClient(%q)", addr)
	}
	return conn, nil
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
		return nil, err
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

// resolveVM resolves the given VM name or alias to a VM name. This enables
// commands to support both raw VM names and context aliases very easy.
func resolveVM(vmOrAlias, project string) string {
	// If the name looks like a pool VM return it directly.
	if isPoolVM(vmOrAlias) {
		return vmOrAlias
	}
	// Always fall back to the given name if we don't know better.
	return vmOrAlias
}
