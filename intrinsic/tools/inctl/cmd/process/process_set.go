// Copyright 2023 Intrinsic Innovation LLC

package process

import (
	"context"
	"fmt"
	"io/ioutil"

	protoregistrypb "intrinsic/proto_tools/proto/proto_registry_go_proto"
	"intrinsic/proto_tools/registry/protoregistryclient"
	"intrinsic/tools/inctl/util/orgutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	btpb "intrinsic/executive/proto/behavior_tree_go_proto"
	execgrpcpb "intrinsic/executive/proto/executive_service_go_proto"
	sgrpcpb "intrinsic/frontend/solution_service/proto/solution_service_go_proto"
	spb "intrinsic/frontend/solution_service/proto/solution_service_go_proto"
)

var viperProcessSet = viper.New()

// resolverToEmpty is a dummy implementation of prototext.UnmarshalOptions.Resolver that always
// returns the Empty message type for any type name or type URL.
type resolverToEmpty struct {
	empty *emptypb.Empty
}

func newResolverToEmpty() *resolverToEmpty {
	return &resolverToEmpty{empty: &emptypb.Empty{}}
}

func (d *resolverToEmpty) FindMessageByName(message protoreflect.FullName) (protoreflect.MessageType, error) {
	return d.empty.ProtoReflect().Type(), nil
}

func (d *resolverToEmpty) FindMessageByURL(url string) (protoreflect.MessageType, error) {
	return d.empty.ProtoReflect().Type(), nil
}

func (d *resolverToEmpty) FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error) {
	return nil, errors.New("dummyResolver.FindExtensionByName is not implemented")
}

func (d *resolverToEmpty) FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error) {
	return nil, errors.New("dummyResolver.FindExtensionByNumber is not implemented")
}

func deserializeFromText(ctx context.Context, content []byte, protoRegistry protoregistrypb.ProtoRegistryClient) (*btpb.BehaviorTree, error) {
	// To unmarshal expanded Any protos from script node parameters in the given
	// behavior tree correctly, we need all the file descriptor sets from the
	// behavior tree. But to get the file descriptor sets from the behavior tree,
	// we need to unmarshal it first. We solve this by unmarshalling the behavior
	// tree in two passes.
	//
	// Pass 1: Unmarshal with a dummy resolver. All expanded Any protos are
	// unmarshalled to empty messages (more precisely, to Any protos with a
	// correct 'type_url' and empty 'data') but the file descriptor sets in the
	// behavior tree are unmarshalled correctly.
	emptyUnmarshaller := prototext.UnmarshalOptions{
		Resolver:     newResolverToEmpty(),
		AllowPartial: true,
		// Unmarshal any expanded Any proto to an Empty proto without errors. Since
		// Empty has no fields, all parsed fields are treated as unknown fields and
		// will simply be discarded.
		DiscardUnknown: true,
	}

	btWithEmptyAnys := &btpb.BehaviorTree{}
	if err := emptyUnmarshaller.Unmarshal(content, btWithEmptyAnys); err != nil {
		return nil, errors.Wrapf(err, "could not parse input file in first pass")
	}

	nodeTypes, err := MergedTypesForAllScriptNodesInTree(ctx, btWithEmptyAnys)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating merged Types from behavior tree script nodes")
	}

	// Pass 2: Unmarshal while resolving Intrinsic type URLs using the proto
	// registry and other type URLs using the descriptors collected from the
	// behavior tree (with the compiled-in types as a last fallback).
	unmarshaller := prototext.UnmarshalOptions{
		Resolver: protoregistryclient.NewProtoRegistryResolver(
			ctx,
			protoRegistry,
			[]protoregistryclient.Resolver{nodeTypes, protoregistry.GlobalTypes},
		),
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	bt := &btpb.BehaviorTree{}
	if err := unmarshaller.Unmarshal(content, bt); err != nil {
		return nil, errors.Wrapf(err, "could not parse input file in second pass")
	}

	return bt, nil
}

func deserializeFromBinary(ctx context.Context, content []byte) (*btpb.BehaviorTree, error) {
	bt := &btpb.BehaviorTree{}
	if err := proto.Unmarshal(content, bt); err != nil {
		return nil, errors.Wrapf(err, "could not parse input file")
	}
	return bt, nil
}

type setProcessParams struct {
	executive       execgrpcpb.ExecutiveServiceClient
	solutionService sgrpcpb.SolutionServiceClient
	protoRegistry   protoregistrypb.ProtoRegistryClient
	name            string
	format          string
	content         []byte
	clearTreeID     bool
	clearNodeIDs    bool
}

