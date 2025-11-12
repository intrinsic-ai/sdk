// Copyright 2023 Intrinsic Innovation LLC

// Package processgen supports generating a Process asset bundle.
package processgen

import (
	"errors"
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
	"intrinsic/assets/bundleio"
	"intrinsic/assets/processes/processutil"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/proto/registryutil"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	processmanifestpb "intrinsic/assets/processes/proto/process_manifest_go_proto"
	behaviortreepb "intrinsic/executive/proto/behavior_tree_go_proto"
)

func readTextProtoWithAnys(path string, message proto.Message, types *protoregistry.Types) error {
	if err := protoio.ReadTextProto(path, message, protoio.WithResolver(types)); err != nil {
		// Ideally, we would include this additional hint only in the case of
		// protoregistry.NotFound. But prototext does not wrap this error
		// correctly and we can only check for the much broader proto.Error.
		hint := ""
		if errors.Is(err, proto.Error) {
			hint = " (if a message type cannot be resolved make sure there are " +
				"no expanded Any protos in the text proto or provide the corresponding textproto_deps " +
				"to the intrinsic_process rule)"
		}

		return fmt.Errorf("failed to read %v from textproto: %w%s", message.ProtoReflect().Descriptor().Name(), err, hint)
	}

	return nil
}

// CreateProcessAssetBundleOptions contains the options for
// CreateProcessAssetBundle().
type CreateProcessAssetBundleOptions struct {
	ManifestPath                    string
	BehaviorTreePath                string
	TextprotoFileDescriptorSetPaths []string
	OutputBundlePath                string
	OutputFileDescriptorSetPath     string
	OutputManifestBinaryPath        string
}

// CreateProcessAssetBundle creates a ProcessAsset tar bundle from a
// ProcessManifest on disk.
func CreateProcessAssetBundle(options CreateProcessAssetBundleOptions) error {
	textprotoFileDescriptorSet, err := registryutil.LoadFileDescriptorSets(options.TextprotoFileDescriptorSetPaths)
	if err != nil {
		return fmt.Errorf("unable to loads FileDescriptorSets: %w", err)
	}

	types, err := registryutil.NewTypesFromFileDescriptorSet(textprotoFileDescriptorSet)
	if err != nil {
		return fmt.Errorf("failed to build types: %w", err)
	}

	// Load ProcessManifest from textprotos.
	manifest := &processmanifestpb.ProcessManifest{}
	if err := readTextProtoWithAnys(options.ManifestPath, manifest, types); err != nil {
		return err
	}

	if options.BehaviorTreePath != "" {
		if manifest.BehaviorTree != nil {
			return fmt.Errorf("behavior tree path given but manifest already contains a behavior tree")
		}

		manifest.BehaviorTree = &behaviortreepb.BehaviorTree{}
		if err := readTextProtoWithAnys(options.BehaviorTreePath, manifest.BehaviorTree, types); err != nil {
			return err
		}
	}

	if err := processutil.ValidateProcessManifest(manifest); err != nil {
		return fmt.Errorf("invalid ProcessManifest: %w", err)
	}

	// Open the file at the output bundle path for writing. Creates the file if it
	// does not already exist.
	outputBundleFile, err := os.OpenFile(options.OutputBundlePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open %q for writing: %w", options.OutputBundlePath, err)
	}
	defer outputBundleFile.Close()

	// Write the ProcessManifest to the output file.
	if err := bundleio.WriteProcessManifest(manifest, outputBundleFile); err != nil {
		return fmt.Errorf("failed to write ProcessManifest: %w", err)
	}

	// Write the FileDescriptorSet of the asset to the output file. Write an empty
	// file if the behavior tree does not have a parameter file descriptor set.
	fileDescriptorSet := manifest.GetBehaviorTree().
		GetDescription().GetParameterDescription().GetParameterDescriptorFileset()
	if fileDescriptorSet == nil {
		fileDescriptorSet = &descriptorpb.FileDescriptorSet{}
	}
	writeOptions := protoio.WithDeterministic(true)
	if err := protoio.WriteBinaryProto(options.OutputFileDescriptorSetPath, fileDescriptorSet, writeOptions); err != nil {
		return fmt.Errorf("failed to write FileDescriptorSet: %w", err)
	}

	if err := protoio.WriteBinaryProto(options.OutputManifestBinaryPath, manifest, writeOptions); err != nil {
		return fmt.Errorf("failed to write ProcessManifest: %w", err)
	}

	return nil
}
