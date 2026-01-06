// Copyright 2023 Intrinsic Innovation LLC

// Package skillgen implements creation of a Skill Asset bundle.
package skillgen

import (
	"fmt"

	"intrinsic/skills/skillbundle"
	"intrinsic/util/proto/protoio"

	smpb "intrinsic/skills/proto/skill_manifest_go_proto"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// CreateSkillBundleOptions provides the data needed to create a Skill Asset bundle.
type CreateSkillBundleOptions struct {
	FileDescriptorSetPath string
	ImageTarPath          string
	ManifestPath          string
	OutputBundlePath      string
}

// CreateSkillBundle creates a Skill Asset bundle on disk.
func CreateSkillBundle(opts *CreateSkillBundleOptions) error {
	fds := &descriptorpb.FileDescriptorSet{}
	if err := protoio.ReadBinaryProto(opts.FileDescriptorSetPath, fds); err != nil {
		return fmt.Errorf("failed to read file descriptor set: %w", err)
	}

	m := &smpb.SkillManifest{}
	if err := protoio.ReadBinaryProto(opts.ManifestPath, m); err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	if err := skillbundle.Write(m, opts.OutputBundlePath,
		skillbundle.WithFileDescriptorSet(fds),
		skillbundle.WithImageTarPath(opts.ImageTarPath),
	); err != nil {
		return fmt.Errorf("failed to write Skill Asset bundle: %w", err)
	}

	return nil
}
