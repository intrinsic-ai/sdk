// Copyright 2023 Intrinsic Innovation LLC

package pool

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"go.opencensus.io/trace"
	"intrinsic/tools/inctl/util/printer"

	vmpoolspb "intrinsic/kubernetes/vmpool/service/api/v1/vmpool_api_go_grpc_proto"
)

var listDesc = `
List all VM pools in a project.

Example:
	inctl vm pool list --org <my-org>
`

func getPoolRowCommandPrinter(cmd *cobra.Command) printer.CommandPrinter {
	ot := printer.GetFlagOutputType(cmd)
	if ot == printer.OutputTypeText {
		ot = printer.OutputTypeTAB
	}
	cp, err := printer.NewPrinterOfType(
		ot,
		cmd,
		printer.WithDefaultsFromValue(&poolRow{}, func(columns []string) []string {
			return []string{"idx", "name", "current_status", "desired_status", "runtime", "intrinsic_os", "hardware_template", "tier"}
		}),
	)
	if err != nil {
		cmd.PrintErrf("Error setting up output: %v\n", err)
		cp = printer.GetDefaultPrinter(cmd)
	}
	return cp
}

type poolRow struct {
	Index            uint16 `json:"idx"`
	Name             string `json:"name"`
	CurrentStatus    string `json:"current_status"`
	DesiredStatus    string `json:"desired_status"`
	Runtime          string `json:"runtime"`
	IntrinsicOS      string `json:"intrinsic_os"`
	HardwareTemplate string `json:"hardware_template"`
	Tier             string `json:"tier"`
}

func asPoolRow(p *vmpoolspb.Pool, idx int) *poolRow {
	return &poolRow{
		Index:            uint16(idx),
		Name:             p.GetName(),
		CurrentStatus:    p.GetCurrentStatus(),
		DesiredStatus:    p.GetDesiredStatus(),
		Runtime:          p.GetSpec().GetRuntime(),
		IntrinsicOS:      p.GetSpec().GetIntrinsicOs(),
		HardwareTemplate: p.GetSpec().GetHardwareTemplate(),
		Tier:             p.GetSpec().GetPoolTier(),
	}
}

var vmpoolsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all VM pools in a project.",
	Long:  listDesc,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ctx, span := trace.StartSpan(ctx, "inctl.vmpools.list")
		defer span.End()
		return listVMPoolsUserfacing(ctx, cmd)
	},
	PreRunE: checkParams,
}

func listVMPoolsUserfacing(ctx context.Context, cmd *cobra.Command) error {
	prtr := getPoolRowCommandPrinter(cmd)
	cl, err := newVmpoolsClient(ctx)
	if err != nil {
		return err
	}
	resp, err := cl.ListPools(ctx, &vmpoolspb.ListPoolsRequest{})
	if err != nil {
		return err
	}
	var view printer.View = nil // this is to reuse reflectors in default views
	for i, p := range resp.GetPools() {
		view = printer.NextView(asPoolRow(p, i), view)
		prtr.Println(view)
	}
	return printer.Flush(prtr)
}

type printablePool struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
	State  string            `json:"state,omitempty"`
}

func (p *printablePool) String() string {
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "Pool %q [%q]\n", p.Name, p.State)
	for k, v := range p.Labels {
		fmt.Fprintf(sb, "\t%s=%s\n", k, v)
	}
	return sb.String()
}

var listTiersDesc = `
List all VM pool tiers in a project.

Example:
	inctl vm pool list-tiers --org <my-org>
`

var vmpoolsListTiersCmd = &cobra.Command{
	Use:   "list-tiers",
	Short: "List all VM pool tiers in a project.",
	Long:  listTiersDesc,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ctx, span := trace.StartSpan(ctx, "inctl.vmpools.list-tiers")
		defer span.End()
		ot := printer.GetFlagOutputType(cmd)
		if ot == printer.OutputTypeText {
			ot = printer.OutputTypeTAB
		}
		prtr, err := printer.NewPrinterOfType(
			ot,
			cmd,
			printer.WithDefaultsFromValue(&vmpoolspb.Tier{}, nil),
		)
		if err != nil {
			return err
		}
		cl, err := newVmpoolsClient(ctx)
		if err != nil {
			return err
		}
		resp, err := cl.ListTiers(ctx, &vmpoolspb.ListTiersRequest{})
		if err != nil {
			return err
		}
		var view printer.View = nil // this is to reuse reflectors in default views
		for _, t := range resp.GetTiers() {
			view = printer.NextView(t, view)
			prtr.Println(view)
		}
		printer.Flush(prtr)
		return nil
	},
	PreRunE: checkParams,
}

var listHardwareTemplatesDesc = `
List all VM hardware templates in a project.

Example:
	inctl vm pool list-hwtemplates --org <my-org>
`

var vmpoolsListHardwareTemplatesCmd = &cobra.Command{
	Use:   "list-hwtemplates",
	Short: "List all VM hardware templates in a project.",
	Long:  listHardwareTemplatesDesc,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		ctx, span := trace.StartSpan(ctx, "inctl.vmpools.list-hwtemplates")
		defer span.End()
		ot := printer.GetFlagOutputType(cmd)
		if ot == printer.OutputTypeText {
			ot = printer.OutputTypeTAB
		}
		prtr, err := printer.NewPrinterOfType(
			ot,
			cmd,
			printer.WithDefaultsFromValue(&vmpoolspb.HardwareTemplate{}, nil),
		)
		if err != nil {
			return err
		}
		cl, err := newVmpoolsClient(ctx)
		if err != nil {
			return err
		}
		resp, err := cl.ListHardwareTemplates(ctx, &vmpoolspb.ListHardwareTemplatesRequest{})
		if err != nil {
			return err
		}
		var view printer.View = nil // this is to reuse reflectors in default views
		for _, hwt := range resp.GetHwTemplates() {
			view = printer.NextView(hwt, view)
			prtr.Println(view)
		}
		printer.Flush(prtr)
		return nil
	},
	PreRunE: checkParams,
}
