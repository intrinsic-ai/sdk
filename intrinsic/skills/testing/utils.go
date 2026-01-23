// Copyright 2023 Intrinsic Innovation LLC

// Package utils provides testing utils for Skills.
package utils

import (
	"testing"

	"intrinsic/util/proto/protoio"
	"intrinsic/util/testing/testio"

	documentationpb "intrinsic/assets/proto/documentation_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	sspb "intrinsic/assets/proto/status_spec_go_proto"
	vendorpb "intrinsic/assets/proto/vendor_go_proto"
	imagepb "intrinsic/kubernetes/workcell_spec/proto/image_go_proto"
	psmpb "intrinsic/skills/proto/processed_skill_manifest_go_proto"
	smpb "intrinsic/skills/proto/skill_manifest_go_proto"

	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

type makeSkillManifestOptions struct {
	displayName   string
	documentation *documentationpb.Documentation
	id            *idpb.Id
	manifestPath  string
	vendor        *vendorpb.Vendor
}

// MakeSkillManifestOption is an option for MakeSkillManifest.
type MakeSkillManifestOption func(*makeSkillManifestOptions)

// WithDisplayName specifies the display name to use in the SkillManifest.
func WithDisplayName(displayName string) MakeSkillManifestOption {
	return func(opts *makeSkillManifestOptions) {
		opts.displayName = displayName
	}
}

// WithDocumentation specifies the documentation to use in the SkillManifest.
func WithDocumentation(documentation *documentationpb.Documentation) MakeSkillManifestOption {
	return func(opts *makeSkillManifestOptions) {
		opts.documentation = documentation
	}
}

// WithID specifies the ID to use in the SkillManifest.
func WithID(id *idpb.Id) MakeSkillManifestOption {
	return func(opts *makeSkillManifestOptions) {
		opts.id = id
	}
}

// WithManifestPath specifies the path to the manifest to load.
func WithManifestPath(path string) MakeSkillManifestOption {
	return func(opts *makeSkillManifestOptions) {
		opts.manifestPath = path
	}
}

// WithVendor specifies the vendor to use in the SkillManifest.
func WithVendor(vendor *vendorpb.Vendor) MakeSkillManifestOption {
	return func(opts *makeSkillManifestOptions) {
		opts.vendor = vendor
	}
}

// MakeSkillManifest returns a new SkillManifest.
func MakeSkillManifest(t *testing.T, options ...MakeSkillManifestOption) *smpb.SkillManifest {
	opts := &makeSkillManifestOptions{}
	for _, opt := range options {
		opt(opts)
	}

	var m *smpb.SkillManifest
	if opts.manifestPath != "" {
		m = mustLoadSkillManifest(t, opts.manifestPath)
	} else {
		m = &smpb.SkillManifest{}
	}

	if opts.displayName == "" && m.GetDisplayName() == "" {
		opts.displayName = "Some Skill"
	}
	if opts.displayName != "" {
		m.DisplayName = opts.displayName
	}

	if opts.documentation != nil {
		m.Documentation = opts.documentation
	}

	if opts.id == nil && m.GetId() == nil {
		opts.id = &idpb.Id{
			Name:    "some_skill",
			Package: "ai.intrinsic",
		}
	}
	if opts.id != nil {
		m.Id = opts.id
	}

	if opts.vendor == nil && m.GetVendor() == nil {
		opts.vendor = &vendorpb.Vendor{
			DisplayName: "Intrinsic",
		}
	}
	if opts.vendor != nil {
		m.Vendor = opts.vendor
	}

	return m
}

func mustLoadSkillManifest(t *testing.T, path string) *smpb.SkillManifest {
	t.Helper()

	m := new(smpb.SkillManifest)
	if err := protoio.ReadBinaryProto(testio.MustCreateRunfilePath(t, path), m); err != nil {
		t.Fatalf("failed to read manifest: %v", err)
	}

	return m
}

type makeProcessedSkillManifestOptions struct {
	fileDescriptorSet *dpb.FileDescriptorSet
	metadata          *psmpb.SkillMetadata
	skillDetails      *psmpb.SkillDetails
}

// MakeProcessedSkillManifestOption is an option for MakeProcessedSkillManifest
type MakeProcessedSkillManifestOption func(*makeProcessedSkillManifestOptions)

// WithFileDescriptorSet specifies the file descriptor set to use in the ProcessedSkillManifest.
func WithFileDescriptorSet(fds *dpb.FileDescriptorSet) MakeProcessedSkillManifestOption {
	return func(opts *makeProcessedSkillManifestOptions) {
		opts.fileDescriptorSet = fds
	}
}

// WithProcessedMetadata specifies the metadata to use in the ProcessedSkillManifest.
func WithProcessedMetadata(m *psmpb.SkillMetadata) MakeProcessedSkillManifestOption {
	return func(opts *makeProcessedSkillManifestOptions) {
		opts.metadata = m
	}
}

// WithProcessedSkillDetails specifies the Skill details to use in the ProcessedSkillManifest.
func WithProcessedSkillDetails(sd *psmpb.SkillDetails) MakeProcessedSkillManifestOption {
	return func(opts *makeProcessedSkillManifestOptions) {
		opts.skillDetails = sd
	}
}

// MakeProcessedSkillManifest makes a ProcessedSkillManifest for testing.
func MakeProcessedSkillManifest(t *testing.T, options ...MakeProcessedSkillManifestOption) *psmpb.ProcessedSkillManifest {
	opts := &makeProcessedSkillManifestOptions{
		fileDescriptorSet: &dpb.FileDescriptorSet{},
		metadata: &psmpb.SkillMetadata{
			Id: &idpb.Id{
				Name:    "some_skill",
				Package: "package.some",
			},
			DisplayName: "Some Skill",
			Vendor: &vendorpb.Vendor{
				DisplayName: "Intrinsic",
			},
		},
		skillDetails: &psmpb.SkillDetails{
			StatusInfo: []*sspb.StatusSpec{
				{
					Code:  10001,
					Title: "My status",
				},
			},
		},
	}
	for _, opt := range options {
		opt(opts)
	}

	return &psmpb.ProcessedSkillManifest{
		Metadata: opts.metadata,
		Details:  opts.skillDetails,
		Assets: &psmpb.ProcessedSkillAssets{
			DeploymentType: &psmpb.ProcessedSkillAssets_Image{
				Image: &imagepb.Image{
					Registry: "gcr.io/test-project",
					Name:     "some_skill",
					Tag:      ":skill",
				},
			},
			FileDescriptorSet: opts.fileDescriptorSet,
		},
	}
}
