// Copyright 2023 Intrinsic Innovation LLC

// Package sceneobjectmanifest contains tools for working with SceneObjectManifest.
package sceneobjectmanifest

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"

	anypb "google.golang.org/protobuf/types/known/anypb"
	sompb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
	sopb "intrinsic/scene/proto/v1/scene_object_go_proto"
)

// ValidateSceneObjectManifestOptions contains options for validating a SceneObjectManifest.
type ValidateSceneObjectManifestOptions struct {
	files           *protoregistry.Files
	defaultUserData *anypb.Any
	gzfPaths        map[string]string
}

// ValidateSceneObjectManifestOption is an option for validating a SceneObjectManifest.
type ValidateSceneObjectManifestOption func(*ValidateSceneObjectManifestOptions)

// WithFiles adds a protoregistry.Files to the validation options.
func WithFiles(files *protoregistry.Files) ValidateSceneObjectManifestOption {
	return func(opts *ValidateSceneObjectManifestOptions) {
		opts.files = files
	}
}

// WithDefaultUserData adds the SceneObject's default user data to the validation options.
//
// Must be specified if the manifest specifies a default user data file.
func WithDefaultUserData(defaultUserData *anypb.Any) ValidateSceneObjectManifestOption {
	return func(opts *ValidateSceneObjectManifestOptions) {
		opts.defaultUserData = defaultUserData
	}
}

// WithGZFPaths adds a map from GZF file paths as specified in the manifest to paths on disk.
//
// Must be specified if the manifest specifies GZF files.
func WithGZFPaths(gzfPaths map[string]string) ValidateSceneObjectManifestOption {
	return func(opts *ValidateSceneObjectManifestOptions) {
		opts.gzfPaths = gzfPaths
	}
}

// ValidateSceneObjectManifest validates a SceneObjectManifest.
func ValidateSceneObjectManifest(m *sompb.SceneObjectManifest, options ...ValidateSceneObjectManifestOption) error {
	opts := &ValidateSceneObjectManifestOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if err := metadatautils.ValidateManifestMetadata(m.GetMetadata()); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetMetadata().GetId())

	if len(m.GetAssets().GetGzfGeometryFilenames()) > 1 {
		return fmt.Errorf("support for multiple gzf files within a SceneObject is not yet implemented")
	}
	if len(m.GetAssets().GetRootSceneObjectName()) != 0 {
		return fmt.Errorf("support for multiple gzf files within a scene_object is not yet implemented, as such do not specify a root_scene_object_name")
	}
	if len(m.GetAssets().GetGzfGeometryFilenames()) > 1 && len(m.GetAssets().GetRootSceneObjectName()) == 0 {
		return fmt.Errorf("root_scene_object_name must be specified for multiple gzf files")
	}

	// Verify that the default user data proto message is in the FileDescriptorSet.
	if m.GetAssets().GetDefaultUserDataFilename() != "" {
		if opts.defaultUserData == nil {
			return fmt.Errorf("default user data file specified (%q), but no default user data provided", m.GetAssets().GetDefaultUserDataFilename())
		}
		if opts.files == nil {
			return fmt.Errorf("default user data proto message specified, but no descriptors provided")
		}
		if _, err := opts.files.FindDescriptorByName(protoreflect.FullName(opts.defaultUserData.MessageName())); err != nil {
			return fmt.Errorf("could not find default user data message %q in provided descriptors for SceneObject %q: %w", opts.defaultUserData.MessageName(), id, err)
		}
	}

	// Verify that any user data in the associated SceneObjects is in the FileDescriptorSet.
	var sceneObjects []*sopb.SceneObject
	for _, gzfManifestPath := range m.GetAssets().GetGzfGeometryFilenames() {
		_, ok := opts.gzfPaths[gzfManifestPath]
		if !ok {
			return fmt.Errorf("gzf file %q specified in manifest, but no on disk path provided", gzfManifestPath)
		}
	}
	for _, sceneObject := range sceneObjects {
		for key, userData := range sceneObject.GetUserData() {
			messageName := protoreflect.FullName(userData.MessageName())
			if opts.files == nil {
				return fmt.Errorf("SceneObject %q has user data (%q, of type %s), but no descriptors provided", id, key, messageName)
			}
			if _, err := opts.files.FindDescriptorByName(messageName); err != nil {
				return fmt.Errorf("could not find user data message %q in provided descriptors for SceneObject %q: %w", messageName, id, err)
			}
		}
	}

	return nil
}
