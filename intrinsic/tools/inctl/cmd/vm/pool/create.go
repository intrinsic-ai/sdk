// Copyright 2023 Intrinsic Innovation LLC

package pool

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.opencensus.io/trace"
	fmpb "google.golang.org/protobuf/types/known/fieldmaskpb"
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

func getUpdatePoolRequest(cmd *cobra.Command) *vmpoolspb.UpdatePoolRequest {
	mask := &fmpb.FieldMask{}
	if cmd.Flags().Changed("runtime") {
		mask.Paths = append(mask.Paths, "spec.runtime")
	}
	if cmd.Flags().Changed("intrinsic-os") {
		mask.Paths = append(mask.Paths, "spec.intrinsic_os")
	}
	if cmd.Flags().Changed("tier") {
		mask.Paths = append(mask.Paths, "spec.pool_tier")
	}
	if cmd.Flags().Changed("hwtemplate") {
		mask.Paths = append(mask.Paths, "spec.hardware_template")
	}
	return &vmpoolspb.UpdatePoolRequest{
		Name: flagPool,
		Spec: &vmpoolspb.Spec{
			Runtime:          flagRuntime,
			IntrinsicOs:      flagIntrinsicOS,
			PoolTier:         flagTier,
			HardwareTemplate: flagHardwareTemplate,
		},
		UpdateMask: mask,
	}
}

var updateDesc = `Update an existing VM pool.

- If no flag is specified, no fields will be updated.
- If a flag is specified, that field will be updated with the value of the flag.
- If a flag is specified with an empty string as value, the endpoint will use the default value (e.g., latest runtime version or latest Intrinsic OS version).

Example: Update a pool named 'usecase0' to the latest runtime and Intrinsic OS versions:
	inctl vm pool update --pool usecase0 --runtime "" --intrinsic-os "" --org=<my-org>

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
		resp, err := cl.UpdatePool(ctx, getUpdatePoolRequest(cmd))
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
