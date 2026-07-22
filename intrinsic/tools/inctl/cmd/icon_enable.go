// Copyright 2023 Intrinsic Innovation LLC

package icon

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func enableIcon(ctx context.Context) error {
	ctx, client, err := makeIconClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.Enable(ctx); err != nil {
		return fmt.Errorf("failed to enable ICON: %w", err)
	}

	fmt.Println("ICON enabled successfully.")
	return nil
}

var iconEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable all parts on the ICON server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return enableIcon(cmd.Context())
	},
}

func init() {
	iconCmd.AddCommand(iconEnableCmd)
}
