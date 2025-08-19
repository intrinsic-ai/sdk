// Copyright 2023 Intrinsic Innovation LLC

// Package pool provides commands to administer VM pools.
package pool

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"intrinsic/kubernetes/vmpool/service/pkg/defaults/defaults"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/util/cobrautil"
	"intrinsic/tools/inctl/util/orgutil"

	leaseapigrpcpb "intrinsic/kubernetes/vmpool/manager/api/v1/lease_api_go_grpc_proto"
	vmpoolsgrpcpb "intrinsic/kubernetes/vmpool/service/api/v1/vmpool_api_go_grpc_proto"
)

var viperLocal = viper.New()

// PoolCmd is the parent command for all VM pool commands.
var PoolCmd = orgutil.WrapCmd(cobrautil.ParentOfNestedSubcommands("pool", "Create and manage pools of virtual machines"), viperLocal)

const (
	poolFlagDesc             = "The name of the VM pool."
	runtimeFlagDesc          = "The intrinsic platform runtime version that the VMs should run. If not specified, the runtime version at the time of execution will be set. The runtime version should be in the form of 'intrinsic.platform.YYYYMMDD.RCXX'"
	intOSFlagDesc            = "The IntrinsicOS version to use for the pool's VMs. If not specified, the IntrinsicOS version at the time of execution will be set. The IntrinsicOS version should be in the form of 'YYYYMMDD.RCXX'."
	tierFlagDesc             = "The tier of VM pool sizing. All tiers -> 'inctl vm pool list-tiers'"
	hardwareTemplateFlagDesc = "The hardware template used for each created VM. All hardware templates -> 'inctl vm pool list-hardware-templates'."
)

var (
	flagPool             string
	flagRuntime          string
	flagIntrinsicOS      string
	flagTier             string
	flagHardwareTemplate string
	flagLease            string
	flagProject          string
	orgID                string
	flagVerbose          bool
	flagCount            int
)

func init() {
	vmpoolsCreateCmd.Flags().StringVar(&flagPool, "pool", "", poolFlagDesc)
	vmpoolsCreateCmd.MarkFlagRequired("pool")
	vmpoolsCreateCmd.Flags().StringVar(&flagRuntime, "runtime", "", runtimeFlagDesc)
	vmpoolsCreateCmd.Flags().StringVar(&flagIntrinsicOS, "intrinsic-os", "", intOSFlagDesc)
	vmpoolsCreateCmd.Flags().StringVar(&flagTier, "tier", defaults.Tier, tierFlagDesc)
	vmpoolsCreateCmd.Flags().StringVar(&flagHardwareTemplate, "hwtemplate", defaults.HardwareTemplate, hardwareTemplateFlagDesc)
	PoolCmd.AddCommand(vmpoolsCreateCmd)

	vmpoolsUpdateCmd.Flags().StringVar(&flagPool, "pool", "", poolFlagDesc)
	vmpoolsUpdateCmd.MarkFlagRequired("pool")
	vmpoolsUpdateCmd.Flags().StringVar(&flagRuntime, "runtime", "", runtimeFlagDesc)
	vmpoolsUpdateCmd.Flags().StringVar(&flagIntrinsicOS, "intrinsic-os", "", intOSFlagDesc)
	vmpoolsUpdateCmd.Flags().StringVar(&flagTier, "tier", defaults.Tier, tierFlagDesc)
	vmpoolsUpdateCmd.Flags().StringVar(&flagHardwareTemplate, "hwtemplate", defaults.HardwareTemplate, hardwareTemplateFlagDesc)
	vmpoolsUpdateCmd.MarkFlagsRequiredTogether("runtime", "pool", "hwtemplate", "tier")
	PoolCmd.AddCommand(vmpoolsUpdateCmd)

	PoolCmd.AddCommand(vmpoolsListCmd)
	PoolCmd.AddCommand(vmpoolsListTiersCmd)
	PoolCmd.AddCommand(vmpoolsListHardwareTemplatesCmd)

	vmpoolsDeleteCmd.Flags().StringVar(&flagPool, "pool", "", poolFlagDesc)
	PoolCmd.AddCommand(vmpoolsDeleteCmd)
	vmpoolsResumeCmd.Flags().StringVar(&flagPool, "pool", "", poolFlagDesc)
	PoolCmd.AddCommand(vmpoolsResumeCmd)
	vmpoolsStopCmd.Flags().StringVar(&flagPool, "pool", "", poolFlagDesc)
	PoolCmd.AddCommand(vmpoolsStopCmd)
	vmpoolsPauseCmd.Flags().StringVar(&flagPool, "pool", "", poolFlagDesc)
	PoolCmd.AddCommand(vmpoolsPauseCmd)

	vmpoolsLeasesListCmd.Flags().StringVar(&flagPool, "pool", "", "Filter for leases in a specific pool.")
	vmpoolsLeasesListCmd.MarkFlagRequired("pool")
	vmpoolsLeasesCmd.AddCommand(vmpoolsLeasesListCmd)
	vmpoolsLeasesStopCmd.Flags().StringVar(&flagLease, "lease", "", "The name of the VM lease.")
	vmpoolsLeasesStopCmd.MarkFlagRequired("lease")
	vmpoolsLeasesCmd.AddCommand(vmpoolsLeasesStopCmd)
	PoolCmd.AddCommand(vmpoolsLeasesCmd)
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
