// Copyright 2023 Intrinsic Innovation LLC

package pool

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.opencensus.io/trace"
	"intrinsic/kubernetes/acl/identity"
	"intrinsic/kubernetes/vmpool/service/pkg/defaults/defaults"
	"intrinsic/tools/inctl/util/printer"

	vmpoolspb "intrinsic/kubernetes/vmpool/service/api/v1/vmpool_api_go_grpc_proto"
)

func textForUpsert(command string) string {
	return fmt.Sprintf(`You can specify:
	- Runtime
	- IntrinsicOS
	- tier (preset for pool size, e.g. how many VMs ready anticipated and maximum amount)
	- hardware template (which hardware to use for the VMs inside that pool)

 Example fully specified:
	 inctl vm pool %s \
	 --pool=usecase1 \
	 --tier=%s \
	 --hwtemplate=%s \
	 --runtime=intrinsic.platform.20241108.RC00 \
	 --intrinsic-os=20241108.RC00 \
	 --org=<my-org>

 Example pool using default values (current runtime & os + default tier & hwtemplate):
	 inctl vm pool %s \
	 --pool=usecase2 \
	 --org=<my-org>
	`, command, defaults.Tier, defaults.HardwareTemplate, command)
}

var createDesc = `Create a new VM pool.

After creating the pool, it will need some time to become ready. You can check the status with:
	inctl vm pool list
Pools become visible in the list command once they reach INITIALIZING, which can take 1-3 minutes.

` + textForUpsert("create")

func validateUpsertParams() error {
	if flagTier == "" {
		flagTier = defaults.Tier
	}
	if flagHardwareTemplate == "" {
		flagHardwareTemplate = defaults.HardwareTemplate
	}
	return nil
}

func getCreatePoolRequest() *vmpoolspb.CreatePoolRequest {
	return &vmpoolspb.CreatePoolRequest{
		Name: flagPool,
		Spec: &vmpoolspb.Spec{
			Runtime:          flagRuntime,
			IntrinsicOs:      flagIntrinsicOS,
			PoolTier:         flagTier,
			HardwareTemplate: flagHardwareTemplate,
		},
	}
}

var vmpoolsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new VM pool.",
	Long:  createDesc,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := identity.OrgToContext(cmd.Context(), poolCmdFlags.GetFlagOrganization())
		if err != nil {
			return err
		}
		ctx, span := trace.StartSpan(ctx, "inctl.vmpools.create")
		defer span.End()
		prtr := printer.GetDefaultPrinter(cmd)
		if err := validateUpsertParams(); err != nil {
			return err
		}
		cl, err := newVmpoolsClient(ctx)
		if err != nil {
			return err
		}
		resp, err := cl.CreatePool(ctx, getCreatePoolRequest())
		if err != nil {
			return err
		}
		prtr.Println("VM pool created.")
		prtr.Println(resp)
		return nil
	},
}

func getUpdatePoolRequest() *vmpoolspb.UpdatePoolRequest {
	return &vmpoolspb.UpdatePoolRequest{
		Name: flagPool,
		Spec: &vmpoolspb.Spec{
			Runtime:          flagRuntime,
			IntrinsicOs:      flagIntrinsicOS,
			PoolTier:         flagTier,
			HardwareTemplate: flagHardwareTemplate,
		},
	}
}

var updateDesc = `Update an existing VM pool.

Update will only result in a pool replacement if the new configuration is different from the old one.

To prevent defaults from being applied by accident, all flags have to be specified.

` + textForUpsert("update")

var vmpoolsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing VM pool.",
	Long:  updateDesc,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ctx, span := trace.StartSpan(ctx, "inctl.vmpools.update")
		defer span.End()
		prtr := printer.GetDefaultPrinter(cmd)
		if err := validateUpsertParams(); err != nil {
			return err
		}
		cl, err := newVmpoolsClient(ctx)
		if err != nil {
			return err
		}
		resp, err := cl.UpdatePool(ctx, getUpdatePoolRequest())
		if err != nil {
			return err
		}
		if resp.GetName() == "" {
			prtr.Println("VM pool unchanged.")
			return nil
		}
		prtr.Println("VM pool updated.")
		prtr.Println(resp)
		return nil
	},
}
