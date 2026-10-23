// Copyright 2026 Intrinsic Innovation LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sceneobject

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	docpb "intrinsic/assets/proto/documentation_go_proto"
	idpb "intrinsic/assets/proto/id_go_proto"
	vpb "intrinsic/assets/proto/vendor_go_proto"
	sompb "intrinsic/assets/scene_objects/proto/scene_object_manifest_go_proto"
	sopb "intrinsic/scene/proto/v1/scene_object_go_proto"
)

func TestManifestFromAsset(t *testing.T) {
	metadata := &sompb.SceneObjectMetadata{
		Id: &idpb.Id{
			Package: "ai.intrinsic",
			Name:    "scene_object_asset",
		},
		DisplayName: "Some Scene Object Asset",
		Documentation: &docpb.Documentation{
			Description: "Some description",
		},
		Vendor: &vpb.Vendor{DisplayName: "Vendor"},
	}

	tests := []struct {
		desc    string
		pm      *sompb.ProcessedSceneObjectManifest
		want    *sompb.SceneObjectManifest
		wantErr bool
	}{
		{
			desc: "populated ProcessedSceneObjectManifest",
			pm: &sompb.ProcessedSceneObjectManifest{
				Metadata: metadata,
				Assets: &sompb.ProcessedSceneObjectAssets{
					SceneObjectModel: &sopb.SceneObject{
						Name: "scene_object_asset",
					},
				},
			},
			want: &sompb.SceneObjectManifest{
				Metadata: metadata,
			},
		},
		{
			desc: "empty ProcessedSceneObjectManifest",
			pm:   &sompb.ProcessedSceneObjectManifest{},
			want: &sompb.SceneObjectManifest{},
		},
		{
			desc:    "nil ProcessedSceneObjectManifest",
			pm:      nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := ManifestFromAsset(tc.pm)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ManifestFromAsset() error = %v, wantErr = %v", err, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("ManifestFromAsset() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
