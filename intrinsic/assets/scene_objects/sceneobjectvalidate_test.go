// Copyright 2023 Intrinsic Innovation LLC

package sceneobjectvalidate

import (
	"os"
	"testing"

	sceneobjecttestutils "intrinsic/assets/scene_objects/testing/utils"
	"intrinsic/util/testing/testio"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"

	sompb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"

	descriptorpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	anypb "google.golang.org/protobuf/types/known/anypb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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

func sceneObjectManifestWithUserData(t *testing.T) *sompb.SceneObjectManifest {
	t.Helper()

	m := mustReadSceneObjectManifestTextProto(t, emptyObjectManifestPath)
	m.Assets = &sompb.SceneObjectAssets{
		GzfGeometryFilenames: []string{emptyObjectGZFPath},
	}

	return m
}

func TestSceneObjectManifest(t *testing.T) {
	userDataFDS := &descriptorpb.FileDescriptorSet{
		File: []*descriptorpb.FileDescriptorProto{
			{
				Name:    proto.String("google.protobuf.Empty"),
				Package: proto.String("google.protobuf"),
				MessageType: []*descriptorpb.DescriptorProto{
					{
						Name: proto.String("Empty"),
					},
				},
			},
		},
	}
	files, err := protodesc.NewFiles(userDataFDS)
	if err != nil {
		t.Fatalf("Failed to create files: %v", err)
	}
	badFiles, err := protodesc.NewFiles(&descriptorpb.FileDescriptorSet{})
	if err != nil {
		t.Fatalf("Failed to create bad files: %v", err)
	}

	tests := []struct {
		desc    string
		m       *sompb.SceneObjectManifest
		opts    []SceneObjectManifestOption
		wantErr bool
	}{
		{
			desc: "valid",
			m:    sceneobjecttestutils.MakeSceneObjectManifest(t),
		},
		{
			desc: "valid with geometry",
			m:    sceneobjecttestutils.MakeSceneObjectManifest(t, sceneobjecttestutils.WithGZFGeometryFilename(boxGZFPath)),
			opts: []SceneObjectManifestOption{
				WithGZFPaths(map[string]string{boxGZFPath: testio.MustCreateRunfilePath(t, boxGZFPath)}),
			},
		},
		{
			desc: "valid with scene object user data",
			m:    sceneObjectManifestWithUserData(t),
			opts: []SceneObjectManifestOption{
				WithGZFPaths(map[string]string{emptyObjectGZFPath: testio.MustCreateRunfilePath(t, emptyObjectGZFPath)}),
				WithFiles(files),
			},
		},
		{
			desc: "scene object user data, missing GZF files",
			m:    sceneObjectManifestWithUserData(t),
			opts: []SceneObjectManifestOption{
				WithFiles(files),
			},
			wantErr: true,
		},
		{
			desc: "scene object user data, missing file descriptors",
			m:    sceneObjectManifestWithUserData(t),
			opts: []SceneObjectManifestOption{
				WithGZFPaths(map[string]string{emptyObjectGZFPath: testio.MustCreateRunfilePath(t, emptyObjectGZFPath)}),
			},
			wantErr: true,
		},
		{
			desc: "scene object user data, invalid file descriptors",
			m:    sceneObjectManifestWithUserData(t),
			opts: []SceneObjectManifestOption{
				WithGZFPaths(map[string]string{emptyObjectGZFPath: testio.MustCreateRunfilePath(t, emptyObjectGZFPath)}),
				WithFiles(badFiles),
			},
			wantErr: true,
		},
		{
			desc:    "missing",
			m:       nil,
			wantErr: true,
		},
		{
			desc: "invalid name",
			m: func() *sompb.SceneObjectManifest {
				m := sceneobjecttestutils.MakeSceneObjectManifest(t)
				m.GetMetadata().GetId().Name = "_invalid_name"
				return m
			}(),
			wantErr: true,
		},
		{
			desc: "invalid package",
			m: func() *sompb.SceneObjectManifest {
				m := sceneobjecttestutils.MakeSceneObjectManifest(t)
				m.GetMetadata().GetId().Package = "_invalid_package"
				return m
			}(),
			wantErr: true,
		},
		{
			desc: "no display name",
			m: func() *sompb.SceneObjectManifest {
				m := sceneobjecttestutils.MakeSceneObjectManifest(t)
				m.Metadata.DisplayName = ""
				return m
			}(),
			wantErr: true,
		},
		{
			desc: "no vendor",
			m: func() *sompb.SceneObjectManifest {
				m := sceneobjecttestutils.MakeSceneObjectManifest(t)
				m.Metadata.Vendor = nil
				return m
			}(),
			wantErr: true,
		},
		{
			desc: "multiple gzf files",
			m: func() *sompb.SceneObjectManifest {
				m := sceneobjecttestutils.MakeSceneObjectManifest(t)
				m.Assets = &sompb.SceneObjectAssets{
					GzfGeometryFilenames: []string{"file1.gzf", "file2.gzf"},
				}
				return m
			}(),
			wantErr: true,
		},
		{
			desc: "root scene object name specified",
			m: func() *sompb.SceneObjectManifest {
				m := sceneobjecttestutils.MakeSceneObjectManifest(t)
				m.Assets = &sompb.SceneObjectAssets{
					RootSceneObjectName: "root",
				}
				return m
			}(),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := SceneObjectManifest(tc.m, tc.opts...)
			if tc.wantErr && err == nil {
				t.Error("SceneObjectManifest() succeeded, want error")
			} else if !tc.wantErr && err != nil {
				t.Errorf("SceneObjectManifest() failed: %v", err)
			}
		})
	}
}

