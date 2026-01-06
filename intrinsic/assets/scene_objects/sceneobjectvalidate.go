// Copyright 2023 Intrinsic Innovation LLC

// Package sceneobjectvalidate provides utils for validating SceneObjects.
package sceneobjectvalidate

import (
	"fmt"

	"intrinsic/assets/idutils"
	"intrinsic/assets/metadatautils"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	sompb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
	sopb "intrinsic/scene/proto/v1/scene_object_go_proto"
)

type sceneobjectManifestOptions struct {
	files    *protoregistry.Files
	gzfPaths map[string]string
}

// SceneObjectManifestOption is an option for validating a SceneObjectManifest.
type SceneObjectManifestOption func(*sceneobjectManifestOptions)

// WithFiles adds a protoregistry.Files for validating proto messages.
func WithFiles(files *protoregistry.Files) SceneObjectManifestOption {
	return func(opts *sceneobjectManifestOptions) {
		opts.files = files
	}
}

// WithGZFPaths adds a map from GZF file paths as specified in the manifest to paths on disk.
//
// Must be specified if the manifest specifies GZF files.
func WithGZFPaths(gzfPaths map[string]string) SceneObjectManifestOption {
	return func(opts *sceneobjectManifestOptions) {
		opts.gzfPaths = gzfPaths
	}
}

// SceneObjectManifest validates a SceneObjectManifest.
func SceneObjectManifest(m *sompb.SceneObjectManifest, options ...SceneObjectManifestOption) error {
	opts := &sceneobjectManifestOptions{}
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
