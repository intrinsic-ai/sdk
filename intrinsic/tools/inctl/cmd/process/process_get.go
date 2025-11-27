// Copyright 2023 Intrinsic Innovation LLC

package process

import (
	"context"
	"fmt"
	"os"

	"intrinsic/executive/go/behaviortree"
	"intrinsic/tools/inctl/util/orgutil"
	"intrinsic/util/proto/registryutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"

	btpb "intrinsic/executive/proto/behavior_tree_go_proto"
	execgrpcpb "intrinsic/executive/proto/executive_service_go_grpc_proto"
	sgrpcpb "intrinsic/frontend/solution_service/proto/solution_service_go_grpc_proto"
	spb "intrinsic/frontend/solution_service/proto/solution_service_go_grpc_proto"
	skillregistrygrpcpb "intrinsic/skills/proto/skill_registry_go_grpc_proto"
)

var viperProcessGet = viper.New()

type serializer interface {
	Serialize(context.Context, *btpb.BehaviorTree) ([]byte, error)
}

type textSerializer struct {
	commonFiles *protoregistry.Files
}

// Serialize serializes the given behavior tree to textproto.
func (t *textSerializer) Serialize(ctx context.Context, bt *btpb.BehaviorTree) ([]byte, error) {
	files := *t.commonFiles

	collector := fileDescriptorSetCollector{}
	if err := behaviortree.Walk(ctx, bt, &collector); err != nil {
		return nil, errors.Wrap(err, "failed walking behavior tree")
	}

	for _, fileDescriptorSet := range collector.fileDescriptorSets {
		if err := addFileDescriptorSetToFiles(fileDescriptorSet, &files); err != nil {
			return nil, errors.Wrap(err, "failed adding file descriptor set to files")
		}
	}

	types := new(protoregistry.Types)
	if err := registryutil.PopulateTypesFromFiles(types, &files); err != nil {
		return nil, errors.Wrapf(err, "failed to populate types from files")
	}

	marshaller := prototext.MarshalOptions{
		Resolver:  types,
		Indent:    "  ",
		Multiline: true,
	}
	s := marshaller.Format(bt)
	return []byte(s), nil
}

func newTextSerializer(ctx context.Context, srC skillregistrygrpcpb.SkillRegistryClient) (*textSerializer, error) {
	skills, err := getSkills(ctx, srC)
	if err != nil {
		return nil, errors.Wrapf(err, "could not list skills")
	}

	files := new(protoregistry.Files)
	for _, skill := range skills {
		if err := addFileDescriptorSetToFiles(skill.GetParameterDescription().GetParameterDescriptorFileset(), files); err != nil {
			return nil, errors.Wrap(err, "failed adding file descriptor set to files")
		}
	}

	return &textSerializer{commonFiles: files}, nil
}

type binarySerializer struct{}

// Serialize serializes the given behavior tree to binary proto.
func (b *binarySerializer) Serialize(ctx context.Context, bt *btpb.BehaviorTree) ([]byte, error) {
	marshaller := proto.MarshalOptions{}
	content, err := marshaller.Marshal(bt)
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal BT")
	}
	return content, nil
}

func newBinarySerializer() *binarySerializer {
	return &binarySerializer{}
}

func serializeBT(ctx context.Context, srC skillregistrygrpcpb.SkillRegistryClient, bt *btpb.BehaviorTree, format string) ([]byte, error) {
	var s serializer
	var err error
	switch format {
	case TextProtoFormat:
		s, err = newTextSerializer(ctx, srC)
		if err != nil {
			return nil, errors.Wrapf(err, "could not create textproto serializer")
		}
	case BinaryProtoFormat:
		s = newBinarySerializer()
	default:
		return nil, fmt.Errorf("unknown format %s", format)
	}

	data, err := s.Serialize(ctx, bt)
	if err != nil {
		return nil, errors.Wrapf(err, "could not serialize BT")
	}

	return data, nil
}

type getProcessParams struct {
	exC          execgrpcpb.ExecutiveServiceClient
	soC          sgrpcpb.SolutionServiceClient
	srC          skillregistrygrpcpb.SkillRegistryClient
	name         string
	format       string
	clearTreeID  bool
	clearNodeIDs bool
}

func getProcess(ctx context.Context, params *getProcessParams) ([]byte, error) {
	var bt *btpb.BehaviorTree
	if params.name == "" {
		activeBT, err := getActiveBT(ctx, params.exC)
		if err != nil {
			return nil, errors.Wrap(err, "could not get active behavior tree")
		}
		bt = activeBT
	} else {
		namedBT, err := params.soC.GetBehaviorTree(ctx, &spb.GetBehaviorTreeRequest{
			Name: params.name,
		})
		if err != nil {
			return nil, errors.Wrap(err, "could not get named behavior tree")
		}
		bt = namedBT
	}

	clearTree(bt, params.clearTreeID, params.clearNodeIDs)

	return serializeBT(ctx, params.srC, bt, params.format)
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
				exC:          execgrpcpb.NewExecutiveServiceClient(conn),
				soC:          sgrpcpb.NewSolutionServiceClient(conn),
				srC:          skillregistrygrpcpb.NewSkillRegistryClient(conn),
				name:         name,
				format:       flagProcessFormat,
				clearTreeID:  flagClearTreeID,
				clearNodeIDs: flagClearNodeIDs,
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
