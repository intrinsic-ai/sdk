// Copyright 2023 Intrinsic Innovation LLC

// Package process contains all commands for handling processes (behavior trees).
package process

import (
	"context"
	"fmt"
	"os"
	"strings"

	"intrinsic/skills/tools/skill/cmd/dialerutil"
	"intrinsic/skills/tools/skill/cmd/solutionutil"
	"intrinsic/tools/inctl/cmd/root"
	"intrinsic/tools/inctl/util/cobrautil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	btpb "intrinsic/executive/proto/behavior_tree_go_proto"
	execgrpcpb "intrinsic/executive/proto/executive_service_go_grpc_proto"
	exsvcpb "intrinsic/executive/proto/executive_service_go_grpc_proto"
	rmdpb "intrinsic/executive/proto/run_metadata_go_proto"
	skillregistrygrpcpb "intrinsic/skills/proto/skill_registry_go_grpc_proto"
	srpb "intrinsic/skills/proto/skill_registry_go_grpc_proto"
	skillspb "intrinsic/skills/proto/skills_go_proto"

	lrpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	fbpb "google.golang.org/genproto/googleapis/api/annotations"
)

const (
	keyFilter = "filter"
)

const (
	// TextProtoFormat is the textproto output format.
	TextProtoFormat = "textproto"
	// BinaryProtoFormat is the binary proto output format.
	BinaryProtoFormat = "binaryproto"
)

var (
	flagServerAddress string
	flagSolutionName  string
	flagClusterName   string
	flagInputFile     string
	flagOutputFile    string
	flagClearTreeID   bool
	flagClearNodeIDs  bool
	flagProcessFormat string
)

var (
	protoNameBehaviorTree     = proto.MessageName(new(btpb.BehaviorTree))
	protoNameBehaviorTreeNode = proto.MessageName(new(btpb.BehaviorTree_Node))
)

func clearField(fieldName string, refl protoreflect.Message) {
	field := refl.Descriptor().Fields().ByTextName(fieldName)
	if refl.Has(field) {
		refl.Clear(field)
	}
}

