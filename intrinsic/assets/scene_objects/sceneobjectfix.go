// Copyright 2023 Intrinsic Innovation LLC

// Package sceneobjectfix contains utils that adapt SceneObjects to meet the requirements of the
// latest platform version.
package sceneobjectfix

import (
	sompb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// Manifest updates a SceneObjectManifest to meet the requirements of the latest platform version.
func Manifest(manifest *sompb.SceneObjectManifest) error {
	return nil
}

// ProcessedManifest updates a ProcessedSceneObjectManifest to meet the requirements of the latest
// platform version.
func ProcessedManifest(manifest *sompb.ProcessedSceneObjectManifest) error {
	// Backfill missing FileDescriptorSet when user data is not specified.
	if len(manifest.GetAssets().GetSceneObjectModel().GetUserData()) == 0 && manifest.GetAssets().GetFileDescriptorSet() == nil {
		if manifest.GetAssets() == nil {
			manifest.Assets = &sompb.ProcessedSceneObjectAssets{}
		}
		manifest.Assets.FileDescriptorSet = &descriptorpb.FileDescriptorSet{}
	}
	return nil
}