func TestProcessedSceneObjectManifest(t *testing.T) {
	userData, err := anypb.New(&emptypb.Empty{})
	if err != nil {
		t.Fatalf("Failed to create user data: %v", err)
	}
	userDataFDS := &descriptorpb.FileDescriptorSet{
		File: []*descriptorpb.FileDescriptorProto{
			{
				Name:    proto.String("google.protobuf.Empty"),
				Package: proto.String("google.protobuf"),
				MessageType: []*descriptorpb.DescriptorProto{
					{
						Name: proto.String("Empty"),
					},
				},
			},
		},
	}

	tests := []struct {
		desc    string
		m       *sompb.ProcessedSceneObjectManifest
		wantErr bool
	}{
		{
			desc: "valid",
			m:    sceneobjecttestutils.MakeProcessedSceneObjectManifest(t),
		},
		{
			desc: "valid with scene object user data",
			m: sceneobjecttestutils.MakeProcessedSceneObjectManifest(t,
				sceneobjecttestutils.WithProcessedUserDataAny("data", userData),
				sceneobjecttestutils.WithProcessedFileDescriptorSet(userDataFDS),
			),
		},
		{
			desc:    "missing",
			m:       nil,
			wantErr: true,
		},
		{
			desc: "fds is missing user data",
			m: sceneobjecttestutils.MakeProcessedSceneObjectManifest(t,
				sceneobjecttestutils.WithProcessedUserDataAny("data", userData),
			),
			wantErr: true,
		},
		{
			desc: "no fds",
			m: sceneobjecttestutils.MakeProcessedSceneObjectManifest(t,
				sceneobjecttestutils.WithProcessedFileDescriptorSet(nil),
			),
			wantErr: true,
		},
		{
			desc: "invalid name",
			m: func() *sompb.ProcessedSceneObjectManifest {
				m := sceneobjecttestutils.MakeProcessedSceneObjectManifest(t)
				m.GetMetadata().GetId().Name = "_invalid_name"
				return m
			}(),
			wantErr: true,
		},
		{
			desc: "invalid package",
			m: func() *sompb.ProcessedSceneObjectManifest {
				m := sceneobjecttestutils.MakeProcessedSceneObjectManifest(t)
				m.GetMetadata().GetId().Package = "_invalid_package"
				return m
			}(),
			wantErr: true,
		},
		{
			desc: "no display name",
			m: func() *sompb.ProcessedSceneObjectManifest {
				m := sceneobjecttestutils.MakeProcessedSceneObjectManifest(t)
				m.Metadata.DisplayName = ""
				return m
			}(),
			wantErr: true,
		},
		{
			desc: "no vendor",
			m: func() *sompb.ProcessedSceneObjectManifest {
				m := sceneobjecttestutils.MakeProcessedSceneObjectManifest(t)
				m.Metadata.Vendor = nil
				return m
			}(),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := ProcessedSceneObjectManifest(tc.m)
			if tc.wantErr && err == nil {
				t.Error("ProcessedSceneObjectManifest() succeeded, want error")
			} else if !tc.wantErr && err != nil {
				t.Errorf("ProcessedSceneObjectManifest() failed: %v", err)
			}
		})
	}
}
