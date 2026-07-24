// Copyright 2023 Intrinsic Innovation LLC

// Package sceneobjectgen implements creation of a SceneObject Asset bundle.
package sceneobjectgen

import (
	"context"
	"fmt"

	"intrinsic/assets/scene_objects/sceneobjectbundle"
	"intrinsic/util/proto/protoio"
	"intrinsic/util/proto/registryutil"
	"intrinsic/util/proto/sourcecodeinfoview"

	sompb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
)

// CreateSceneObjectBundleOptions provides the data needed to create a SceneObject Asset bundle.
type CreateSceneObjectBundleOptions struct {
	FileDescriptorSetPaths []string
	GzfGeometryFilePaths   []string
	ManifestPath           string
	OutputBundlePath       string
	RootSceneObjectName    string
}

// CreateSceneObjectBundle creates a SceneObject Asset bundle on disk.
func CreateSceneObjectBundle(ctx context.Context, opts *CreateSceneObjectBundleOptions) error {
	m := &sompb.SceneObjectManifest{}
	if err := protoio.ReadTextProto(opts.ManifestPath, m); err != nil {
		return fmt.Errorf("failed to read manifest: %v", err)
	}

	fds, err := registryutil.LoadFileDescriptorSets(opts.FileDescriptorSetPaths)
	if err != nil {
		return fmt.Errorf("failed to load FileDescriptorSets: %v", err)
	}
	if err := sourcecodeinfoview.PruneSourceCodeInfo(fds); err != nil {
		return fmt.Errorf("failed to prune source code info: %v", err)
	}

	if err := sceneobjectbundle.WriteFile(ctx, m, opts.OutputBundlePath,
		sceneobjectbundle.WithFileDescriptorSet(fds),
		sceneobjectbundle.WithRootSceneObjectName(opts.RootSceneObjectName),
		sceneobjectbundle.WithGZFGeometryFilePaths(opts.GzfGeometryFilePaths),
	); err != nil {
		return fmt.Errorf("failed to write SceneObject Asset bundle: %v", err)
	}

	return nil
}
