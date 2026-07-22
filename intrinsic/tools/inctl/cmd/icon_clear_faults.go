// Copyright 2023 Intrinsic Innovation LLC

package icon

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func clearFaults(ctx context.Context) error {
	ctx, client, err := makeIconClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.ClearFaults(ctx); err != nil {
		return fmt.Errorf("failed to clear faults: %w", err)
	}

	fmt.Println("Faults cleared successfully.")
	return nil
}

var iconClearFaultsCmd = &cobra.Command{
	Use:   "clear-faults",
	Short: "Clear all faults on the ICON server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return clearFaults(cmd.Context())
	},
}

func init() {
	iconCmd.AddCommand(iconClearFaultsCmd)
}
