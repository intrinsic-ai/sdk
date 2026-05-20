// Copyright 2023 Intrinsic Innovation LLC

package pool

import (
	"context"
	"time"

	"intrinsic/tools/inctl/util/printer"

	"github.com/spf13/cobra"

	vmpoolspb "intrinsic/kubernetes/vmpool/service/api/v1/vmpool_api_go_proto"
)

var pollInterval = 10 * time.Second

func waitForPoolReady(ctx context.Context, cmd *cobra.Command, poolName string, desiredSpec *vmpoolspb.Spec) error {
	prtr := printer.GetDefaultPrinter(cmd)
	prtr.Printf("Waiting for VM pool %q to become ready ", poolName)

	for {
		select {
		case <-ctx.Done():
			prtr.Println()
			return ctx.Err()
		default:
		}

		pools, err := fetchPoolsUserfacing(ctx)
		if err != nil {
			prtr.Printf("? (error fetching pools: %v, retrying)", err)
			time.Sleep(pollInterval)
			continue
		}

		var found *vmpoolspb.Pool
		for _, p := range pools {
			if p.GetName() == poolName {
				found = p
				break
			}
		}

		if found == nil {
			// Pool not found yet (takes time to become visible after creation)
			prtr.Printf(".")
			time.Sleep(pollInterval)
			continue
		}

		status := found.GetCurrentStatus()
		spec := found.GetSpec()

		if status == "RUNNING" && specsEqual(spec, desiredSpec) {
			prtr.Printf("\nVM pool %q is ready.\n", poolName)
			return nil
		}

		prtr.Printf(".")
		time.Sleep(pollInterval)
	}
}

func waitForPoolDeletion(ctx context.Context, cmd *cobra.Command, poolName string) error {
	prtr := printer.GetDefaultPrinter(cmd)
	prtr.Printf("Waiting for VM pool %q to be deleted ", poolName)

	for {
		select {
		case <-ctx.Done():
			prtr.Println()
			return ctx.Err()
		default:
		}

		pools, err := fetchPoolsUserfacing(ctx)
		if err != nil {
			prtr.Printf("? (error fetching pools: %v, retrying)", err)
			time.Sleep(pollInterval)
			continue
		}

		var found *vmpoolspb.Pool
		for _, p := range pools {
			if p.GetName() == poolName {
				found = p
				break
			}
		}

		if found == nil {
			prtr.Printf("\nVM pool %q has been deleted.\n", poolName)
			return nil
		}

		if found.GetCurrentStatus() == "DELETED" {
			prtr.Printf("\nVM pool %q is DELETED.\n", poolName)
			return nil
		}

		prtr.Printf(".")
		time.Sleep(pollInterval)
	}
}

func specsEqual(a, b *vmpoolspb.Spec) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.GetRuntime() == b.GetRuntime() &&
		a.GetIntrinsicOs() == b.GetIntrinsicOs() &&
		a.GetHardwareTemplate() == b.GetHardwareTemplate() &&
		a.GetPoolTier() == b.GetPoolTier()
}
