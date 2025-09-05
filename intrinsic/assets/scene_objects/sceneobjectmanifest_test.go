// Copyright 2023 Intrinsic Innovation LLC

package sceneobjectmanifest

import (
	"os"
	"testing"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"intrinsic/util/testing/testio"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	idpb "intrinsic/assets/proto/id_go_proto"
	vendorpb "intrinsic/assets/proto/vendor_go_proto"
	sompb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
)

const (
	boxGZFPath              = "intrinsic/scene/build_defs/tests/box.gzf"
	emptyObjectManifestPath = "intrinsic/assets/scene_objects/build_defs/tests/empty_object_with_user_data_manifest.textproto"
	emptyObjectGZFPath      = "intrinsic/assets/scene_objects/build_defs/tests/empty_object_with_user_data.gzf"
)

func mustReadSceneObjectManifestTextProto(t *testing.T, path string) *sompb.SceneObjectManifest {
	t.Helper()
	b, err := os.ReadFile(testio.MustCreateRunfilePath(t, path))
	if err != nil {
		t.Fatalf("Failed to read SceneObjectManifest: %v", err)
	}
	m := &sompb.SceneObjectManifest{}
	if err := prototext.Unmarshal(b, m); err != nil {
		t.Fatalf("Failed to unmarshal SceneObjectManifest: %v", err)
	}
	return m
}

func TestValidateSceneObjectManifest(t *testing.T) {
	m := &sompb.SceneObjectManifest{
		Metadata: &sompb.SceneObjectMetadata{
			Id: &idpb.Id{
				Name:    "test",
				Package: "package.some",
			},
			DisplayName: "Some Scene Object",
			Vendor: &vendorpb.Vendor{
				DisplayName: "Intrinsic",
			},
		},
	}
	mWithGeometry := &sompb.SceneObjectManifest{
		Metadata: &sompb.SceneObjectMetadata{
			Id: &idpb.Id{
				Name:    "test2",
				Package: "package.some",
			},
			DisplayName: "Some Other Scene Object",
			Vendor: &vendorpb.Vendor{
				DisplayName: "Intrinsic",
			},
		},
		Assets: &sompb.SceneObjectAssets{
			GzfGeometryFilenames: []string{boxGZFPath},
		},
	}
	mWithSceneObjectUserData := mustReadSceneObjectManifestTextProto(t, emptyObjectManifestPath)
	mWithSceneObjectUserData.Assets = &sompb.SceneObjectAssets{
		GzfGeometryFilenames: []string{emptyObjectGZFPath},
	}

	userDataFDS := &descriptorpb.FileDescriptorSet{
		File: []*descriptorpb.FileDescriptorProto{
			&descriptorpb.FileDescriptorProto{
				Name:    proto.String("google.protobuf.Empty"),
				Package: proto.String("google.protobuf"),
				MessageType: []*descriptorpb.DescriptorProto{
					&descriptorpb.DescriptorProto{
						Name: proto.String("Empty"),
					},
				},
			},
		},
	}
	files, err := protodesc.NewFiles(userDataFDS)
	if err != nil {
		t.Fatalf("Failed to create FileDescriptorSet: %v", err)
	}
	badFiles, err := protodesc.NewFiles(&descriptorpb.FileDescriptorSet{})
	if err != nil {
		t.Fatalf("Failed to create FileDescriptorSet: %v", err)
	}

	mInvalidName := proto.Clone(m).(*sompb.SceneObjectManifest)
	mInvalidName.GetMetadata().GetId().Name = "_invalid_name"
	mInvalidPackage := proto.Clone(m).(*sompb.SceneObjectManifest)
	mInvalidPackage.GetMetadata().GetId().Package = "_invalid_package"
	mNoDisplayName := proto.Clone(m).(*sompb.SceneObjectManifest)
	mNoDisplayName.GetMetadata().DisplayName = ""
	mNoVendor := proto.Clone(m).(*sompb.SceneObjectManifest)
	mNoVendor.GetMetadata().Vendor = nil
	mWithMultipleGZFFiles := proto.Clone(m).(*sompb.SceneObjectManifest)
	mWithMultipleGZFFiles.Assets = &sompb.SceneObjectAssets{
		GzfGeometryFilenames: []string{"file1.gzf", "file2.gzf"},
	}
	mWithRoot := proto.Clone(m).(*sompb.SceneObjectManifest)
	mWithRoot.Assets = &sompb.SceneObjectAssets{
		RootSceneObjectName: "root",
	}

	tests := []struct {
		desc    string
		m       *sompb.SceneObjectManifest
		opts    []ValidateSceneObjectManifestOption
		wantErr bool
	}{
		{
			desc: "valid",
			m:    m,
		},
		{
			desc: "valid with geometry",
			m:    mWithGeometry,
			opts: []ValidateSceneObjectManifestOption{
				WithGZFPaths(map[string]string{boxGZFPath: testio.MustCreateRunfilePath(t, boxGZFPath)}),
			},
		},
		{
			desc: "valid with scene object user data",
			m:    mWithSceneObjectUserData,
			opts: []ValidateSceneObjectManifestOption{
				WithGZFPaths(map[string]string{emptyObjectGZFPath: testio.MustCreateRunfilePath(t, emptyObjectGZFPath)}),
				WithFiles(files),
			},
		},
		{
			desc: "scene object user data, missing GZF files",
			m:    mWithSceneObjectUserData,
			opts: []ValidateSceneObjectManifestOption{
				WithFiles(files),
			},
			wantErr: true,
		},
		{
			desc: "scene object user data, missing file descriptors",
			m:    mWithSceneObjectUserData,
			opts: []ValidateSceneObjectManifestOption{
				WithGZFPaths(map[string]string{emptyObjectGZFPath: testio.MustCreateRunfilePath(t, emptyObjectGZFPath)}),
			},
			wantErr: true,
		},
		{
			desc: "scene object user data, invalid file descriptors",
			m:    mWithSceneObjectUserData,
			opts: []ValidateSceneObjectManifestOption{
				WithGZFPaths(map[string]string{emptyObjectGZFPath: testio.MustCreateRunfilePath(t, emptyObjectGZFPath)}),
				WithFiles(badFiles),
			},
			wantErr: true,
		},
		{
			desc:    "invalid name",
			m:       mInvalidName,
			wantErr: true,
		},
		{
			desc:    "invalid package",
			m:       mInvalidPackage,
			wantErr: true,
		},
		{
			desc:    "no display name",
			m:       mNoDisplayName,
			wantErr: true,
		},
		{
			desc:    "no vendor",
			m:       mNoVendor,
			wantErr: true,
		},
		{
			desc:    "multiple gzf files",
			m:       mWithMultipleGZFFiles,
			wantErr: true,
		},
		{
			desc:    "root scene object name specified",
			m:       mWithRoot,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := ValidateSceneObjectManifest(tc.m, tc.opts...)
			if tc.wantErr && err == nil {
				t.Error("ValidateSceneObjectManifest() succeeded, want error")
			} else if !tc.wantErr && err != nil {
				t.Errorf("ValidateSceneObjectManifest() failed: %v", err)
			}
		})
	}
}