func clearTree(m proto.Message, clearTreeID bool, clearNodeIDs bool) error {
	refl := m.ProtoReflect()

	n := proto.MessageName(m)
	if clearTreeID && n == protoNameBehaviorTree {
		clearField("tree_id", refl)
	}
	if clearNodeIDs && n == protoNameBehaviorTreeNode {
		clearField("id", refl)
	}

	for i := 0; i < refl.Descriptor().Fields().Len(); i++ {
		field := refl.Descriptor().Fields().Get(i)
		if !refl.Has(field) {
			continue
		}
		options := field.Options().(*descriptorpb.FieldOptions)
		if proto.HasExtension(options, fbpb.E_FieldBehavior) {
			behaviors := proto.GetExtension(options, fbpb.E_FieldBehavior).([]fbpb.FieldBehavior)
			for _, behavior := range behaviors {
				if behavior == fbpb.FieldBehavior_OUTPUT_ONLY {
					refl.Clear(field)
				}
			}
		}

		if field.Kind() == protoreflect.MessageKind {
			if field.IsList() {
				list := refl.Get(field).List()
				for j := 0; j < list.Len(); j++ {
					if err := clearTree(list.Get(j).Message().Interface(), clearTreeID, clearNodeIDs); err != nil {
						return err
					}
				}
			} else if !field.IsMap() {
				if err := clearTree(refl.Get(field).Message().Interface(), clearTreeID, clearNodeIDs); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func connectToCluster(ctx context.Context, projectName string, orgName string, address string, solutionName string, clusterName string) (context.Context, *grpc.ClientConn, error) {
	if solutionName != "" {
		// Look up solution name via cloud portal.
		ctx, conn, err := dialerutil.DialConnectionCtx(ctx, dialerutil.DialInfoParams{
			CredName: projectName,
			CredOrg:  orgName,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create client connection: %w", err)
		}

		clusterName, err = solutionutil.GetClusterNameFromSolution(ctx, conn, solutionName)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "could not resolve solution to cluster")
		}
	}

	// Establish a gRPC connection to server, cluster, or cloud.
	ctx, conn, err := dialerutil.DialConnectionCtx(ctx, dialerutil.DialInfoParams{
		Address:  address,
		Cluster:  clusterName,
		CredName: projectName,
		CredOrg:  orgName,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create client connection: %w", err)
	}

	return ctx, conn, nil
}

func getActiveBT(ctx context.Context, exC execgrpcpb.ExecutiveServiceClient) (*btpb.BehaviorTree, error) {
	listOpResp, err := exC.ListOperations(ctx, &lrpb.ListOperationsRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "unable to list executive operations")
	}

	if len(listOpResp.Operations) == 0 {
		return nil, fmt.Errorf("no operations found. Did you load a behavior tree into the executive?")
	}

	if len(listOpResp.Operations) > 1 {
		fmt.Fprintf(os.Stderr, "Found %d concurrent operations, getting first one", len(listOpResp.Operations))
	}
	operation := listOpResp.Operations[0]

	metadata := new(rmdpb.RunMetadata)
	if err := operation.GetMetadata().UnmarshalTo(metadata); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal RunMetadata proto")
	}

	return metadata.GetBehaviorTree(), nil
}

func setBT(ctx context.Context, exC execgrpcpb.ExecutiveServiceClient, bt *btpb.BehaviorTree) error {
	listOpResp, err := exC.ListOperations(ctx, &lrpb.ListOperationsRequest{})
	if err != nil {
		return errors.Wrap(err, "unable to list executive operations")
	}

	if len(listOpResp.Operations) > 1 {
		return errors.Errorf("More than one concurrently loaded BT/executive operation, please delete all but one")
	}

	if len(listOpResp.Operations) == 1 {
		operationToDelete := listOpResp.Operations[0]
		if _, err = exC.DeleteOperation(ctx, &lrpb.DeleteOperationRequest{
			Name: operationToDelete.Name,
		}); err != nil {
			return errors.Wrap(err, "unable to delete operation")
		}
	}

	req := &exsvcpb.CreateOperationRequest{}
	req.RunnableType = &exsvcpb.CreateOperationRequest_BehaviorTree{BehaviorTree: bt}

	if _, err = exC.CreateOperation(ctx, req); err != nil {
		return errors.Wrap(err, "unable to create executive operation")
	}

	return nil
}

func getSkills(ctx context.Context, srC skillregistrygrpcpb.SkillRegistryClient) ([]*skillspb.Skill, error) {
	var (
		skills        []*skillspb.Skill
		nextPageToken string
	)
	for {
		resp, err := srC.ListSkills(ctx, &srpb.ListSkillsRequest{
			PageToken: nextPageToken,
		})
		if err != nil {
			return nil, fmt.Errorf("could not list skills: %w", err)
		}
		skills = append(skills, resp.GetSkills()...)
		nextPageToken = resp.GetNextPageToken()
		if nextPageToken == "" {
			break
		}
	}
	return skills, nil
}

// fileDescriptorSetCollector is a behavior tree visitor that collects all file descriptor sets in a
// behavior tree.
type fileDescriptorSetCollector struct {
	fileDescriptorSets []*descriptorpb.FileDescriptorSet
}

func (c *fileDescriptorSetCollector) VisitCondition(ctx context.Context, cond *btpb.BehaviorTree_Condition) error {
	return nil
}

func (c *fileDescriptorSetCollector) VisitNode(ctx context.Context, node *btpb.BehaviorTree_Node) error {
	fileDescriptorSet := node.GetTask().GetExecuteCode().GetFileDescriptorSet()
	if fileDescriptorSet != nil {
		c.fileDescriptorSets = append(c.fileDescriptorSets, fileDescriptorSet)
	}
	return nil
}

func addFileDescriptorSetToFiles(fileDescriptorSet *descriptorpb.FileDescriptorSet, files *protoregistry.Files) error {
	for _, fileDescriptor := range fileDescriptorSet.GetFile() {
		file, err := protodesc.NewFile(fileDescriptor, files)
		if err != nil {
			return errors.Wrap(err, "failed creating file from file descriptor")
		}
		// Add file if not already present. The error returned by RegisterFile() below cannot easily be
		// classified into "not found" vs "other error". So we check for the file's presence first using
		// FindFileByPath() which does return a specific error for "not found".
		fileExists := true
		_, err = files.FindFileByPath(file.Path())
		if errors.Is(err, protoregistry.NotFound) {
			fileExists = false
		} else if err != nil {
			return errors.Wrap(err, "failed finding file by path")
		}

		if !fileExists {
			if err = files.RegisterFile(file); err != nil {
				return errors.Wrap(err, "failed registering file")
			}
		}
	}
	return nil
}

func addCommonGetSetFlags(cmd *cobra.Command) {
	allowedFormats := []string{TextProtoFormat, BinaryProtoFormat}
	cmd.Flags().StringVar(
		&flagProcessFormat, "process_format", TextProtoFormat,
		fmt.Sprintf("(optional) input/output format. One of: (%s)", strings.Join(allowedFormats, ", ")))
	cmd.Flags().StringVar(&flagSolutionName, "solution", "", "Id of the solution to interact with. For example, use `inctl solutions list --org orgname@projectname --output json [--filter running_in_sim]` to see the list of solutions.")
	cmd.Flags().StringVar(&flagClusterName, "cluster", "", "Name of the cluster to interact with.")
	cmd.Flags().StringVar(&flagServerAddress, "server", "", "Server address of the cluster. Format is {ADDRESS}:{PORT}, for example 'localhost:17080'")
	cmd.Flags().BoolVar(&flagClearTreeID, "clear_tree_id", true, "Clear the tree_id field from the BT proto.")
	cmd.Flags().BoolVar(&flagClearNodeIDs, "clear_node_ids", true, "Clear the nodes' id fields from the BT proto.")
}

var processCmd = cobrautil.ParentOfNestedSubcommands(
	root.ProcessCmdName,
	"Interacts with processes (behavior trees)",
)

func init() {
	processCmd.AddCommand(processGetCmd)
	processCmd.AddCommand(processSetCmd)
	root.RootCmd.AddCommand(processCmd)
}
