// Copyright 2023 Intrinsic Innovation LLC

// Package icon contains all commands for icon handling.
package icon

import (
	"context"

	"intrinsic/assets/clientutils"
	"intrinsic/assets/cmdutils"
	"intrinsic/icon/go/icon"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"

	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

var (
	flagVerbose      bool
	flagInstanceName string
	flags            *cmdutils.CmdFlags
)

// makeIconClient is a variable instead of a regular function to allow tests
// to override it and inject a mock client (e.g., connecting via a Unix domain socket).
var makeIconClient = func(ctx context.Context) (context.Context, icon.Client, error) {
	ctx, conn, address, err := clientutils.DialClusterFromInctl(ctx, flags)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to dial cluster")
	}

	authCtx := clientutils.AuthInsecureConn(ctx, address, flags.GetFlagProject())
	if flagInstanceName != "" {
		authCtx = metadata.AppendToOutgoingContext(authCtx, "x-resource-instance-name", flagInstanceName)
	}
	client := icon.InitClientFromConn(conn)

	return authCtx, client, nil
}

var iconCmd = cobrautil.ParentOfNestedSubcommands("icon", "Introspect and operate ICON")

func init() {
	flags = cmdutils.NewCmdFlags()
	flags.SetCommand(iconCmd)
	flags.AddFlagsAddressClusterSolution()
	flags.AddFlagsProjectOrg()

	iconCmd.PersistentFlags().StringVar(&flagInstanceName, "instance_name", "", "name of the ICON instance to connect to")

	root.RootCmd.AddCommand(iconCmd)
}
