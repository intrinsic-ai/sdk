// Copyright 2023 Intrinsic Innovation LLC

package icon

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func listParts(ctx context.Context) error {
	ctx, client, err := makeIconClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	parts, err := client.Parts(ctx)
	if err != nil {
		return err
	}

	for _, part := range parts {
		fmt.Println(part)
	}

	return nil
}

var iconListPartsCmd = &cobra.Command{
	Use:   "list-parts",
	Short: "Print a list of available parts",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listParts(cmd.Context())
	},
}

func init() {
	iconCmd.AddCommand(iconListPartsCmd)
}
