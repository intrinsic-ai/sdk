// Copyright 2023 Intrinsic Innovation LLC

// Package solution contains all commands for solution handling.
package solution

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"intrinsic/tools/inctl/auth/auth"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/orgutil"
)

const (
	keyFilter = "filter"
)

var (
	viperLocal = viper.New()
)

// SolutionCmd is the `inctl solution` command.
var SolutionCmd = orgutil.WrapCmd(&cobra.Command{
	Use:                root.SolutionCmdName,
	Aliases:            []string{root.SolutionsCmdName},
	Short:              "Solution interacts with solutions",
	DisableFlagParsing: true,
}, viperLocal)

func init() {
	root.RootCmd.AddCommand(SolutionCmd)
}

func newCloudConn(ctx context.Context) (*grpc.ClientConn, error) {
	return auth.NewCloudConnection(ctx, auth.WithFlagValues(viperLocal))
}