func deserializeBT(ctx context.Context, format string, content []byte, protoRegistry protoregistrypb.ProtoRegistryClient) (*btpb.BehaviorTree, error) {
	var err error
	var bt *btpb.BehaviorTree

	switch format {
	case TextProtoFormat:
		bt, err = deserializeFromText(ctx, content, protoRegistry)
		if err != nil {
			return nil, errors.Wrapf(err, "could not deserialize BT from text")
		}
	case BinaryProtoFormat:
		bt, err = deserializeFromBinary(ctx, content)
		if err != nil {
			return nil, errors.Wrapf(err, "could not deserialize BT from binary")
		}
	default:
		return nil, fmt.Errorf("unknown format %s", format)
	}

	return bt, nil
}

func setProcess(ctx context.Context, params *setProcessParams) error {
	bt, err := deserializeBT(ctx, params.format, params.content, params.protoRegistry)
	if err != nil {
		return errors.Wrapf(err, "could not deserialize BT")
	}

	clearTree(bt, params.clearTreeID, params.clearNodeIDs)

	if params.name == "" {
		if err := setBT(ctx, params.executive, bt); err != nil {
			return errors.Wrapf(err, "could not set active behavior tree")
		}
	} else {
		if _, err := params.solutionService.CreateBehaviorTree(ctx, &spb.CreateBehaviorTreeRequest{
			BehaviorTreeId: params.name,
			BehaviorTree:   bt,
		}); err != nil {
			return errors.Wrapf(err, "could not create behavior tree in the solution")
		}
	}

	return nil
}

var processSetCmd = orgutil.WrapCmd(
	&cobra.Command{
		Use:   "set",
		Short: "Set process (behavior tree) of a solution. ",
		Long: `Set the process (behavior tree) of a currently deployed solution.

There are two main operation modes. The first one is to set the "active" process
in the executive. This prepares the process for execution. This is the default
behavior if no name is provided as the first argument.

inctl process set --solution my-solution --cluster my-cluster --input_file /tmp/my-process.textproto [--process_format textproto|binaryproto]

---

Alternatively, the process can be added to the solution. The command will do
this if you specify a name for the process as the first argument. This makes the
process available in the list of processes in the Flowstate frontend. The
process will NOT be loaded into the executive. It can instead be executed by
selecting it in the frontend and running from there.

Note: The name you provide as an argument will be set as the "name" field in the
process regardless of the value that may or may not already be present. If there
is already a process with the same name this will fail.

inctl process set name_to_store_with --solution my-solution --cluster my-cluster --input_file /tmp/my-process.textproto [--process_format textproto|binaryproto]`,
		Args: cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if flagInputFile == "" {
				return fmt.Errorf("--input_file must be specified")
			}

			name := ""
			if len(args) == 1 {
				name = args[0]
			}
			projectName := viperProcessSet.GetString(orgutil.KeyProject)
			orgName := viperProcessSet.GetString(orgutil.KeyOrganization)
			ctx, conn, err := connectToCluster(cmd.Context(), projectName,
				orgName, flagServerAddress,
				flagSolutionName, flagClusterName)
			if err != nil {
				return errors.Wrapf(err, "could not dial connection")
			}
			defer conn.Close()

			content, err := ioutil.ReadFile(flagInputFile)
			if err != nil {
				return errors.Wrapf(err, "could not read input file")
			}

			if err = setProcess(ctx, &setProcessParams{
				executive:       execgrpcpb.NewExecutiveServiceClient(conn),
				solutionService: sgrpcpb.NewSolutionServiceClient(conn),
				protoRegistry:   protoregistrypb.NewProtoRegistryClient(conn),
				content:         content,
				name:            name,
				format:          flagProcessFormat,
				clearTreeID:     flagClearTreeID,
				clearNodeIDs:    flagClearNodeIDs,
			}); err != nil {
				return errors.Wrapf(err, "could not set BT")
			}

			if name == "" {
				fmt.Println("BT loaded successfully to the executive. To edit behavior tree in the frontend: Process -> Load -> From executive")
			} else {
				fmt.Println("BT added to the solution. To edit and execute the process in the frontend: Process -> Load -> <process name>")
			}

			return nil
		},
	},
	viperProcessSet,
)

func init() {
	addCommonGetSetFlags(processSetCmd)
	processSetCmd.Flags().StringVar(&flagInputFile, "input_file", "", "File from which to read the process.")
}
