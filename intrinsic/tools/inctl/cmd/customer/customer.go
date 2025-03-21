// Copyright 2023 Intrinsic Innovation LLC

// Package customer provides access to Flowstate features to create and manage organizations.
package customer

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	lropb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"intrinsic/config/environments"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"
	"intrinsic/tools/inctl/util/orgutil"
)

var vipr = viper.New()

// AdminCmd is the `inctl customer` command.
var customerCmd = cobrautil.ParentOfNestedSubcommands("customer", "Manage your Flowstate customers.")

var (
	flagEnvironment   string
	flagDebugRequests bool
)

func init() {
	customerCmd.Hidden = true

	customerCmd.PersistentFlags().StringVar(&flagEnvironment, orgutil.KeyEnvironment, environments.Prod, "The environment to use for the command.")
	customerCmd.PersistentFlags().BoolVar(&flagDebugRequests, "debug-requests", false, "If true, print the full request and response for each API call.")
	customerCmd = orgutil.WrapCmd(customerCmd, vipr)
	root.RootCmd.AddCommand(customerCmd)
}

func addPrefix(s string, prefix string) string {
	if strings.HasPrefix(s, prefix) {
		return s
	}
	return prefix + s
}

func addPrefixes(s []string, prefix string) []string {
	ps := slices.Clone(s)
	for i := range ps {
		ps[i] = addPrefix(ps[i], prefix)
	}
	return ps
}

func protoPrint(p proto.Message) {
	fmt.Println(p.ProtoReflect().Descriptor().Name())
	ms, err := protojson.MarshalOptions{
		Multiline:         true,
		UseProtoNames:     true,
		EmitUnpopulated:   true,
		EmitDefaultValues: true,
	}.Marshal(p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(ms))
}

type getOperationFunc func(ctx context.Context, in *lropb.GetOperationRequest, opts ...grpc.CallOption) (*lropb.Operation, error)

const (
	pollInterval = time.Second * 5
)

func waitForOperation(ctx context.Context, getLongOp getOperationFunc, lro *lropb.Operation, timeout time.Duration) (*lropb.Operation, error) {
	if lro == nil {
		return nil, fmt.Errorf("no operation to wait for")
	}
	if lro.Done {
		fmt.Printf("Operation (%q) completed\n", lro.Name)
		return lro, nil
	}

	fmt.Printf("Waiting for operation (%q) to complete (%.1f seconds timeout, %v poll interval).\n",
		lro.Name, timeout.Seconds(), pollInterval)
	ts := time.Now()
	defer func() {
		fmt.Printf("Waited %.1f seconds for operation.\n", time.Since(ts).Seconds())
	}()

	ctx, stop := context.WithTimeout(ctx, timeout)
	defer stop()
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	req := lropb.GetOperationRequest{Name: lro.Name}
	for {
		select {
		case <-ticker.C:
			lro, err := getLongOp(ctx, &req)
			if err != nil {
				return nil, err
			}
			if !lro.GetDone() {
				continue
			}
			if lro.GetError() != nil {
				return nil, fmt.Errorf("operation %q failed: %v", lro.GetName(), lro.GetError())
			}
			return lro, nil
		case <-ctx.Done():
			return nil, fmt.Errorf("operation %q timed out", lro.GetName())
		}
	}
}
