// Copyright 2023 Intrinsic Innovation LLC

// Package sceneobjectvalidate provides utils for validating SceneObjects.
package sceneobjectvalidate

import (
	"fmt"

	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"

	sompb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
	sopb "intrinsic/scene/proto/v1/scene_object_go_proto"
)

type sceneObjectManifestOptions struct {
	gzfPaths map[string]string
	files    *protoregistry.Files
}

// SceneObjectManifestOption is an option for validating a SceneObjectManifest.
type SceneObjectManifestOption func(*sceneObjectManifestOptions)

// WithFiles provides a Files for validating proto messages.
func WithFiles(files *protoregistry.Files) SceneObjectManifestOption {
	return func(opts *sceneObjectManifestOptions) {
		opts.files = files
	}
}

// WithGZFPaths adds a map from GZF file paths as specified in the manifest to paths on disk.
//
// Must be specified if the manifest specifies GZF files.
func WithGZFPaths(gzfPaths map[string]string) SceneObjectManifestOption {
	return func(opts *sceneObjectManifestOptions) {
		opts.gzfPaths = gzfPaths
	}
}

// SceneObjectManifest validates a SceneObjectManifest.
func SceneObjectManifest(m *sompb.SceneObjectManifest, options ...SceneObjectManifestOption) error {
	opts := &sceneObjectManifestOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if m == nil {
		return fmt.Errorf("SceneObjectManifest must not be nil")
	}

	if err := metadatautils.ValidateManifestMetadata(m.GetMetadata()); err != nil {
		return fmt.Errorf("invalid SceneObjectManifest metadata: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetMetadata().GetId())

	if numGZF := len(m.GetAssets().GetGzfGeometryFilenames()); numGZF > 1 {
		return fmt.Errorf("support for multiple gzf files within a SceneObject is not yet implemented (got %d files)", numGZF)
	}
	if name := m.GetAssets().GetRootSceneObjectName(); len(name) != 0 {
		return fmt.Errorf("support for multiple gzf files within a scene_object is not yet implemented, so do not specify a root_scene_object_name (got: %q)", name)
	}
	if numGZF := len(m.GetAssets().GetGzfGeometryFilenames()); numGZF > 1 && len(m.GetAssets().GetRootSceneObjectName()) == 0 {
		return fmt.Errorf("root_scene_object_name must be specified for multiple gzf files")
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
			if name := userData.MessageName(); name == "" {
				return fmt.Errorf("user data %q message must not be an empty Any for %q", key, id)
			} else if opts.files == nil {
				return fmt.Errorf("SceneObject %q has user data (%q, of type %s), but no descriptors provided", id, key, name)
			} else if _, err := opts.files.FindDescriptorByName(name); err != nil {
				return fmt.Errorf("cannot find user data message %q for %q: %w", name, id, err)
			}
		}
	}

	return nil
}

// ProcessedSceneObjectManifest validates a ProcessedSceneObjectManifest.
func ProcessedSceneObjectManifest(m *sompb.ProcessedSceneObjectManifest) error {
	if m == nil {
		return fmt.Errorf("ProcessedSceneObjectManifest must not be nil")
	}

	if err := metadatautils.ValidateManifestMetadata(m.GetMetadata()); err != nil {
		return fmt.Errorf("invalid ProcessedSceneObjectManifest metadata: %w", err)
	}
	id := idutils.IDFromProtoUnchecked(m.GetMetadata().GetId())

	fds := m.GetAssets().GetFileDescriptorSet()
	if fds == nil {
		return fmt.Errorf("FileDescriptorSet must not be nil for %q", id)
	}
	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return fmt.Errorf("failed to populate registry for %q: %v", id, err)
	}

	// Verify that any user data in the SceneObject is in the FileDescriptorSet.
	for key, userData := range m.GetAssets().GetSceneObjectModel().GetUserData() {
		if name := userData.MessageName(); name == "" {
			return fmt.Errorf("user data %q message must not be an empty Any for %q", key, id)
		} else if _, err := files.FindDescriptorByName(name); err != nil {
			return fmt.Errorf("cannot find user data message %q for %q: %w", name, id, err)
		}
	}

	return nil
}
