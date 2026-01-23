// Copyright 2023 Intrinsic Innovation LLC

// Package utils provides testing utils for Processes.
package utils

import (
	"testing"

	"intrinsic/assets/idutils"

	papb "intrinsic/assets/processes/proto/process_asset_go_proto"
	pmpb "intrinsic/assets/processes/proto/process_manifest_go_proto"
	atypepb "intrinsic/assets/proto/asset_type_go_proto"
	documentationpb "intrinsic/assets/proto/documentation_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	mpb "intrinsic/assets/proto/metadata_go_proto"
	vendorpb "intrinsic/assets/proto/vendor_go_proto"
	btpb "intrinsic/executive/proto/behavior_tree_go_proto"
	skpb "intrinsic/skills/proto/skills_go_proto"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/protobuf/proto"
)

type makeProcessManifestOptions struct {
	behaviorTree *btpb.BehaviorTree
	metadata     *pmpb.ProcessMetadata
}

// MakeProcessManifestOption is an option for MakeProcessManifest.
type MakeProcessManifestOption func(*makeProcessManifestOptions)

// WithBehaviorTree specifies the behavior tree to use in the ProcessManifest.
func WithBehaviorTree(behaviorTree *btpb.BehaviorTree) MakeProcessManifestOption {
	return func(opts *makeProcessManifestOptions) {
		opts.behaviorTree = behaviorTree
	}
}

// WithMetadata specifies the metadata to use in the ProcessManifest.
func WithMetadata(metadata *pmpb.ProcessMetadata) MakeProcessManifestOption {
	return func(opts *makeProcessManifestOptions) {
		opts.metadata = metadata
	}
}

// MakeProcessManifest makes a ProcessManifest for testing.
func MakeProcessManifest(t *testing.T, options ...MakeProcessManifestOption) *pmpb.ProcessManifest {
	opts := &makeProcessManifestOptions{
		metadata: &pmpb.ProcessMetadata{
			Id: &idpb.Id{
				Name:    "some_process",
				Package: "package.some",
			},
			DisplayName: "Some Process",
			Vendor: &vendorpb.Vendor{
				DisplayName: "Intrinsic",
			},
		},
	}
	for _, opt := range options {
		opt(opts)
	}
	if opts.behaviorTree == nil {
		opts.behaviorTree = &btpb.BehaviorTree{
			Name: opts.metadata.GetDisplayName(),
			Description: &skpb.Skill{
				SkillName:               opts.metadata.GetId().GetName(),
				PackageName:             opts.metadata.GetId().GetPackage(),
				Id:                      idutils.IDFromProtoUnchecked(opts.metadata.GetId()),
				Description:             opts.metadata.GetDocumentation().GetDescription(),
				BehaviorTreeDescription: &skpb.BehaviorTreeDescription{},
				ParameterDescription: &skpb.ParameterDescription{
					ParameterDescriptorFileset: &dpb.FileDescriptorSet{
						File: []*dpb.FileDescriptorProto{
							{
								Name: proto.String("my_process_params.proto"),
							},
						},
					},
				},
			},
		}
	}

	return &pmpb.ProcessManifest{
		BehaviorTree: opts.behaviorTree,
		Metadata:     opts.metadata,
	}
}

type makeProcessAssetOptions struct {
	behaviorTree *btpb.BehaviorTree
	metadata     *mpb.Metadata
}

// MakeProcessAssetOption is an option for MakeProcessAsset.
type MakeProcessAssetOption func(*makeProcessAssetOptions)

// WithProcessAssetBehaviorTree specifies the behavior tree to use in the ProcessAsset.
func WithProcessAssetBehaviorTree(behaviorTree *btpb.BehaviorTree) MakeProcessAssetOption {
	return func(opts *makeProcessAssetOptions) {
		opts.behaviorTree = behaviorTree
	}
}

// WithProcessAssetMetadata specifies the metadata to use in the ProcessAsset.
func WithProcessAssetMetadata(metadata *mpb.Metadata) MakeProcessAssetOption {
	return func(opts *makeProcessAssetOptions) {
		opts.metadata = metadata
	}
}

// MakeProcessAsset makes a ProcessAsset for testing.
func MakeProcessAsset(t *testing.T, options ...MakeProcessAssetOption) *papb.ProcessAsset {
	opts := &makeProcessAssetOptions{
		metadata: &mpb.Metadata{
			IdVersion: &idpb.IdVersion{
				Id: &idpb.Id{
					Name:    "some_process",
					Package: "package.some",
				},
			},
			DisplayName: "Some Process",
			Vendor: &vendorpb.Vendor{
				DisplayName: "Intrinsic",
			},
			AssetType: atypepb.AssetType_ASSET_TYPE_PROCESS,
			Documentation: &documentationpb.Documentation{
				Description: "My process description",
			},
		},
	}
	for _, opt := range options {
		opt(opts)
	}
	if opts.behaviorTree == nil {
		opts.behaviorTree = &btpb.BehaviorTree{
			Name: opts.metadata.GetDisplayName(),
			Description: &skpb.Skill{
				SkillName:               opts.metadata.GetIdVersion().GetId().GetName(),
				PackageName:             opts.metadata.GetIdVersion().GetId().GetPackage(),
				Id:                      idutils.IDFromProtoUnchecked(opts.metadata.GetIdVersion().GetId()),
				Description:             opts.metadata.GetDocumentation().GetDescription(),
				BehaviorTreeDescription: &skpb.BehaviorTreeDescription{},
				ParameterDescription: &skpb.ParameterDescription{
					ParameterDescriptorFileset: &dpb.FileDescriptorSet{
						File: []*dpb.FileDescriptorProto{
							{
								Name: proto.String("my_process_params.proto"),
							},
						},
					},
				},
			},
		}
	}

	return &papb.ProcessAsset{
		Metadata:     opts.metadata,
		BehaviorTree: opts.behaviorTree,
	}
}
