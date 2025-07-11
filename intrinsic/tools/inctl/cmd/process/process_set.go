// Copyright 2023 Intrinsic Innovation LLC

package process

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"intrinsic/executive/go/behaviortree"
	"intrinsic/tools/inctl/util/orgutil"
	"intrinsic/util/proto/registryutil"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
	btpb "intrinsic/executive/proto/behavior_tree_go_proto"
	execgrpcpb "intrinsic/executive/proto/executive_service_go_grpc_proto"
	sgrpcpb "intrinsic/frontend/solution_service/proto/solution_service_go_grpc_proto"
	spb "intrinsic/frontend/solution_service/proto/solution_service_go_grpc_proto"
	skillregistrygrpcpb "intrinsic/skills/proto/skill_registry_go_grpc_proto"
)

var allowedSetFormats = []string{TextProtoFormat, BinaryProtoFormat}

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

type deserializer interface {
	deserialize(ctx context.Context, content []byte) (*btpb.BehaviorTree, error)
}

type textDeserializer struct {
	srC skillregistrygrpcpb.SkillRegistryClient
}

func (t *textDeserializer) deserialize(ctx context.Context, content []byte) (*btpb.BehaviorTree, error) {
	// Collect all file descriptor sets from skills.
	skills, err := getSkills(ctx, t.srC)
	if err != nil {
		return nil, errors.Wrapf(err, "could not list skills")
	}

	files := new(protoregistry.Files)
	for _, skill := range skills {
		if err := addFileDescriptorSetToFiles(skill.GetParameterDescription().GetParameterDescriptorFileset(), files); err != nil {
			return nil, errors.Wrap(err, "failed adding file descriptor set to files")
		}
	}

	// To unmarshal expanded Any protos in the given behavior tree correctly, we need all the file
	// descriptor sets from the behavior tree. But to get the file descriptor sets from the behavior
	// tree, we need to unmarshal it first. We solve this by unmarshalling the behavior tree in two
	// passes.
	//
	// Pass 1: Unmarshal with a dummy resolver. All expanded Any protos are unmarshalled to empty
	// messages (more precisely, to Any protos with a correct 'type_url' and empty 'data') but the
	// file descriptor sets in the behavior tree are unmarshalled correctly.
	dummyUnmarshaller := prototext.UnmarshalOptions{
		Resolver:       newResolverToEmpty(),
		AllowPartial:   true,
		DiscardUnknown: true, // To unmarshal any text format to an Empty proto without errors
	}

	btWithEmptyAnys := &btpb.BehaviorTree{}
	if err := dummyUnmarshaller.Unmarshal(content, btWithEmptyAnys); err != nil {
		return nil, errors.Wrapf(err, "could not parse input file in first pass")
	}

	// Collect all file descriptor sets from the behavior tree.
	collector := fileDescriptorSetCollector{}
	behaviortree.Walk(btWithEmptyAnys, &collector)

	for _, fileDescriptorSet := range collector.fileDescriptorSets {
		if err := addFileDescriptorSetToFiles(fileDescriptorSet, files); err != nil {
			return nil, errors.Wrap(err, "failed adding file descriptor set to files")
		}
	}

	types := new(protoregistry.Types)
	if err := registryutil.PopulateTypesFromFiles(types, files); err != nil {
		return nil, errors.Wrapf(err, "failed to populate types from files")
	}

	// Pass 2: Unmarshal with a proper resolver that now uses the file descriptors sets from all
	// skills and from the behavior tree.
	unmarshaller := prototext.UnmarshalOptions{
		Resolver:       types,
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	bt := &btpb.BehaviorTree{}
	if err := unmarshaller.Unmarshal(content, bt); err != nil {
		return nil, errors.Wrapf(err, "could not parse input file in second pass")
	}

	return bt, nil
}

func newTextDeserializer(srC skillregistrygrpcpb.SkillRegistryClient) *textDeserializer {
	return &textDeserializer{srC: srC}
}

type binaryDeserializer struct {
}

func (b *binarySerializer) deserialize(ctx context.Context, content []byte) (*btpb.BehaviorTree, error) {
	bt := &btpb.BehaviorTree{}
	if err := proto.Unmarshal(content, bt); err != nil {
		return nil, errors.Wrapf(err, "could not parse input file")
	}
	return bt, nil
}

func newBinaryDeserializer() *binarySerializer {
	return &binarySerializer{}
}

type setProcessParams struct {
	exC          execgrpcpb.ExecutiveServiceClient
	srC          skillregistrygrpcpb.SkillRegistryClient
	soC          sgrpcpb.SolutionServiceClient
	name         string
	format       string
	content      []byte
	clearTreeID  bool
	clearNodeIDs bool
}

func deserializeBT(ctx context.Context, srC skillregistrygrpcpb.SkillRegistryClient, format string, content []byte) (*btpb.BehaviorTree, error) {
	var d deserializer
	switch format {
	case TextProtoFormat:
		d = newTextDeserializer(srC)
	case BinaryProtoFormat:
		d = newBinaryDeserializer()
	default:
		return nil, fmt.Errorf("unknown format %s", format)
	}

	bt, err := d.deserialize(ctx, content)
	if err != nil {
		return nil, errors.Wrapf(err, "could not serialize BT")
	}
	return bt, nil
}

func setProcess(ctx context.Context, params *setProcessParams) error {
	bt, err := deserializeBT(ctx, params.srC, params.format, params.content)
	if err != nil {
		return errors.Wrapf(err, "could not deserialize BT")
	}

	clearTree(bt, params.clearTreeID, params.clearNodeIDs)

	if params.name == "" {
		if err := setBT(ctx, params.exC, bt); err != nil {
			return errors.Wrapf(err, "could not set active behavior tree")
		}
	} else {
		if _, err := params.soC.CreateBehaviorTree(ctx, &spb.CreateBehaviorTreeRequest{
			BehaviorTreeId: params.name,
			BehaviorTree:   bt,
		}); err != nil {
			return errors.Wrapf(err, "could not create behavior tree in the solution")
		}
	}

	return nil
}

var processSetCmd = &cobra.Command{
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
		projectName := viperLocal.GetString(orgutil.KeyProject)
		orgName := viperLocal.GetString(orgutil.KeyOrganization)
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
			exC:          execgrpcpb.NewExecutiveServiceClient(conn),
			srC:          skillregistrygrpcpb.NewSkillRegistryClient(conn),
			soC:          sgrpcpb.NewSolutionServiceClient(conn),
			content:      content,
			name:         name,
			format:       flagProcessFormat,
			clearTreeID:  flagClearTreeID,
			clearNodeIDs: flagClearNodeIDs,
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
}

func init() {
	processSetCmd.Flags().StringVar(
		&flagProcessFormat, "process_format", TextProtoFormat,
		fmt.Sprintf("(optional) input format. One of: (%s)", strings.Join(allowedSetFormats, ", ")))
	processSetCmd.Flags().StringVar(&flagSolutionName, "solution", "", "Solution to set the process on. For example, use `inctl solutions list --org orgname@projectname --output json [--filter running_in_sim]` to see the list of solutions.")
	processSetCmd.Flags().StringVar(&flagClusterName, "cluster", "", "Cluster to set the process on.")
	processSetCmd.Flags().StringVar(&flagInputFile, "input_file", "", "File from which to read the process.")
	processCmd.AddCommand(processSetCmd)

}
