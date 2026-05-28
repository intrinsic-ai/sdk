// Copyright 2023 Intrinsic Innovation LLC

// Package utils provides testing utils for SceneObjects.
package utils

import (
	"testing"

	documentationpb "intrinsic/assets/proto/documentation_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	vendorpb "intrinsic/assets/proto/vendor_go_proto"
	sompb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
	epb "intrinsic/scene/proto/v1/entity_go_proto"
	socpb "intrinsic/scene/proto/v1/scene_object_config_go_proto"
	sopb "intrinsic/scene/proto/v1/scene_object_go_proto"

	"google.golang.org/protobuf/proto"
	dpb "google.golang.org/protobuf/types/descriptorpb"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

type makeSceneObjectManifestOptions struct {
	gzfGeometryFilenames []string
	metadata             *sompb.SceneObjectMetadata
}

// MakeSceneObjectManifestOption is an option for MakeSceneObjectManifest.
type MakeSceneObjectManifestOption func(*makeSceneObjectManifestOptions)

// WithGZFGeometryFilename appends a value to the GzfGeometryFilenames assets field.
func WithGZFGeometryFilename(filename string) MakeSceneObjectManifestOption {
	return func(opts *makeSceneObjectManifestOptions) {
		opts.gzfGeometryFilenames = append(opts.gzfGeometryFilenames, filename)
	}
}

// WithMetadata specifies the metadata to use for the SceneObjectManifest.
func WithMetadata(metadata *sompb.SceneObjectMetadata) MakeSceneObjectManifestOption {
	return func(opts *makeSceneObjectManifestOptions) {
		opts.metadata = metadata
	}
}

// MakeSceneObjectManifest makes a SceneObjectManifest for testing.
func MakeSceneObjectManifest(t *testing.T, options ...MakeSceneObjectManifestOption) *sompb.SceneObjectManifest {
	t.Helper()

	opts := &makeSceneObjectManifestOptions{
		metadata: &sompb.SceneObjectMetadata{
			Id: &idpb.Id{
				Name:    "some_scene_object",
				Package: "package.some",
			},
			DisplayName: "Some SceneObject",
			Documentation: &documentationpb.Documentation{
				Description: "Some documentation",
			},
			Vendor: &vendorpb.Vendor{
				DisplayName: "Some Company",
			},
		},
	}
	for _, opt := range options {
		opt(opts)
	}

	m := &sompb.SceneObjectManifest{
		Assets:   &sompb.SceneObjectAssets{},
		Metadata: opts.metadata,
	}

	if len(opts.gzfGeometryFilenames) > 0 {
		m.Assets.GzfGeometryFilenames = opts.gzfGeometryFilenames
	}

	return m
}

type makeProcessedSceneObjectManifestOptions struct {
	defaultSceneObjectConfig *socpb.SceneObjectConfig
	fileDescriptorSet        *dpb.FileDescriptorSet
	metadata                 *sompb.SceneObjectMetadata
	sceneObject              *sopb.SceneObject
	userData                 map[string]*anypb.Any
}

// MakeProcessedSceneObjectManifestOption is an option for MakeProcessedSceneObjectManifest.
type MakeProcessedSceneObjectManifestOption func(*makeProcessedSceneObjectManifestOptions)

// WithDefaultSceneObjectConfig specifies the default SceneObjectConfig to use.
func WithDefaultSceneObjectConfig(config *socpb.SceneObjectConfig) MakeProcessedSceneObjectManifestOption {
	return func(opts *makeProcessedSceneObjectManifestOptions) {
		opts.defaultSceneObjectConfig = config
	}
}

// WithProcessedMetadata specifies the metadata to use for the ProcessedSceneObjectManifest.
func WithProcessedMetadata(metadata *sompb.SceneObjectMetadata) MakeProcessedSceneObjectManifestOption {
	return func(opts *makeProcessedSceneObjectManifestOptions) {
		opts.metadata = metadata
	}
}

// WithProcessedFileDescriptorSet specifies the manifest's FileDescriptorSet.
func WithProcessedFileDescriptorSet(fds *dpb.FileDescriptorSet) MakeProcessedSceneObjectManifestOption {
	return func(opts *makeProcessedSceneObjectManifestOptions) {
		opts.fileDescriptorSet = fds
	}
}

// WithProcessedSceneObject specifies the ProcessedSceneObjectManifest's SceneObject.
func WithProcessedSceneObject(so *sopb.SceneObject) MakeProcessedSceneObjectManifestOption {
	return func(opts *makeProcessedSceneObjectManifestOptions) {
		opts.sceneObject = so
	}
}

// WithProcessedSceneObjectName specifies the ProcessedSceneObjectManifest's SceneObject's name.
func WithProcessedSceneObjectName(name string) MakeProcessedSceneObjectManifestOption {
	return func(opts *makeProcessedSceneObjectManifestOptions) {
		opts.sceneObject.Name = name
	}
}

// WithProcessedUserData converts the specifies user data to an Any and adds it to the manifest's
// SceneObject.
func WithProcessedUserData(t *testing.T, key string, value proto.Message) MakeProcessedSceneObjectManifestOption {
	valueAny, err := anypb.New(value)
	if err != nil {
		t.Fatalf("anypb.New(%v) failed: %v", value, err)
	}
	return WithProcessedUserDataAny(key, valueAny)
}

// WithProcessedUserDataAny adds user data to the manifest's SceneObject.
func WithProcessedUserDataAny(key string, valueAny *anypb.Any) MakeProcessedSceneObjectManifestOption {
	return func(opts *makeProcessedSceneObjectManifestOptions) {
		if opts.userData == nil {
			opts.userData = make(map[string]*anypb.Any)
		}
		opts.userData[key] = valueAny
	}
}

// MakeProcessedSceneObjectManifest makes a ProcessedSceneObjectManifest for testing.
func MakeProcessedSceneObjectManifest(t *testing.T, options ...MakeProcessedSceneObjectManifestOption) *sompb.ProcessedSceneObjectManifest {
	t.Helper()

	opts := &makeProcessedSceneObjectManifestOptions{
		defaultSceneObjectConfig: &socpb.SceneObjectConfig{},
		fileDescriptorSet:        &dpb.FileDescriptorSet{},
		metadata: &sompb.SceneObjectMetadata{
			Id: &idpb.Id{
				Name:    "some_scene_object",
				Package: "package.some",
			},
			DisplayName: "Some SceneObject",
			Documentation: &documentationpb.Documentation{
				Description: "Some documentation",
			},
			Vendor: &vendorpb.Vendor{
				DisplayName: "Some Company",
			},
		},
		sceneObject: &sopb.SceneObject{
			Name: "scene_object",
			Entities: []*epb.Entity{
				{
					Name:       "root",
					EntityType: &epb.Entity_Link{},
				},
			},
		},
	}
	for _, opt := range options {
		opt(opts)
	}

	m := &sompb.ProcessedSceneObjectManifest{
		Metadata: opts.metadata,
		Assets: &sompb.ProcessedSceneObjectAssets{
			DefaultSceneObjectConfig: opts.defaultSceneObjectConfig,
			FileDescriptorSet:        opts.fileDescriptorSet,
			SceneObjectModel:         opts.sceneObject,
		},
	}

	if len(opts.userData) > 0 {
		m.Assets.SceneObjectModel.UserData = opts.userData
	}

	return m
}
