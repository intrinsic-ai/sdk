// Copyright 2023 Intrinsic Innovation LLC

// package main is the entrypoint for creating SceneObject Asset bundles.
package main

import (
	"context"
	"flag"

	"intrinsic/assets/scene_objects/build_defs/sceneobjectgen"
	"intrinsic/production/intrinsic"
	intrinsicflag "intrinsic/util/flag"

	log "github.com/golang/glog"
)

var (
	manifestPath           = flag.String("manifest", "", "Path to the SceneObjectManifest textproto file.")
	fileDescriptorSetPaths = intrinsicflag.MultiString("file_descriptor_set", nil, "Path to a binary file descriptor set proto associated with the manifest. Can be repeated.")
	rootSceneObjectName    = flag.String("root_scene_object_name", "", "Name of the root scene object.")
	gzfGeometryFilePaths   = intrinsicflag.MultiString("geometry_gzfile", nil, "Full path to .gzf file holding geometry. Can be repeated.")
	outputBundlePath       = flag.String("output_bundle", "", "Output path for the .tar bundle.")
)

func main() {
	intrinsic.Init()

	ctx := context.Background()
	if err := sceneobjectgen.CreateSceneObjectBundle(ctx, &sceneobjectgen.CreateSceneObjectBundleOptions{
		ManifestPath:           *manifestPath,
		FileDescriptorSetPaths: *fileDescriptorSetPaths,
		RootSceneObjectName:    *rootSceneObjectName,
		GzfGeometryFilePaths:   *gzfGeometryFilePaths,
		OutputBundlePath:       *outputBundlePath,
	}); err != nil {
		log.Exitf("failed to create SceneObject Asset bundle: %v", err)
	}
}
