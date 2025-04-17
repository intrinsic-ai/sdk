// Copyright 2023 Intrinsic Innovation LLC

package sceneobjectmanifest

import (
	"testing"

	"google.golang.org/protobuf/proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	vendorpb "intrinsic/assets/proto/vendor_go_proto"
	sompb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
)

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
		wantErr bool
	}{
		{
			desc: "valid",
			m:    m,
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
			err := ValidateSceneObjectManifest(tc.m)
			if tc.wantErr && err == nil {
				t.Error("ValidateSceneObjectManifest() succeeded, want error")
			} else if !tc.wantErr && err != nil {
				t.Errorf("ValidateSceneObjectManifest() failed: %v", err)
			}
		})
	}
}
