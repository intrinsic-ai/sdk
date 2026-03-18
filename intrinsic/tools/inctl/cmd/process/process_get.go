// Copyright 2023 Intrinsic Innovation LLC

package process

import (
	"context"
	"fmt"
	"os"

	"intrinsic/proto_tools/registry/protoregistryclient"
	"intrinsic/tools/inctl/util/orgutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"

	btpb "intrinsic/executive/proto/behavior_tree_go_proto"
	execgrpcpb "intrinsic/executive/proto/executive_service_go_proto"
	sgrpcpb "intrinsic/frontend/solution_service/proto/solution_service_go_proto"
	spb "intrinsic/frontend/solution_service/proto/solution_service_go_proto"
	protoregistrygrpcpb "intrinsic/proto_tools/proto/proto_registry_go_proto"
)

var viperProcessGet = viper.New()

func serializeToTextProto(ctx context.Context, bt *btpb.BehaviorTree, protoRegistry protoregistrygrpcpb.ProtoRegistryClient) ([]byte, error) {
	nodeTypes, err := MergedTypesForAllScriptNodesInTree(ctx, bt)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating merged Types from behavior tree script nodes")
	}

	// Marshal while resolving Intrinsic type URLs using the proto registry and
	// other type URLs using the descriptors collected from the behavior tree
	// (with the compiled-in types as a last fallback).
	marshaller := prototext.MarshalOptions{
		Resolver: protoregistryclient.NewProtoRegistryResolver(
			ctx,
			protoRegistry,
			[]protoregistryclient.Resolver{nodeTypes, protoregistry.GlobalTypes},
		),
		Indent:    "  ",
		Multiline: true,
	}
	s := marshaller.Format(bt)
	return []byte(s), nil
}

func serializeToBinaryProto(bt *btpb.BehaviorTree) ([]byte, error) {
	marshaller := proto.MarshalOptions{}
	content, err := marshaller.Marshal(bt)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal BT")
	}
	return content, nil
}

func serializeBT(ctx context.Context, protoRegistry protoregistrygrpcpb.ProtoRegistryClient, bt *btpb.BehaviorTree, format string) ([]byte, error) {
	var data []byte
	var err error
	switch format {
	case TextProtoFormat:
		data, err = serializeToTextProto(ctx, bt, protoRegistry)
		if err != nil {
			return nil, errors.Wrapf(err, "could not serialize BT to text")
		}
	case BinaryProtoFormat:
		data, err = serializeToBinaryProto(bt)
		if err != nil {
			return nil, errors.Wrapf(err, "could not serialize BT to binary")
		}
	default:
		return nil, fmt.Errorf("unknown format %s", format)
	}

	return data, nil
}

type getProcessParams struct {
	executive       execgrpcpb.ExecutiveServiceClient
	solutionService sgrpcpb.SolutionServiceClient
	protoRegistry   protoregistrygrpcpb.ProtoRegistryClient
	name            string
	format          string
	clearTreeID     bool
	clearNodeIDs    bool
}

func getProcess(ctx context.Context, params *getProcessParams) ([]byte, error) {
	var bt *btpb.BehaviorTree
	if params.name == "" {
		activeBT, err := getActiveBT(ctx, params.executive)
		if err != nil {
			return nil, errors.Wrap(err, "could not get active behavior tree")
		}
		bt = activeBT
	} else {
		namedBT, err := params.solutionService.GetBehaviorTree(ctx, &spb.GetBehaviorTreeRequest{
			Name: params.name,
		})
		if err != nil {
			return nil, errors.Wrap(err, "could not get named behavior tree")
		}
		bt = namedBT
	}

	clearTree(bt, params.clearTreeID, params.clearNodeIDs)

	return serializeBT(ctx, params.protoRegistry, bt, params.format)
}

var processGetCmd = orgutil.WrapCmd(
	&cobra.Command{
		Use:   "get",
		Short: "Get process (behavior tree) of a solution. ",
		Long: `Get the process (behavior tree) of a currently deployed solution.

There are two main operation modes. The first one is to get the "active" process
in the executive. This is the default behavior if no name is provided as the
first argument.

inctl process get --solution my-solution-id --cluster my-cluster [--output_file /tmp/process.textproto] [--process_format textproto|binaryproto]

---

Alternatively, the process can be retrieved from the solution. The command will
do this if you specify the name of the process as the first argument. The
process must already exist in the solution.

inctl process get my_process --solution my-solution-id --cluster my-cluster [--output_file /tmp/process.textproto] [--process_format textproto|binaryproto]`,
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) == 1 {
				name = args[0]
			}
			projectName := viperProcessGet.GetString(orgutil.KeyProject)
			orgName := viperProcessGet.GetString(orgutil.KeyOrganization)
			ctx, conn, err := connectToCluster(cmd.Context(), projectName,
				orgName, flagServerAddress,
				flagSolutionName, flagClusterName)
			if err != nil {
				return errors.Wrapf(err, "could not dial connection")
			}
			defer conn.Close()

			content, err := getProcess(ctx, &getProcessParams{
				executive:       execgrpcpb.NewExecutiveServiceClient(conn),
				solutionService: sgrpcpb.NewSolutionServiceClient(conn),
				protoRegistry:   protoregistrygrpcpb.NewProtoRegistryClient(conn),
				name:            name,
				format:          flagProcessFormat,
				clearTreeID:     flagClearTreeID,
				clearNodeIDs:    flagClearNodeIDs,
			})
			if err != nil {
				return errors.Wrapf(err, "could not get BT")
			}

			if flagOutputFile != "" {
				if err := os.WriteFile(flagOutputFile, content, 0o644); err != nil {
					return errors.Wrapf(err, "could not write to file %s", flagOutputFile)
				}
				return nil
			}

			fmt.Println(string(content))

			return nil
		},
	},
	viperProcessGet,
)

func init() {
	addCommonGetSetFlags(processGetCmd)
	processGetCmd.Flags().StringVar(&flagOutputFile, "output_file", "", "If set, writes the process to the given file instead of stdout.")
}
