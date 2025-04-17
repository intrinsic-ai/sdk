// Copyright 2023 Intrinsic Innovation LLC

// Package sceneobjectmanifest contains tools for working with SceneObjectManifest.
package sceneobjectmanifest

import (
	"fmt"

	"intrinsic/assets/metadatautils"
	sompb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
)

// ValidateSceneObjectManifestOptions contains options for validating a SceneObjectManifest.
type ValidateSceneObjectManifestOptions struct {
}

// ValidateSceneObjectManifestOption is an option for validating a SceneObjectManifest.
type ValidateSceneObjectManifestOption func(*ValidateSceneObjectManifestOptions)

// ValidateSceneObjectManifest validates a SceneObjectManifest.
func ValidateSceneObjectManifest(m *sompb.SceneObjectManifest, options ...ValidateSceneObjectManifestOption) error {
	opts := &ValidateSceneObjectManifestOptions{}
	for _, opt := range options {
		opt(opts)
	}

	if err := metadatautils.ValidateManifestMetadata(m.GetMetadata()); err != nil {
		return fmt.Errorf("invalid manifest: %w", err)
	}

	if len(m.GetAssets().GetGzfGeometryFilenames()) > 1 {
		return fmt.Errorf("support for multiple gzf files within a SceneObject is not yet implemented")
	}
	if len(m.GetAssets().GetRootSceneObjectName()) != 0 {
		return fmt.Errorf("support for multiple gzf files within a scene_object is not yet implemented, as such do not specify a root_scene_object_name")
	}
	if len(m.GetAssets().GetGzfGeometryFilenames()) > 1 && len(m.GetAssets().GetRootSceneObjectName()) == 0 {
		return fmt.Errorf("root_scene_object_name must be specified for multiple gzf files")
	}

	return nil
}
