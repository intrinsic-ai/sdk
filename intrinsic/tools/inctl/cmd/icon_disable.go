// Copyright 2023 Intrinsic Innovation LLC

package icon

import (
	"context"
	"fmt"

	"intrinsic/icon/go/icon"

	"github.com/spf13/cobra"
)

var flagOperationalOnly bool

func disableIcon(ctx context.Context) error {
	ctx, client, err := makeIconClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	group := icon.AllHardware
	if flagOperationalOnly {
		group = icon.OperationalHardwareOnly
	}

	if err := client.Disable(ctx, group); err != nil {
		return fmt.Errorf("failed to disable ICON: %w", err)
	}

	fmt.Println("ICON disabled successfully.")
	return nil
}

var iconDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable all parts on the ICON server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return disableIcon(cmd.Context())
	},
}

func init() {
	iconDisableCmd.Flags().BoolVar(&flagOperationalOnly, "operational-only", false, "Disable only operational hardware modules")
	iconCmd.AddCommand(iconDisableCmd)
}
